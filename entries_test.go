package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"miniflux.app/v2/client"
)

func TestGetEntriesFilters(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v1/entries" {
			t.Errorf("path = %s, want /v1/entries", r.URL.Path)
		}

		query := r.URL.Query()
		expectedValues := map[string]string{
			"feed_id":          "42",
			"category_id":      "7",
			"limit":            "1",
			"offset":           "2",
			"published_after":  "1700000000",
			"published_before": "1700003600",
			"changed_after":    "1700000100",
			"changed_before":   "1700003500",
			"after_entry_id":   "10",
			"before_entry_id":  "100",
			"search":           "weather test",
			"starred":          client.FilterOnlyStarred,
			"order":            "published_at",
			"direction":        "desc",
			"globally_visible": "true",
		}
		for name, expected := range expectedValues {
			if actual := query.Get(name); actual != expected {
				t.Errorf("%s = %q, want %q", name, actual, expected)
			}
		}
		if actual := query["status"]; !reflect.DeepEqual(actual, []string{"read", "unread"}) {
			t.Errorf("status = %#v, want read and unread", actual)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&client.EntryResultSet{Total: 12}); err != nil {
			t.Errorf("encode response body: %v", err)
		}
	}))
	defer apiServer.Close()

	minifluxServer := &MinifluxServer{client: client.NewClient(apiServer.URL, "test-api-key")}
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"status":           "removed",
				"statuses":         []interface{}{"read", "unread"},
				"feed_id":          float64(42),
				"category_id":      float64(7),
				"limit":            float64(1),
				"offset":           float64(2),
				"published_after":  float64(1700000000),
				"published_before": float64(1700003600),
				"changed_after":    float64(1700000100),
				"changed_before":   float64(1700003500),
				"after_entry_id":   float64(10),
				"before_entry_id":  float64(100),
				"search":           "weather test",
				"starred":          true,
				"order":            "published_at",
				"direction":        "desc",
				"globally_visible": true,
			},
		},
	}

	result, err := minifluxServer.GetEntries(context.Background(), request)
	if err != nil {
		t.Fatalf("GetEntries returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("GetEntries returned tool error: %#v", result.Content)
	}

	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatalf("result content type = %T, want text", result.Content[0])
	}
	var entries client.EntryResultSet
	if err := json.Unmarshal([]byte(textContent.Text), &entries); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if entries.Total != 12 {
		t.Errorf("total = %d, want 12", entries.Total)
	}
}
