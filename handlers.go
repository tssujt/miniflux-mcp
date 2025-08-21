package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"miniflux.app/v2/client"
)

// Feed Management Methods (Additional)
func (s *MinifluxServer) GetFeed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("feed_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	feedIDFloat, ok := argsMap["feed_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("feed_id must be a number"), nil
	}

	feedID := int64(feedIDFloat)
	feed, err := s.client.Feed(feedID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch feed: %v", err)), nil
	}

	feedJSON, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal feed: %v", err)), nil
	}

	return mcp.NewToolResultText(string(feedJSON)), nil
}

func (s *MinifluxServer) DeleteFeed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("feed_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	feedIDFloat, ok := argsMap["feed_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("feed_id must be a number"), nil
	}

	feedID := int64(feedIDFloat)
	err := s.client.DeleteFeed(feedID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete feed: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Feed %d deleted successfully", feedID)), nil
}

func (s *MinifluxServer) GetFeedEntries(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("feed_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	feedIDFloat, ok := argsMap["feed_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("feed_id must be a number"), nil
	}

	feedID := int64(feedIDFloat)

	// Parse optional filter parameters
	var filter *client.Filter
	if statusStr, ok := argsMap["status"].(string); ok {
		filter = &client.Filter{Status: statusStr}
	}
	if limitFloat, ok := argsMap["limit"].(float64); ok {
		if filter == nil {
			filter = &client.Filter{}
		}
		limit := int(limitFloat)
		filter.Limit = limit
	}
	if offsetFloat, ok := argsMap["offset"].(float64); ok {
		if filter == nil {
			filter = &client.Filter{}
		}
		offset := int(offsetFloat)
		filter.Offset = offset
	}

	entries, err := s.client.FeedEntries(feedID, filter)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch feed entries: %v", err)), nil
	}

	entriesJSON, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal entries: %v", err)), nil
	}

	return mcp.NewToolResultText(string(entriesJSON)), nil
}

func (s *MinifluxServer) GetFeedEntry(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("feed_id and entry_id are required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	feedIDFloat, ok := argsMap["feed_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("feed_id must be a number"), nil
	}

	entryIDFloat, ok := argsMap["entry_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("entry_id must be a number"), nil
	}

	feedID := int64(feedIDFloat)
	entryID := int64(entryIDFloat)

	entry, err := s.client.FeedEntry(feedID, entryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch feed entry: %v", err)), nil
	}

	entryJSON, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal entry: %v", err)), nil
	}

	return mcp.NewToolResultText(string(entryJSON)), nil
}

func (s *MinifluxServer) GetFeedIcon(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("feed_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	feedIDFloat, ok := argsMap["feed_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("feed_id must be a number"), nil
	}

	feedID := int64(feedIDFloat)
	icon, err := s.client.FeedIcon(feedID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch feed icon: %v", err)), nil
	}

	iconJSON, err := json.MarshalIndent(icon, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal icon: %v", err)), nil
	}

	return mcp.NewToolResultText(string(iconJSON)), nil
}

func (s *MinifluxServer) MarkFeedAsRead(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("feed_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	feedIDFloat, ok := argsMap["feed_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("feed_id must be a number"), nil
	}

	feedID := int64(feedIDFloat)
	err := s.client.MarkFeedAsRead(feedID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to mark feed as read: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Feed %d marked as read", feedID)), nil
}

func (s *MinifluxServer) RefreshAllFeeds(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	err := s.client.RefreshAllFeeds()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to refresh all feeds: %v", err)), nil
	}

	return mcp.NewToolResultText("All feeds refreshed successfully"), nil
}

// Entry Management Methods (Additional)
func (s *MinifluxServer) GetCategoryEntry(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("category_id and entry_id are required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	categoryIDFloat, ok := argsMap["category_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("category_id must be a number"), nil
	}

	entryIDFloat, ok := argsMap["entry_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("entry_id must be a number"), nil
	}

	categoryID := int64(categoryIDFloat)
	entryID := int64(entryIDFloat)

	entry, err := s.client.CategoryEntry(categoryID, entryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch category entry: %v", err)), nil
	}

	entryJSON, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal entry: %v", err)), nil
	}

	return mcp.NewToolResultText(string(entryJSON)), nil
}

func (s *MinifluxServer) ToggleBookmark(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("entry_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	entryIDFloat, ok := argsMap["entry_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("entry_id must be a number"), nil
	}

	entryID := int64(entryIDFloat)
	err := s.client.ToggleBookmark(entryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to toggle bookmark: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Bookmark toggled for entry %d", entryID)), nil
}

func (s *MinifluxServer) SaveEntry(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("entry_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	entryIDFloat, ok := argsMap["entry_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("entry_id must be a number"), nil
	}

	entryID := int64(entryIDFloat)
	err := s.client.SaveEntry(entryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to save entry: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Entry %d saved successfully", entryID)), nil
}

func (s *MinifluxServer) FetchEntryOriginalContent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("entry_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	entryIDFloat, ok := argsMap["entry_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("entry_id must be a number"), nil
	}

	entryID := int64(entryIDFloat)
	content, err := s.client.FetchEntryOriginalContent(entryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch original content: %v", err)), nil
	}

	return mcp.NewToolResultText(content), nil
}

func (s *MinifluxServer) MarkAllAsRead(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("user_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	userIDFloat, ok := argsMap["user_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("user_id must be a number"), nil
	}

	userID := int64(userIDFloat)
	err := s.client.MarkAllAsRead(userID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to mark all as read: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("All entries marked as read for user %d", userID)), nil
}

// System and Utility Methods
func (s *MinifluxServer) GetVersion(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	version, err := s.client.Version()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch version: %v", err)), nil
	}

	versionJSON, err := json.MarshalIndent(version, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal version: %v", err)), nil
	}

	return mcp.NewToolResultText(string(versionJSON)), nil
}

func (s *MinifluxServer) Healthcheck(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	err := s.client.Healthcheck()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Healthcheck failed: %v", err)), nil
	}

	return mcp.NewToolResultText("Healthcheck passed"), nil
}

func (s *MinifluxServer) FetchCounters(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	counters, err := s.client.FetchCounters()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch counters: %v", err)), nil
	}

	countersJSON, err := json.MarshalIndent(counters, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal counters: %v", err)), nil
	}

	return mcp.NewToolResultText(string(countersJSON)), nil
}

func (s *MinifluxServer) Discover(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("url is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	url, ok := argsMap["url"].(string)
	if !ok {
		return mcp.NewToolResultError("url must be a string"), nil
	}

	subscriptions, err := s.client.Discover(url)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to discover feeds: %v", err)), nil
	}

	subscriptionsJSON, err := json.MarshalIndent(subscriptions, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal subscriptions: %v", err)), nil
	}

	return mcp.NewToolResultText(string(subscriptionsJSON)), nil
}

func (s *MinifluxServer) Export(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	data, err := s.client.Export()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to export: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

func (s *MinifluxServer) FlushHistory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	err := s.client.FlushHistory()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to flush history: %v", err)), nil
	}

	return mcp.NewToolResultText("History flushed successfully"), nil
}

// API Key Management Methods
func (s *MinifluxServer) GetAPIKeys(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKeys, err := s.client.APIKeys()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch API keys: %v", err)), nil
	}

	apiKeysJSON, err := json.MarshalIndent(apiKeys, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal API keys: %v", err)), nil
	}

	return mcp.NewToolResultText(string(apiKeysJSON)), nil
}

func (s *MinifluxServer) CreateAPIKey(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("description is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	description, ok := argsMap["description"].(string)
	if !ok {
		return mcp.NewToolResultError("description must be a string"), nil
	}

	apiKey, err := s.client.CreateAPIKey(description)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create API key: %v", err)), nil
	}

	apiKeyJSON, err := json.MarshalIndent(apiKey, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal API key: %v", err)), nil
	}

	return mcp.NewToolResultText(string(apiKeyJSON)), nil
}

func (s *MinifluxServer) DeleteAPIKey(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("api_key_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	apiKeyIDFloat, ok := argsMap["api_key_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("api_key_id must be a number"), nil
	}

	apiKeyID := int64(apiKeyIDFloat)
	err := s.client.DeleteAPIKey(apiKeyID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete API key: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("API key %d deleted successfully", apiKeyID)), nil
}

// Icon Methods
func (s *MinifluxServer) GetIcon(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("icon_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	iconIDFloat, ok := argsMap["icon_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("icon_id must be a number"), nil
	}

	iconID := int64(iconIDFloat)
	icon, err := s.client.Icon(iconID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch icon: %v", err)), nil
	}

	iconJSON, err := json.MarshalIndent(icon, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal icon: %v", err)), nil
	}

	return mcp.NewToolResultText(string(iconJSON)), nil
}

// Enclosure Methods
func (s *MinifluxServer) GetEnclosure(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("enclosure_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	enclosureIDFloat, ok := argsMap["enclosure_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("enclosure_id must be a number"), nil
	}

	enclosureID := int64(enclosureIDFloat)
	enclosure, err := s.client.Enclosure(enclosureID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch enclosure: %v", err)), nil
	}

	enclosureJSON, err := json.MarshalIndent(enclosure, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal enclosure: %v", err)), nil
	}

	return mcp.NewToolResultText(string(enclosureJSON)), nil
}
