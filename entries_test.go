package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"miniflux.app/v2/client"
)

func TestEntryListsUseCompactFeeds(t *testing.T) {
	feed := &client.Feed{
		ID: 42, Title: "Weather", FeedURL: "https://example.com/feed", Disabled: true,
		Category:        &client.Category{ID: 7, Title: "News", UserID: 1},
		ParsingErrorMsg: strings.Repeat("fetch failed ", 100), ParsingErrorCount: 3,
		ScraperRules: "article", Cookie: "session=secret", Username: "user", Password: "secret",
	}
	entry := &client.Entry{
		ID: 123, FeedID: 42, UserID: 1, Feed: feed, Title: "Forecast",
		Content: "<p>Full article content</p>", URL: "https://example.com/article",
		CommentsURL: "https://example.com/comments", Author: "Author", Status: "unread",
		Starred: true, Tags: []string{"weather"}, ReadingTime: 5, Language: "en",
		Date:       time.Date(2026, 9, 4, 0, 0, 0, 0, time.UTC),
		Enclosures: client.Enclosures{&client.Enclosure{ID: 2, URL: "https://example.com/audio", MimeType: "audio/mpeg"}},
	}
	for _, method := range []struct {
		name string
		path string
		call func(*MinifluxServer, context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
	}{
		{"get_entries", "/v1/entries", (*MinifluxServer).GetEntries},
		{"get_feed_entries", "/v1/feeds/42/entries", (*MinifluxServer).GetFeedEntries},
		{"get_category_entries", "/v1/categories/7/entries", (*MinifluxServer).GetCategoryEntries},
	} {
		t.Run(method.name, func(t *testing.T) {
			for _, fixture := range []struct {
				name    string
				entries client.Entries
			}{
				{"populated", client.Entries{entry, &client.Entry{ID: 124, FeedID: 43}, nil}},
				{"empty", client.Entries{}},
				{"null", nil},
			} {
				t.Run(fixture.name, func(t *testing.T) {
					upstream := &client.EntryResultSet{Total: 50, Entries: fixture.entries}
					payload, err := json.Marshal(upstream)
					if err != nil {
						t.Fatal(err)
					}
					var expected map[string]any
					if err := json.Unmarshal(payload, &expected); err != nil {
						t.Fatal(err)
					}
					if fixture.name == "populated" {
						expected["entries"].([]any)[0].(map[string]any)["feed"] = map[string]any{
							"id": float64(42), "title": "Weather", "feed_url": "https://example.com/feed", "disabled": true,
							"category": map[string]any{"id": float64(7), "title": "News"},
						}
					}
					apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.Method != http.MethodGet || r.URL.Path != method.path {
							t.Errorf("request = %s %s, want GET %s", r.Method, r.URL.Path, method.path)
						}
						if r.URL.Query().Get("limit") != "3" || r.URL.Query().Get("status") != "unread" {
							t.Errorf("filters not preserved: %s", r.URL.RawQuery)
						}
						w.Header().Set("Content-Type", "application/json")
						if _, err := w.Write(payload); err != nil {
							t.Error(err)
						}
					}))
					defer apiServer.Close()
					s := &MinifluxServer{client: client.NewClient(apiServer.URL, "test-api-key")}
					result, err := method.call(s, context.Background(), mcp.CallToolRequest{Params: mcp.CallToolParams{
						Arguments: map[string]any{"feed_id": float64(42), "category_id": float64(7), "limit": float64(3), "status": "unread"},
					}})
					if err != nil {
						t.Fatal(err)
					}
					if result.IsError || len(result.Content) != 1 {
						t.Fatalf("unexpected tool result: %#v", result)
					}
					content, ok := mcp.AsTextContent(result.Content[0])
					if !ok {
						t.Fatalf("expected text, got %T", result.Content[0])
					}
					var actual map[string]any
					if err := json.Unmarshal([]byte(content.Text), &actual); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(actual, expected) {
						t.Error("entry list must preserve article fields, total, and null/empty values while replacing nested feed with its compact summary")
					}
				})
			}
		})
	}
}

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
