package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ToolDefinition struct {
	Tool    mcp.Tool
	Handler server.ToolHandlerFunc
}

func (s *MinifluxServer) RegisterAllTools(mcpServer *server.MCPServer) {
	tools := []ToolDefinition{
		// Feed Operations
		{
			Tool: mcp.Tool{
				Name:        "get_feeds",
				Description: "Get all RSS/Atom feeds from Miniflux",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.GetFeeds,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_feed",
				Description: "Get a specific feed by ID",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"feed_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the feed to retrieve",
						},
					},
					Required: []string{"feed_id"},
				},
			},
			Handler: s.GetFeed,
		},
		{
			Tool: mcp.Tool{
				Name:        "create_feed",
				Description: "Add a new RSS/Atom feed to Miniflux",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"feed_url": map[string]interface{}{
							"type":        "string",
							"description": "The URL of the RSS/Atom feed to add",
						},
						"category_id": map[string]interface{}{
							"type":        "number",
							"description": "The category ID to assign the feed to (default: 1)",
						},
						"crawler": map[string]interface{}{
							"type":        "boolean",
							"description": "Enable web scraper for full content",
						},
						"user_agent": map[string]interface{}{
							"type":        "string",
							"description": "Custom user agent for feed fetching",
						},
						"username": map[string]interface{}{
							"type":        "string",
							"description": "Username for HTTP basic authentication",
						},
						"password": map[string]interface{}{
							"type":        "string",
							"description": "Password for HTTP basic authentication",
						},
					},
					Required: []string{"feed_url"},
				},
			},
			Handler: s.CreateFeed,
		},
		{
			Tool: mcp.Tool{
				Name:        "delete_feed",
				Description: "Delete a specific feed",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"feed_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the feed to delete",
						},
					},
					Required: []string{"feed_id"},
				},
			},
			Handler: s.DeleteFeed,
		},
		{
			Tool: mcp.Tool{
				Name:        "refresh_feed",
				Description: "Manually refresh a specific feed",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"feed_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the feed to refresh",
						},
					},
					Required: []string{"feed_id"},
				},
			},
			Handler: s.RefreshFeed,
		},
		{
			Tool: mcp.Tool{
				Name:        "refresh_all_feeds",
				Description: "Refresh all feeds",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.RefreshAllFeeds,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_feed_entries",
				Description: "Get entries from a specific feed",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"feed_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the feed",
						},
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Filter by entry status (read, unread, removed)",
						},
						"limit": map[string]interface{}{
							"type":        "number",
							"description": "Limit the number of entries returned",
						},
						"offset": map[string]interface{}{
							"type":        "number",
							"description": "Offset for pagination",
						},
					},
					Required: []string{"feed_id"},
				},
			},
			Handler: s.GetFeedEntries,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_feed_entry",
				Description: "Get a specific entry from a feed",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"feed_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the feed",
						},
						"entry_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the entry",
						},
					},
					Required: []string{"feed_id", "entry_id"},
				},
			},
			Handler: s.GetFeedEntry,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_feed_icon",
				Description: "Get the icon of a specific feed",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"feed_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the feed",
						},
					},
					Required: []string{"feed_id"},
				},
			},
			Handler: s.GetFeedIcon,
		},
		{
			Tool: mcp.Tool{
				Name:        "mark_feed_as_read",
				Description: "Mark all entries in a feed as read",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"feed_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the feed",
						},
					},
					Required: []string{"feed_id"},
				},
			},
			Handler: s.MarkFeedAsRead,
		},

		// Entry Operations
		{
			Tool: mcp.Tool{
				Name:        "get_entries",
				Description: "Get entries (articles) from Miniflux with optional filtering",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Filter by entry status (read, unread, removed)",
						},
						"feed_id": map[string]interface{}{
							"type":        "number",
							"description": "Filter by specific feed ID",
						},
						"limit": map[string]interface{}{
							"type":        "number",
							"description": "Limit the number of entries returned",
						},
						"offset": map[string]interface{}{
							"type":        "number",
							"description": "Offset for pagination",
						},
					},
				},
			},
			Handler: s.GetEntries,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_entry",
				Description: "Get a specific entry by ID",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"entry_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the entry to retrieve",
						},
					},
					Required: []string{"entry_id"},
				},
			},
			Handler: s.GetEntry,
		},
		{
			Tool: mcp.Tool{
				Name:        "update_entry_status",
				Description: "Update the status of an entry (mark as read/unread/removed)",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"entry_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the entry to update",
						},
						"status": map[string]interface{}{
							"type":        "string",
							"description": "New status for the entry (read, unread, removed)",
							"enum":        []string{"read", "unread", "removed"},
						},
					},
					Required: []string{"entry_id", "status"},
				},
			},
			Handler: s.UpdateEntryStatus,
		},
		{
			Tool: mcp.Tool{
				Name:        "toggle_bookmark",
				Description: "Toggle bookmark status of an entry",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"entry_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the entry",
						},
					},
					Required: []string{"entry_id"},
				},
			},
			Handler: s.ToggleBookmark,
		},
		{
			Tool: mcp.Tool{
				Name:        "save_entry",
				Description: "Save an entry",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"entry_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the entry",
						},
					},
					Required: []string{"entry_id"},
				},
			},
			Handler: s.SaveEntry,
		},
		{
			Tool: mcp.Tool{
				Name:        "fetch_original_content",
				Description: "Fetch the original content of an entry",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"entry_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the entry",
						},
					},
					Required: []string{"entry_id"},
				},
			},
			Handler: s.FetchEntryOriginalContent,
		},
		{
			Tool: mcp.Tool{
				Name:        "mark_all_as_read",
				Description: "Mark all entries as read for a user",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the user",
						},
					},
					Required: []string{"user_id"},
				},
			},
			Handler: s.MarkAllAsRead,
		},

		// Category Operations
		{
			Tool: mcp.Tool{
				Name:        "get_categories",
				Description: "Get all feed categories from Miniflux",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.GetCategories,
		},
		{
			Tool: mcp.Tool{
				Name:        "create_category",
				Description: "Create a new category",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"title": map[string]interface{}{
							"type":        "string",
							"description": "The title of the category",
						},
					},
					Required: []string{"title"},
				},
			},
			Handler: s.CreateCategory,
		},
		{
			Tool: mcp.Tool{
				Name:        "update_category",
				Description: "Update a category title",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the category",
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "The new title of the category",
						},
					},
					Required: []string{"category_id", "title"},
				},
			},
			Handler: s.UpdateCategory,
		},
		{
			Tool: mcp.Tool{
				Name:        "delete_category",
				Description: "Delete a category",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the category",
						},
					},
					Required: []string{"category_id"},
				},
			},
			Handler: s.DeleteCategory,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_category_feeds",
				Description: "Get all feeds in a specific category",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the category",
						},
					},
					Required: []string{"category_id"},
				},
			},
			Handler: s.GetCategoryFeeds,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_category_entries",
				Description: "Get all entries in a specific category",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the category",
						},
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Filter by entry status (read, unread, removed)",
						},
						"limit": map[string]interface{}{
							"type":        "number",
							"description": "Limit the number of entries returned",
						},
					},
					Required: []string{"category_id"},
				},
			},
			Handler: s.GetCategoryEntries,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_category_entry",
				Description: "Get a specific entry from a category",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the category",
						},
						"entry_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the entry",
						},
					},
					Required: []string{"category_id", "entry_id"},
				},
			},
			Handler: s.GetCategoryEntry,
		},
		{
			Tool: mcp.Tool{
				Name:        "mark_category_as_read",
				Description: "Mark all entries in a category as read",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the category",
						},
					},
					Required: []string{"category_id"},
				},
			},
			Handler: s.MarkCategoryAsRead,
		},
		{
			Tool: mcp.Tool{
				Name:        "refresh_category",
				Description: "Refresh all feeds in a category",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the category",
						},
					},
					Required: []string{"category_id"},
				},
			},
			Handler: s.RefreshCategory,
		},

		// User Management
		{
			Tool: mcp.Tool{
				Name:        "get_users",
				Description: "Get all users",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.GetUsers,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_me",
				Description: "Get current user information",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.GetMe,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_user_by_id",
				Description: "Get a specific user by ID",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the user",
						},
					},
					Required: []string{"user_id"},
				},
			},
			Handler: s.GetUserByID,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_user_by_username",
				Description: "Get a specific user by username",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "The username of the user",
						},
					},
					Required: []string{"username"},
				},
			},
			Handler: s.GetUserByUsername,
		},
		{
			Tool: mcp.Tool{
				Name:        "create_user",
				Description: "Create a new user",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "The username for the new user",
						},
						"password": map[string]interface{}{
							"type":        "string",
							"description": "The password for the new user",
						},
						"is_admin": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether the user should be an admin",
						},
					},
					Required: []string{"username", "password"},
				},
			},
			Handler: s.CreateUser,
		},
		{
			Tool: mcp.Tool{
				Name:        "delete_user",
				Description: "Delete a user",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the user",
						},
					},
					Required: []string{"user_id"},
				},
			},
			Handler: s.DeleteUser,
		},

		// System & Utility
		{
			Tool: mcp.Tool{
				Name:        "get_version",
				Description: "Get Miniflux version information",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.GetVersion,
		},
		{
			Tool: mcp.Tool{
				Name:        "healthcheck",
				Description: "Perform a health check",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.Healthcheck,
		},
		{
			Tool: mcp.Tool{
				Name:        "fetch_counters",
				Description: "Fetch feed counters",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.FetchCounters,
		},
		{
			Tool: mcp.Tool{
				Name:        "discover",
				Description: "Discover feeds from a URL",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"url": map[string]interface{}{
							"type":        "string",
							"description": "The URL to discover feeds from",
						},
					},
					Required: []string{"url"},
				},
			},
			Handler: s.Discover,
		},
		{
			Tool: mcp.Tool{
				Name:        "export",
				Description: "Export feeds as OPML",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.Export,
		},
		{
			Tool: mcp.Tool{
				Name:        "flush_history",
				Description: "Flush the read history",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.FlushHistory,
		},

		// API Key Management
		{
			Tool: mcp.Tool{
				Name:        "get_api_keys",
				Description: "Get all API keys",
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: map[string]interface{}{},
				},
			},
			Handler: s.GetAPIKeys,
		},
		{
			Tool: mcp.Tool{
				Name:        "create_api_key",
				Description: "Create a new API key",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Description for the API key",
						},
					},
					Required: []string{"description"},
				},
			},
			Handler: s.CreateAPIKey,
		},
		{
			Tool: mcp.Tool{
				Name:        "delete_api_key",
				Description: "Delete an API key",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"api_key_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the API key",
						},
					},
					Required: []string{"api_key_id"},
				},
			},
			Handler: s.DeleteAPIKey,
		},

		// Icons & Media
		{
			Tool: mcp.Tool{
				Name:        "get_icon",
				Description: "Get an icon by ID",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"icon_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the icon",
						},
					},
					Required: []string{"icon_id"},
				},
			},
			Handler: s.GetIcon,
		},
		{
			Tool: mcp.Tool{
				Name:        "get_enclosure",
				Description: "Get an enclosure by ID",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"enclosure_id": map[string]interface{}{
							"type":        "number",
							"description": "The ID of the enclosure",
						},
					},
					Required: []string{"enclosure_id"},
				},
			},
			Handler: s.GetEnclosure,
		},
	}

	// Register all tools
	for _, toolDef := range tools {
		mcpServer.AddTool(toolDef.Tool, toolDef.Handler)
	}
}