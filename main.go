package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"miniflux.app/v2/client"
)

type MinifluxServer struct {
	client *client.Client
}

func NewMinifluxServer() *MinifluxServer {
	// Get configuration from environment variables
	baseURL := os.Getenv("MINIFLUX_URL")
	if baseURL == "" {
		log.Fatal("MINIFLUX_URL environment variable is required")
	}

	apiKey := os.Getenv("MINIFLUX_API_KEY")
	username := os.Getenv("MINIFLUX_USERNAME")
	password := os.Getenv("MINIFLUX_PASSWORD")

	if apiKey == "" && (username == "" || password == "") {
		log.Fatal("Either MINIFLUX_API_KEY or both MINIFLUX_USERNAME and MINIFLUX_PASSWORD must be set")
	}

	var minifluxClient *client.Client
	if apiKey != "" {
		minifluxClient = client.NewClient(baseURL, apiKey)
	} else {
		minifluxClient = client.NewClient(baseURL, username, password)
	}

	return &MinifluxServer{
		client: minifluxClient,
	}
}

func (s *MinifluxServer) GetFeeds(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	feeds, err := s.client.Feeds()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch feeds: %v", err)), nil
	}

	feedsJSON, err := json.MarshalIndent(feeds, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal feeds: %v", err)), nil
	}

	return mcp.NewToolResultText(string(feedsJSON)), nil
}

func (s *MinifluxServer) GetEntries(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments

	// Parse optional parameters
	var filter *client.Filter
	if args != nil {
		argsMap, ok := args.(map[string]interface{})
		if ok {
			filter = &client.Filter{}

			if statusStr, ok := argsMap["status"].(string); ok {
				filter.Status = statusStr
			}

			if feedIDFloat, ok := argsMap["feed_id"].(float64); ok {
				feedID := int64(feedIDFloat)
				filter.FeedID = feedID
			}

			if limitFloat, ok := argsMap["limit"].(float64); ok {
				limit := int(limitFloat)
				filter.Limit = limit
			}

			if offsetFloat, ok := argsMap["offset"].(float64); ok {
				offset := int(offsetFloat)
				filter.Offset = offset
			}
		}
	}

	entries, err := s.client.Entries(filter)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch entries: %v", err)), nil
	}

	entriesJSON, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal entries: %v", err)), nil
	}

	return mcp.NewToolResultText(string(entriesJSON)), nil
}

func (s *MinifluxServer) GetEntry(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	entry, err := s.client.Entry(entryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch entry: %v", err)), nil
	}

	entryJSON, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal entry: %v", err)), nil
	}

	return mcp.NewToolResultText(string(entryJSON)), nil
}

func (s *MinifluxServer) UpdateEntryStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("entry_id and status are required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	entryIDFloat, ok := argsMap["entry_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("entry_id must be a number"), nil
	}

	status, ok := argsMap["status"].(string)
	if !ok {
		return mcp.NewToolResultError("status must be a string"), nil
	}

	entryID := int64(entryIDFloat)
	err := s.client.UpdateEntries([]int64{entryID}, status)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update entry status: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Entry %d status updated to: %s", entryID, status)), nil
}

func (s *MinifluxServer) CreateFeed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("feed_url is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	feedURL, ok := argsMap["feed_url"].(string)
	if !ok {
		return mcp.NewToolResultError("feed_url must be a string"), nil
	}

	var categoryID int64 = 1 // Default category
	if categoryIDFloat, ok := argsMap["category_id"].(float64); ok {
		categoryID = int64(categoryIDFloat)
	}

	feedRequest := &client.FeedCreationRequest{
		FeedURL:    feedURL,
		CategoryID: categoryID,
	}

	if crawler, ok := argsMap["crawler"].(bool); ok {
		feedRequest.Crawler = crawler
	}

	if userAgent, ok := argsMap["user_agent"].(string); ok {
		feedRequest.UserAgent = userAgent
	}

	if username, ok := argsMap["username"].(string); ok {
		feedRequest.Username = username
	}

	if password, ok := argsMap["password"].(string); ok {
		feedRequest.Password = password
	}

	createdFeed, err := s.client.CreateFeed(feedRequest)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create feed: %v", err)), nil
	}

	feedJSON, err := json.MarshalIndent(createdFeed, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal created feed: %v", err)), nil
	}

	return mcp.NewToolResultText(string(feedJSON)), nil
}

func (s *MinifluxServer) GetCategories(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	categories, err := s.client.Categories()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch categories: %v", err)), nil
	}

	categoriesJSON, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal categories: %v", err)), nil
	}

	return mcp.NewToolResultText(string(categoriesJSON)), nil
}

func (s *MinifluxServer) RefreshFeed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	err := s.client.RefreshFeed(feedID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to refresh feed: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Feed %d refreshed successfully", feedID)), nil
}

// User Management Methods
func (s *MinifluxServer) GetUsers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	users, err := s.client.Users()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch users: %v", err)), nil
	}

	usersJSON, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal users: %v", err)), nil
	}

	return mcp.NewToolResultText(string(usersJSON)), nil
}

func (s *MinifluxServer) GetMe(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	user, err := s.client.Me()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch current user: %v", err)), nil
	}

	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal user: %v", err)), nil
	}

	return mcp.NewToolResultText(string(userJSON)), nil
}

func (s *MinifluxServer) GetUserByID(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	user, err := s.client.UserByID(userID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch user: %v", err)), nil
	}

	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal user: %v", err)), nil
	}

	return mcp.NewToolResultText(string(userJSON)), nil
}

func (s *MinifluxServer) GetUserByUsername(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("username is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	username, ok := argsMap["username"].(string)
	if !ok {
		return mcp.NewToolResultError("username must be a string"), nil
	}

	user, err := s.client.UserByUsername(username)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch user: %v", err)), nil
	}

	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal user: %v", err)), nil
	}

	return mcp.NewToolResultText(string(userJSON)), nil
}

func (s *MinifluxServer) CreateUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("username and password are required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	username, ok := argsMap["username"].(string)
	if !ok {
		return mcp.NewToolResultError("username must be a string"), nil
	}

	password, ok := argsMap["password"].(string)
	if !ok {
		return mcp.NewToolResultError("password must be a string"), nil
	}

	var isAdmin bool
	if adminVal, ok := argsMap["is_admin"].(bool); ok {
		isAdmin = adminVal
	}

	user, err := s.client.CreateUser(username, password, isAdmin)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create user: %v", err)), nil
	}

	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal user: %v", err)), nil
	}

	return mcp.NewToolResultText(string(userJSON)), nil
}

func (s *MinifluxServer) DeleteUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	err := s.client.DeleteUser(userID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete user: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("User %d deleted successfully", userID)), nil
}

// Category Management Methods
func (s *MinifluxServer) CreateCategory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("title is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	title, ok := argsMap["title"].(string)
	if !ok {
		return mcp.NewToolResultError("title must be a string"), nil
	}

	category, err := s.client.CreateCategory(title)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create category: %v", err)), nil
	}

	categoryJSON, err := json.MarshalIndent(category, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal category: %v", err)), nil
	}

	return mcp.NewToolResultText(string(categoryJSON)), nil
}

func (s *MinifluxServer) UpdateCategory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("category_id and title are required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	categoryIDFloat, ok := argsMap["category_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("category_id must be a number"), nil
	}

	title, ok := argsMap["title"].(string)
	if !ok {
		return mcp.NewToolResultError("title must be a string"), nil
	}

	categoryID := int64(categoryIDFloat)
	category, err := s.client.UpdateCategory(categoryID, title)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update category: %v", err)), nil
	}

	categoryJSON, err := json.MarshalIndent(category, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal category: %v", err)), nil
	}

	return mcp.NewToolResultText(string(categoryJSON)), nil
}

func (s *MinifluxServer) DeleteCategory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("category_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	categoryIDFloat, ok := argsMap["category_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("category_id must be a number"), nil
	}

	categoryID := int64(categoryIDFloat)
	err := s.client.DeleteCategory(categoryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete category: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Category %d deleted successfully", categoryID)), nil
}

func (s *MinifluxServer) GetCategoryFeeds(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("category_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	categoryIDFloat, ok := argsMap["category_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("category_id must be a number"), nil
	}

	categoryID := int64(categoryIDFloat)
	feeds, err := s.client.CategoryFeeds(categoryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch category feeds: %v", err)), nil
	}

	feedsJSON, err := json.MarshalIndent(feeds, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal feeds: %v", err)), nil
	}

	return mcp.NewToolResultText(string(feedsJSON)), nil
}

func (s *MinifluxServer) GetCategoryEntries(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("category_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	categoryIDFloat, ok := argsMap["category_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("category_id must be a number"), nil
	}

	categoryID := int64(categoryIDFloat)

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

	entries, err := s.client.CategoryEntries(categoryID, filter)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch category entries: %v", err)), nil
	}

	entriesJSON, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal entries: %v", err)), nil
	}

	return mcp.NewToolResultText(string(entriesJSON)), nil
}

func (s *MinifluxServer) MarkCategoryAsRead(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("category_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	categoryIDFloat, ok := argsMap["category_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("category_id must be a number"), nil
	}

	categoryID := int64(categoryIDFloat)
	err := s.client.MarkCategoryAsRead(categoryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to mark category as read: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Category %d marked as read", categoryID)), nil
}

func (s *MinifluxServer) RefreshCategory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	if args == nil {
		return mcp.NewToolResultError("category_id is required"), nil
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	categoryIDFloat, ok := argsMap["category_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("category_id must be a number"), nil
	}

	categoryID := int64(categoryIDFloat)
	err := s.client.RefreshCategory(categoryID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to refresh category: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Category %d refreshed successfully", categoryID)), nil
}

func main() {
	minifluxServer := NewMinifluxServer()

	s := server.NewMCPServer(
		"miniflux-mcp",
		"1.0.0",
		server.WithLogging(),
	)

	// Register all tools
	minifluxServer.RegisterAllTools(s)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

