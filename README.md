# Miniflux MCP Server

A Model Context Protocol (MCP) server for interacting with Miniflux RSS reader. This server provides tools to manage feeds, entries, users, and categories through the MCP protocol.

## Features

- **Feed Management**: List, create, and refresh RSS/Atom feeds
- **Entry Operations**: Read entries, update status (read/unread/removed)
- **Category Management**: List and organize feed categories
- **Flexible Authentication**: Support for both API key and username/password authentication

## Installation

### Prerequisites

- Go 1.25 or higher
- Access to a Miniflux instance
- Miniflux API key or user credentials

### Build from Source

```bash
git clone <repository-url>
cd miniflux-mcp
go mod download
go build -o miniflux-mcp .
```

### Using Docker

```bash
docker build -t miniflux-mcp .
docker run --env-file .env miniflux-mcp
```

## Configuration

Create a `.env` file based on `.env.example`:

```bash
cp .env.example .env
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `MINIFLUX_URL` | Your Miniflux instance URL | Yes |
| `MINIFLUX_API_KEY` | API key for authentication | Yes* |
| `MINIFLUX_USERNAME` | Username for basic auth | Yes* |
| `MINIFLUX_PASSWORD` | Password for basic auth | Yes* |

*Either use `MINIFLUX_API_KEY` OR both `MINIFLUX_USERNAME` and `MINIFLUX_PASSWORD`

### Getting a Miniflux API Key

1. Log into your Miniflux instance
2. Go to Settings → API Keys
3. Create a new API key
4. Copy the generated key to your configuration

## Usage

### Starting the Server

```bash
# Using environment variables
export MINIFLUX_URL="https://your-miniflux-instance.com"
export MINIFLUX_API_KEY="your_api_key_here"
./miniflux-mcp

# Or using a .env file
./miniflux-mcp
```

### Available Tools

The Miniflux MCP Server provides **40+ tools** covering all Miniflux API functionality:

## Feed Management (10 tools)
- `get_feeds` - Get all RSS/Atom feeds
- `get_feed` - Get a specific feed by ID
- `create_feed` - Add a new RSS/Atom feed
- `delete_feed` - Delete a specific feed
- `refresh_feed` - Manually refresh a specific feed
- `refresh_all_feeds` - Refresh all feeds
- `get_feed_entries` - Get entries from a specific feed
- `get_feed_entry` - Get a specific entry from a feed
- `get_feed_icon` - Get the icon of a specific feed
- `mark_feed_as_read` - Mark all entries in a feed as read

## Entry Management (8 tools)
- `get_entries` - Get entries with optional filtering
- `get_entry` - Get a specific entry by ID
- `update_entry_status` - Update entry status (read/unread/removed)
- `toggle_bookmark` - Toggle bookmark status of an entry
- `save_entry` - Save an entry
- `fetch_original_content` - Fetch original content of an entry
- `mark_all_as_read` - Mark all entries as read for a user
- `get_category_entry` - Get a specific entry from a category

## Category Management (8 tools)
- `get_categories` - Get all feed categories
- `create_category` - Create a new category
- `update_category` - Update a category title
- `delete_category` - Delete a category
- `get_category_feeds` - Get all feeds in a specific category
- `get_category_entries` - Get all entries in a specific category
- `mark_category_as_read` - Mark all entries in a category as read
- `refresh_category` - Refresh all feeds in a category

## User Management (6 tools)
- `get_users` - Get all users
- `get_me` - Get current user information
- `get_user_by_id` - Get a specific user by ID
- `get_user_by_username` - Get a specific user by username
- `create_user` - Create a new user
- `delete_user` - Delete a user

## System & Utility (7 tools)
- `get_version` - Get Miniflux version information
- `healthcheck` - Perform a health check
- `fetch_counters` - Fetch feed counters
- `discover` - Discover feeds from a URL
- `export` - Export feeds as OPML
- `flush_history` - Flush the read history

## API Key Management (3 tools)
- `get_api_keys` - Get all API keys
- `create_api_key` - Create a new API key
- `delete_api_key` - Delete an API key

## Icons & Media (2 tools)
- `get_icon` - Get an icon by ID
- `get_enclosure` - Get an enclosure by ID

## Integration with Claude Desktop

To use this MCP server with Claude Desktop, add the following to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "miniflux": {
      "command": "/path/to/miniflux-mcp",
      "env": {
        "MINIFLUX_URL": "https://your-miniflux-instance.com",
        "MINIFLUX_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
1. Check the Miniflux documentation: https://miniflux.app/docs/
2. Review the MCP specification: https://spec.modelcontextprotocol.io/
3. Open an issue in this repository
