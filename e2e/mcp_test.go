//go:build e2e

package e2e

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"time"
)

type rpcClient struct {
	input  io.Writer
	output *bufio.Reader
	nextID int
}

type rpcResponse struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type toolCallResult struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	IsError bool `json:"isError"`
}

type httpRPCClient struct {
	endpoint string
	token    string
	client   *http.Client
	nextID   int
}

func (c *rpcClient) request(method string, params any, result any) error {
	c.nextID++
	request := map[string]any{
		"jsonrpc": "2.0",
		"id":      c.nextID,
		"method":  method,
		"params":  params,
	}
	if err := json.NewEncoder(c.input).Encode(request); err != nil {
		return fmt.Errorf("encode %s request: %w", method, err)
	}

	for {
		line, err := c.output.ReadBytes('\n')
		if err != nil {
			return fmt.Errorf("read %s response: %w", method, err)
		}

		var response rpcResponse
		if err := json.Unmarshal(line, &response); err != nil {
			return fmt.Errorf("decode %s response: %w", method, err)
		}
		if response.ID != c.nextID {
			continue
		}
		if response.Error != nil {
			return fmt.Errorf("%s failed (%d): %s", method, response.Error.Code, response.Error.Message)
		}
		if result == nil {
			return nil
		}
		if err := json.Unmarshal(response.Result, result); err != nil {
			return fmt.Errorf("decode %s result: %w", method, err)
		}
		return nil
	}
}

func (c *rpcClient) notify(method string) error {
	return json.NewEncoder(c.input).Encode(map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
	})
}

func (c *rpcClient) callTool(name string, arguments map[string]any) (string, error) {
	var result toolCallResult
	if err := c.request("tools/call", map[string]any{
		"name":      name,
		"arguments": arguments,
	}, &result); err != nil {
		return "", err
	}
	if result.IsError {
		return "", fmt.Errorf("%s returned a tool error: %+v", name, result.Content)
	}
	for _, content := range result.Content {
		if content.Type == "text" {
			return content.Text, nil
		}
	}
	return "", fmt.Errorf("%s returned no text content", name)
}

func (c *httpRPCClient) request(method string, params any, result any) error {
	c.nextID++
	request := map[string]any{
		"jsonrpc": "2.0",
		"id":      c.nextID,
		"method":  method,
		"params":  params,
	}
	return c.send(request, http.StatusOK, result)
}

func (c *httpRPCClient) notify(method string) error {
	request := map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
	}
	return c.send(request, http.StatusAccepted, nil)
}

func (c *httpRPCClient) send(request map[string]any, expectedStatus int, result any) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("encode %s request: %w", request["method"], err)
	}
	httpRequest, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create %s request: %w", request["method"], err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+c.token)
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("send %s request: %w", request["method"], err)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read %s response: %w", request["method"], err)
	}
	if response.StatusCode != expectedStatus {
		return fmt.Errorf("%s returned HTTP %d: %s", request["method"], response.StatusCode, strings.TrimSpace(string(body)))
	}
	if result == nil {
		return nil
	}

	var rpcResult rpcResponse
	if err := json.Unmarshal(body, &rpcResult); err != nil {
		return fmt.Errorf("decode %s response: %w", request["method"], err)
	}
	if rpcResult.Error != nil {
		return fmt.Errorf("%s failed (%d): %s", request["method"], rpcResult.Error.Code, rpcResult.Error.Message)
	}
	if err := json.Unmarshal(rpcResult.Result, result); err != nil {
		return fmt.Errorf("decode %s result: %w", request["method"], err)
	}
	return nil
}

func (c *httpRPCClient) callTool(name string, arguments map[string]any) (string, error) {
	var result toolCallResult
	if err := c.request("tools/call", map[string]any{
		"name":      name,
		"arguments": arguments,
	}, &result); err != nil {
		return "", err
	}
	if result.IsError {
		return "", fmt.Errorf("%s returned a tool error: %+v", name, result.Content)
	}
	for _, content := range result.Content {
		if content.Type == "text" {
			return content.Text, nil
		}
	}
	return "", fmt.Errorf("%s returned no text content", name)
}

func TestMCPServerWithMiniflux(t *testing.T) {
	serverPath := os.Getenv("MCP_SERVER_PATH")
	if serverPath == "" {
		t.Fatal("MCP_SERVER_PATH is required")
	}
	for _, name := range []string{"MINIFLUX_URL", "MINIFLUX_USERNAME", "MINIFLUX_PASSWORD"} {
		if os.Getenv(name) == "" {
			t.Fatalf("%s is required", name)
		}
	}

	command := exec.Command(serverPath)
	command.Env = os.Environ()
	stdin, err := command.StdinPipe()
	if err != nil {
		t.Fatalf("create server stdin: %v", err)
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		t.Fatalf("create server stdout: %v", err)
	}
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		t.Fatalf("start MCP server: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- command.Wait()
	}()
	t.Cleanup(func() {
		_ = stdin.Close()
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("MCP server exited with error: %v\n%s", err, stderr.String())
			}
		case <-time.After(5 * time.Second):
			_ = command.Process.Kill()
			<-done
			t.Errorf("MCP server did not stop after stdin was closed")
		}
	})

	client := &rpcClient{input: stdin, output: bufio.NewReader(stdout)}
	var initialized struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if err := client.request("initialize", map[string]any{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]any{},
		"clientInfo": map[string]any{
			"name":    "miniflux-mcp-ci",
			"version": "1.0.0",
		},
	}, &initialized); err != nil {
		t.Fatalf("initialize MCP session: %v\n%s", err, stderr.String())
	}
	if initialized.ProtocolVersion == "" {
		t.Fatal("initialize response did not include a protocol version")
	}
	if err := client.notify("notifications/initialized"); err != nil {
		t.Fatalf("send initialized notification: %v", err)
	}

	var listed struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	if err := client.request("tools/list", map[string]any{}, &listed); err != nil {
		t.Fatalf("list tools: %v", err)
	}
	toolNames := make([]string, 0, len(listed.Tools))
	for _, tool := range listed.Tools {
		toolNames = append(toolNames, tool.Name)
	}
	for _, name := range []string{"get_version", "create_category", "get_categories", "delete_category"} {
		if !slices.Contains(toolNames, name) {
			t.Fatalf("tools/list did not include %q", name)
		}
	}

	versionText, err := client.callTool("get_version", map[string]any{})
	if err != nil {
		t.Fatalf("get Miniflux version: %v", err)
	}
	var version struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal([]byte(versionText), &version); err != nil {
		t.Fatalf("decode Miniflux version: %v", err)
	}
	if version.Version == "" {
		t.Fatal("get_version returned an empty version")
	}

	title := fmt.Sprintf("mcp-e2e-%d", time.Now().UnixNano())
	categoryText, err := client.callTool("create_category", map[string]any{"title": title})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	var category struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}
	if err := json.Unmarshal([]byte(categoryText), &category); err != nil {
		t.Fatalf("decode created category: %v", err)
	}
	if category.ID == 0 || category.Title != title {
		t.Fatalf("created category = %+v, want title %q and a non-zero ID", category, title)
	}
	createdCategoryID := category.ID
	t.Cleanup(func() {
		if createdCategoryID == 0 {
			return
		}
		if _, err := client.callTool("delete_category", map[string]any{"category_id": createdCategoryID}); err != nil {
			t.Errorf("clean up category %d: %v", createdCategoryID, err)
		}
	})

	categoriesText, err := client.callTool("get_categories", map[string]any{})
	if err != nil {
		t.Fatalf("get categories: %v", err)
	}
	var categories []struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal([]byte(categoriesText), &categories); err != nil {
		t.Fatalf("decode categories: %v", err)
	}
	found := slices.ContainsFunc(categories, func(candidate struct {
		ID int64 `json:"id"`
	}) bool {
		return candidate.ID == category.ID
	})
	if !found {
		t.Fatalf("get_categories did not return created category %d", category.ID)
	}

	if _, err := client.callTool("delete_category", map[string]any{"category_id": category.ID}); err != nil {
		t.Fatalf("delete category: %v", err)
	}
	createdCategoryID = 0
}

func TestRemoteMCPServerWithMiniflux(t *testing.T) {
	serverPath := os.Getenv("MCP_SERVER_PATH")
	if serverPath == "" {
		t.Fatal("MCP_SERVER_PATH is required")
	}
	for _, name := range []string{"MINIFLUX_URL", "MINIFLUX_USERNAME", "MINIFLUX_PASSWORD"} {
		if os.Getenv(name) == "" {
			t.Fatalf("%s is required", name)
		}
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve HTTP port: %v", err)
	}
	address := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("release HTTP port: %v", err)
	}

	const token = "e2e-secret"
	command := exec.Command(serverPath)
	command.Env = append(os.Environ(),
		"MCP_TRANSPORT=streamable-http",
		"MCP_HTTP_ADDR="+address,
		"MCP_HTTP_PATH=/mcp",
		"MCP_AUTH_TOKEN="+token,
	)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		t.Fatalf("start remote MCP server: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- command.Wait()
	}()
	t.Cleanup(func() {
		if command.ProcessState == nil {
			_ = command.Process.Signal(os.Interrupt)
		}
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = command.Process.Kill()
			<-done
			t.Errorf("remote MCP server did not stop after interrupt")
		}
	})

	baseURL := "http://" + address
	httpClient := &http.Client{Timeout: 5 * time.Second}
	deadline := time.Now().Add(5 * time.Second)
	for {
		response, requestErr := httpClient.Get(baseURL + "/healthz")
		if requestErr == nil {
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				break
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("remote MCP server did not become healthy: %v\n%s", requestErr, stderr.String())
		}
		time.Sleep(100 * time.Millisecond)
	}

	unauthorizedRequest, err := http.NewRequest(http.MethodPost, baseURL+"/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`))
	if err != nil {
		t.Fatalf("create unauthorized request: %v", err)
	}
	unauthorizedRequest.Header.Set("Content-Type", "application/json")
	unauthorizedResponse, err := httpClient.Do(unauthorizedRequest)
	if err != nil {
		t.Fatalf("send unauthorized request: %v", err)
	}
	_ = unauthorizedResponse.Body.Close()
	if unauthorizedResponse.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthorized request returned HTTP %d, want 401", unauthorizedResponse.StatusCode)
	}

	client := &httpRPCClient{
		endpoint: baseURL + "/mcp",
		token:    token,
		client:   httpClient,
	}
	var initialized struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if err := client.request("initialize", map[string]any{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]any{},
		"clientInfo": map[string]any{
			"name":    "miniflux-mcp-http-ci",
			"version": "1.0.0",
		},
	}, &initialized); err != nil {
		t.Fatalf("initialize remote MCP session: %v", err)
	}
	if initialized.ProtocolVersion == "" {
		t.Fatal("initialize response did not include a protocol version")
	}
	if err := client.notify("notifications/initialized"); err != nil {
		t.Fatalf("send initialized notification: %v", err)
	}

	var listed struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	if err := client.request("tools/list", map[string]any{}, &listed); err != nil {
		t.Fatalf("list remote tools: %v", err)
	}
	toolNames := make([]string, 0, len(listed.Tools))
	for _, tool := range listed.Tools {
		toolNames = append(toolNames, tool.Name)
	}
	for _, name := range []string{"get_version", "create_category", "delete_category"} {
		if !slices.Contains(toolNames, name) {
			t.Fatalf("tools/list did not include %q", name)
		}
	}

	title := fmt.Sprintf("mcp-http-e2e-%d", time.Now().UnixNano())
	categoryText, err := client.callTool("create_category", map[string]any{"title": title})
	if err != nil {
		t.Fatalf("create category through remote MCP: %v", err)
	}
	var category struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}
	if err := json.Unmarshal([]byte(categoryText), &category); err != nil {
		t.Fatalf("decode created category: %v", err)
	}
	if category.ID == 0 || category.Title != title {
		t.Fatalf("created category = %+v, want title %q and a non-zero ID", category, title)
	}
	createdCategoryID := category.ID
	t.Cleanup(func() {
		if createdCategoryID == 0 {
			return
		}
		if _, err := client.callTool("delete_category", map[string]any{"category_id": createdCategoryID}); err != nil {
			t.Errorf("clean up category %d: %v", createdCategoryID, err)
		}
	})

	if _, err := client.callTool("delete_category", map[string]any{"category_id": category.ID}); err != nil {
		t.Fatalf("delete category through remote MCP: %v", err)
	}
	createdCategoryID = 0
}
