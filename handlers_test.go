package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"miniflux.app/v2/client"
)

func TestUpdateFeed(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/v1/feeds/42" {
			t.Errorf("path = %s, want /v1/feeds/42", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		if body["category_id"] != float64(7) {
			t.Errorf("category_id = %#v, want 7", body["category_id"])
		}
		if body["title"] != "Updated title" {
			t.Errorf("title = %#v, want Updated title", body["title"])
		}
		if body["crawler"] != false {
			t.Errorf("crawler = %#v, want false", body["crawler"])
		}
		if body["feed_url"] != nil {
			t.Errorf("feed_url = %#v, want null", body["feed_url"])
		}
		if _, exists := body["feed_id"]; exists {
			t.Errorf("request body must not contain feed_id: %#v", body)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&client.Feed{
			ID:       42,
			Title:    "Updated title",
			Crawler:  false,
			Category: &client.Category{ID: 7},
		}); err != nil {
			t.Errorf("encode response body: %v", err)
		}
	}))
	defer apiServer.Close()

	minifluxServer := &MinifluxServer{client: client.NewClient(apiServer.URL, "test-api-key")}
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"feed_id":     float64(42),
				"category_id": float64(7),
				"title":       "Updated title",
				"crawler":     false,
			},
		},
	}

	result, err := minifluxServer.UpdateFeed(context.Background(), request)
	if err != nil {
		t.Fatalf("UpdateFeed returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("UpdateFeed returned tool error: %#v", result.Content)
	}
	if len(result.Content) != 1 {
		t.Fatalf("result content length = %d, want 1", len(result.Content))
	}

	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatalf("result content type = %T, want text", result.Content[0])
	}
	var updatedFeed client.Feed
	if err := json.Unmarshal([]byte(textContent.Text), &updatedFeed); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if updatedFeed.ID != 42 || updatedFeed.Category == nil || updatedFeed.Category.ID != 7 {
		t.Errorf("updated feed = %#v, want feed 42 in category 7", updatedFeed)
	}
}

func TestUpdateFeedRequiresAChange(t *testing.T) {
	minifluxServer := &MinifluxServer{}
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"feed_id": float64(42),
			},
		},
	}

	result, err := minifluxServer.UpdateFeed(context.Background(), request)
	if err != nil {
		t.Fatalf("UpdateFeed returned error: %v", err)
	}
	if !result.IsError {
		t.Fatal("UpdateFeed succeeded without a field to update")
	}
}

func TestToggleStarredUsesStarEndpoint(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/v1/entries/42/star" {
			t.Errorf("path = %s, want /v1/entries/42/star", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer apiServer.Close()

	minifluxServer := &MinifluxServer{client: client.NewClient(apiServer.URL, "test-api-key")}
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"entry_id": float64(42),
			},
		},
	}

	result, err := minifluxServer.ToggleStarred(context.Background(), request)
	if err != nil {
		t.Fatalf("ToggleStarred returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("ToggleStarred returned tool error: %#v", result.Content)
	}
}
