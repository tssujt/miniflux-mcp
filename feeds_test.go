package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"miniflux.app/v2/client"
)

func TestGetFeedsReturnsCompactSummaries(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v1/feeds" {
			t.Errorf("path = %s, want /v1/feeds", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{
				"id":                    42,
				"title":                 "Weather",
				"feed_url":              "https://example.com/feed.xml",
				"disabled":              true,
				"parsing_error_message": strings.Repeat("long parsing error ", 100),
				"scraper_rules":         "article",
				"username":              "feed-user",
				"password":              "feed-password",
				"category": map[string]any{
					"id":      7,
					"title":   "News",
					"user_id": 1,
				},
			},
			{
				"id":                    43,
				"title":                 "Uncategorized",
				"parsing_error_message": "another error",
			},
		}); err != nil {
			t.Errorf("encode response body: %v", err)
		}
	}))
	defer apiServer.Close()

	minifluxServer := &MinifluxServer{client: client.NewClient(apiServer.URL, "test-api-key")}
	result, err := minifluxServer.GetFeeds(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("GetFeeds returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("GetFeeds returned tool error: %#v", result.Content)
	}
	if len(result.Content) != 1 {
		t.Fatalf("result content length = %d, want 1", len(result.Content))
	}

	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatalf("result content type = %T, want text", result.Content[0])
	}

	const expected = `[
  {
    "id": 42,
    "title": "Weather",
    "feed_url": "https://example.com/feed.xml",
    "disabled": true,
    "category": {
      "id": 7,
      "title": "News"
    }
  },
  {
    "id": 43,
    "title": "Uncategorized",
    "feed_url": "",
    "disabled": false
  }
]`
	if textContent.Text != expected {
		t.Errorf("GetFeeds result = %s, want %s", textContent.Text, expected)
	}
}
