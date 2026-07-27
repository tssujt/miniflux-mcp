//go:build e2e

package e2e

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
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
