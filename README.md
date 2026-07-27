# Miniflux MCP Server

A Model Context Protocol (MCP) server for interacting with Miniflux RSS reader. This server provides tools to manage feeds, entries, users, and categories through the MCP protocol using [Miniflux Client](https://github.com/miniflux/v2/tree/main/client).

## Features

- **Feed Management**: List, create, and refresh RSS/Atom feeds
- **Entry Operations**: Read entries, update status (read/unread/removed)
- **Category Management**: List and organize feed categories
- **Flexible Authentication**: Support for both API key and username/password authentication

## Setup

### Getting a Miniflux API Key

1. Log into your Miniflux instance
2. Go to Settings → API Keys
3. Create a new API key
4. Copy the generated key to your configuration

### Miniflux Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `MINIFLUX_URL` | Your Miniflux instance URL | Yes |
| `MINIFLUX_API_KEY` | API key for authentication | Yes* |
| `MINIFLUX_USERNAME` | Username for basic auth | Yes* |
| `MINIFLUX_PASSWORD` | Password for basic auth | Yes* |

*Either use `MINIFLUX_API_KEY` OR both `MINIFLUX_USERNAME` and `MINIFLUX_PASSWORD`

## Local stdio Server

`stdio` is the default transport and is intended for an MCP client that starts the server locally.

### Using Docker

```bash
docker build -t miniflux-mcp .
docker run -i --rm --env-file .env miniflux-mcp
```

Or use the published image:

```bash
docker run -i --rm --env-file .env jwonder/miniflux-mcp:latest
```

### Claude Code (`.mcp.json`)

Add the following `.mcp.json` file to your project root:

```json
{
  "mcpServers": {
    "miniflux": {
      "type": "stdio",
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "MINIFLUX_URL",
        "-e",
        "MINIFLUX_API_KEY",
        "jwonder/miniflux-mcp:latest"
      ],
      "env": {
        "MINIFLUX_URL": "${MINIFLUX_URL}",
        "MINIFLUX_API_KEY": "${MINIFLUX_API_KEY}"
      }
    }
  }
}
```

## Remote Streamable HTTP Server

The remote server exposes a Streamable HTTP MCP endpoint protected by a static Bearer token.

| Variable | Description | Default |
|----------|-------------|---------|
| `MCP_TRANSPORT` | Set to `streamable-http` | `stdio` |
| `MCP_HTTP_ADDR` | HTTP listen address | `:8080` |
| `MCP_HTTP_PATH` | MCP endpoint path | `/mcp` |
| `MCP_AUTH_TOKEN` | Bearer token protecting the MCP endpoint; required in HTTP mode | None |

Set a strong token and start the container with the Streamable HTTP transport:

```bash
export MCP_AUTH_TOKEN='replace-with-a-strong-secret'

docker run --rm \
  -p 8080:8080 \
  --env-file .env \
  -e MCP_TRANSPORT=streamable-http \
  -e MCP_AUTH_TOKEN \
  jwonder/miniflux-mcp:latest
```

Configure an MCP client to connect to `http://your-server:8080/mcp` and send:

```http
Authorization: Bearer replace-with-a-strong-secret
```

### Claude Code (`.mcp.json`)

Add the remote server to your project-level `.mcp.json`:

```json
{
  "mcpServers": {
    "miniflux": {
      "type": "http",
      "url": "https://mcp.example.com/mcp",
      "headers": {
        "Authorization": "Bearer ${MCP_AUTH_TOKEN}"
      }
    }
  }
}
```

The unauthenticated health endpoint is available at `/healthz`. For deployment outside a trusted private network, put the server behind an HTTPS reverse proxy so the Bearer token is encrypted in transit. One server process uses one configured Miniflux identity, so every connected MCP client has that identity's permissions.

## Available Tools

The Miniflux MCP Server provides **40+ tools** covering all Miniflux API functionality, which can be found in the [Miniflux API Reference](https://miniflux.app/docs/api.html#go-client).

### Feed Management (11 tools)
- `get_feeds` - Get all RSS/Atom feeds
- `get_feed` - Get a specific feed by ID
- `create_feed` - Add a new RSS/Atom feed
- `update_feed` - Update an existing feed
- `delete_feed` - Delete a specific feed
- `refresh_feed` - Manually refresh a specific feed
- `refresh_all_feeds` - Refresh all feeds
- `get_feed_entries` - Get entries from a specific feed
- `get_feed_entry` - Get a specific entry from a feed
- `get_feed_icon` - Get the icon of a specific feed
- `mark_feed_as_read` - Mark all entries in a feed as read

### Entry Management (8 tools)
- `get_entries` - Get entries with optional filtering
- `get_entry` - Get a specific entry by ID
- `update_entry_status` - Update entry status (read/unread/removed)
- `toggle_starred` - Toggle starred status of an entry
- `save_entry` - Save an entry
- `fetch_original_content` - Fetch original content of an entry
- `mark_all_as_read` - Mark all entries as read for a user
- `get_category_entry` - Get a specific entry from a category

### Category Management (8 tools)
- `get_categories` - Get all feed categories
- `create_category` - Create a new category
- `update_category` - Update a category title
- `delete_category` - Delete a category
- `get_category_feeds` - Get all feeds in a specific category
- `get_category_entries` - Get all entries in a specific category
- `mark_category_as_read` - Mark all entries in a category as read
- `refresh_category` - Refresh all feeds in a category

### User Management (6 tools)
- `get_users` - Get all users
- `get_me` - Get current user information
- `get_user_by_id` - Get a specific user by ID
- `get_user_by_username` - Get a specific user by username
- `create_user` - Create a new user
- `delete_user` - Delete a user

### System & Utility (7 tools)
- `get_version` - Get Miniflux version information
- `healthcheck` - Perform a health check
- `fetch_counters` - Fetch feed counters
- `discover` - Discover feeds from a URL
- `export` - Export feeds as OPML
- `flush_history` - Flush the read history

### API Key Management (3 tools)
- `get_api_keys` - Get all API keys
- `create_api_key` - Create a new API key
- `delete_api_key` - Delete an API key

### Icons & Media (2 tools)
- `get_icon` - Get an icon by ID
- `get_enclosure` - Get an enclosure by ID

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
1. Check the Miniflux documentation: https://miniflux.app/docs/
2. Review the MCP specification: https://spec.modelcontextprotocol.io/
3. Open an issue in this repository
