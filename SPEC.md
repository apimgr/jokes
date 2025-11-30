# JOKES API Specification

## Overview

The JOKES API is a REST API that provides access to over 5,000 jokes across 16 different categories. The API is freely accessible without authentication and uses file-based YAML configuration.

## Base Information

- **Version**: 1.0.0
- **Base URL**: `http://localhost:3009`
- **Protocol**: HTTP/HTTPS
- **Response Format**: JSON
- **Authentication**: None (freely available)
- **Rate Limiting**: 2000 requests per hour per IP (configurable)

## CLI Commands

The application provides several command-line flags for configuration and control:

### Display Commands (No Root/Admin Required)

**Show Help**
```bash
./jokes-api --help
```
Displays usage information and available commands.

**Show Version**
```bash
./jokes-api --version
```
Displays the current version of the application.

**Show Status**
```bash
./jokes-api --status
```
Shows server status and health check. Exits with standard exit codes:
- `0`: Service is running and healthy
- `1`: Service is not running
- `2`: Service is running but unhealthy

### Configuration Commands

**Set Data Directory**
```bash
./jokes-api --data /path/to/data
```
Specifies the directory containing jokes data (jokes.json). Defaults to embedded data if not specified.

**Set Config Directory**
```bash
./jokes-api --config /path/to/config
```
Specifies the directory containing server.yaml configuration. Defaults to:
- Root users: `/etc/apimgr/jokes/`
- Regular users: `~/.config/apimgr/jokes/`

**Set Listen Address**
```bash
./jokes-api --address 127.0.0.1
```
Sets the bind address for the server. Default: `0.0.0.0` (all interfaces).

**Set Port**
```bash
./jokes-api --port 8080
```
Sets the listening port. Default: Random port in 64000-64999 range, or as specified in config file.

### Service Management

**Service Control**
```bash
./jokes-api --service {command}
```

Available service commands:

- `start` - Start the service
- `stop` - Stop the service
- `restart` - Restart the service
- `reload` - Reload configuration without restart
- `--install` - Install the service (creates service file/unit)
- `--uninstall` - Uninstall the service (removes service file/unit)
- `--disable` - Disable the service from starting at boot
- `help` - Show service management help

**Supported Service Managers:**

- **Linux systemd** - Creates `/etc/systemd/system/jokes-api.service`
- **Linux runit** - Creates `/etc/sv/jokes-api/run`
- **macOS launchd** - Creates `/Library/LaunchDaemons/com.apimgr.jokes.plist`
- **Windows Service Manager** - Creates Windows Service via sc.exe
- **BSD rc.d** - Creates `/usr/local/etc/rc.d/jokes-api`

**Examples:**
```bash
# Install the service (requires root/admin)
sudo ./jokes-api --service --install

# Start the service
sudo ./jokes-api --service start

# Check status (no root required)
./jokes-api --status

# Restart the service
sudo ./jokes-api --service restart

# Uninstall the service
sudo ./jokes-api --service --uninstall
```

**Note:** Service management commands typically require root/administrator privileges except for status checks.

## Configuration

### Configuration File Location

- **Root users**: `/etc/apimgr/jokes/server.yaml`
- **Regular users**: `~/.config/apimgr/jokes/server.yaml`

### Configuration Schema

```yaml
server:
  host: string      # Bind address (default: "0.0.0.0")
  port: integer     # Port number (default: random 64000-64999)
  rate_limit: integer # Requests per hour (default: 2000)
```

### Auto-Configuration

The application automatically creates the configuration file on first run if it doesn't exist:
1. Detects if running as root or regular user
2. Creates appropriate config directory
3. Generates server.yaml with default settings
4. Selects random available port in 64000-64999 range

## Data Model

### Joke Object

```json
{
  "id": 1,
  "joke": "Chuck Norris doesn't read books...",
  "categories": ["chuck-norris"]
}
```

**Fields:**
- `id` (integer): Unique identifier (1-5160)
- `joke` (string): The joke text
- `categories` (array of strings): Categories this joke belongs to

### Response Wrapper

All API responses use a consistent wrapper:

```json
{
  "type": "success|error",
  "value": <data or error message>
}
```

## Endpoints

### Health Check

**Endpoint**: `GET /healthz`

**Description**: Returns server health status and statistics.

**Response**: 200 OK

```json
{
  "status": "healthy",
  "timestamp": "2025-11-25T16:00:00Z",
  "version": "1.0.0",
  "jokes_loaded": 5160
}
```

### API Documentation

**Endpoint**: `GET /docs`

**Description**: Returns complete API documentation as JSON.

**Response**: 200 OK

### Root Information

**Endpoint**: `GET /`

**Description**: Returns API information and available endpoints.

**Response**: 200 OK

## Jokes Endpoints

### Get Random Joke

**Endpoint**: `GET /api/v1/jokes/random`

**Query Parameters**:
- `firstName` (optional): Replace "Chuck" with this name
- `lastName` (optional): Replace "Norris" with this name
- `category` (optional): Filter by specific category
- `limitTo` (optional): Limit to categories. Format: `[category]` or `[cat1,cat2]`
- `exclude` (optional): Exclude categories. Format: `cat1,cat2`

**Response**: 200 OK

```json
{
  "type": "success",
  "value": {
    "id": 42,
    "joke": "Chuck Norris can divide by zero.",
    "categories": ["chuck-norris"]
  }
}
```

**Example Requests**:
```
GET /api/v1/jokes/random
GET /api/v1/jokes/random?limitTo=[animal,nerdy]
GET /api/v1/jokes/random?exclude=explicit
GET /api/v1/jokes/random?firstName=John&lastName=Doe
```

### Get Multiple Random Jokes

**Endpoint**: `GET /api/v1/jokes/random/:count`

**Path Parameters**:
- `count` (required): Number of jokes to return (1-100)

**Query Parameters**: Same as Get Random Joke

**Response**: 200 OK

```json
{
  "type": "success",
  "value": [
    {
      "id": 1,
      "joke": "...",
      "categories": ["chuck-norris"]
    },
    {
      "id": 2,
      "joke": "...",
      "categories": ["nerdy"]
    }
  ]
}
```

**Example Requests**:
```
GET /api/v1/jokes/random/5
GET /api/v1/jokes/random/10?limitTo=[nerdy,movie]&exclude=explicit
```

### Get Joke by ID

**Endpoint**: `GET /api/v1/jokes/:id`

**Path Parameters**:
- `id` (required): Joke ID (1-5160)

**Query Parameters**:
- `firstName` (optional): Replace "Chuck" with this name
- `lastName` (optional): Replace "Norris" with this name

**Response**: 200 OK

```json
{
  "type": "success",
  "value": {
    "id": 1,
    "joke": "Chuck Norris doesn't read books...",
    "categories": ["chuck-norris"]
  }
}
```

**Error Responses**:
- 400 Bad Request: Invalid ID format
- 404 Not Found: Joke ID not found

**Example Requests**:
```
GET /api/v1/jokes/1
GET /api/v1/jokes/42?firstName=John&lastName=Doe
```

### Get All Jokes

**Endpoint**: `GET /api/v1/jokes/all`

**Query Parameters**:
- `limitTo` (optional): Limit to categories. Format: `[category]` or `[cat1,cat2]`
- `exclude` (optional): Exclude categories. Format: `cat1,cat2`
- `firstName` (optional): Replace "Chuck" with this name
- `lastName` (optional): Replace "Norris" with this name

**Response**: 200 OK

```json
{
  "type": "success",
  "value": {
    "jokes": [
      {
        "id": 1,
        "joke": "...",
        "categories": ["chuck-norris"]
      }
    ],
    "meta": {
      "total_in_database": 5160,
      "returned": 100,
      "limited_to_categories": ["animal", "nerdy"],
      "excluded_categories": ["explicit"]
    }
  }
}
```

**Example Requests**:
```
GET /api/v1/jokes/all
GET /api/v1/jokes/all?limitTo=[lawyer]
GET /api/v1/jokes/all?exclude=explicit,adult
```

### Get Categories

**Endpoint**: `GET /api/v1/jokes/categories`

**Description**: Returns list of all available joke categories.

**Response**: 200 OK

```json
{
  "type": "success",
  "value": [
    "explicit",
    "nerdy",
    "movie",
    "history",
    "animal",
    "food",
    "sports",
    "work",
    "travel",
    "music",
    "medical",
    "lawyer",
    "school",
    "science",
    "chuck-norris",
    "general"
  ]
}
```

### Get Joke Count

**Endpoint**: `GET /api/v1/jokes/count`

**Description**: Returns total joke count and per-category statistics.

**Response**: 200 OK

```json
{
  "type": "success",
  "value": {
    "total": 5160,
    "categories": [
      {
        "name": "animal",
        "count": 670
      },
      {
        "name": "chuck-norris",
        "count": 593
      }
    ]
  }
}
```

## Available Categories

| Category | Count | Description |
|----------|-------|-------------|
| animal | 670 | Animal-related humor |
| chuck-norris | 593 | Classic Chuck Norris jokes |
| nerdy | 554 | Programming and tech jokes |
| food | 515 | Food and cooking humor |
| movie | 510 | Film and entertainment |
| sports | 380 | Athletic and sports jokes |
| explicit | 360 | Adult humor |
| history | 340 | Historical references |
| work | 340 | Office and workplace |
| travel | 180 | Travel and tourism |
| general | 158 | General humor |
| medical | 120 | Healthcare jokes |
| lawyer | 120 | Legal profession humor |
| school | 120 | Education jokes |
| science | 120 | Scientific humor |
| music | 80 | Musical references |

## HTTP Status Codes

- **200 OK**: Request successful
- **400 Bad Request**: Invalid parameters or request format
- **404 Not Found**: Resource not found
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Server error

## Error Responses

All errors follow the standard response format:

```json
{
  "type": "error",
  "value": "Error message describing what went wrong"
}
```

## Query Parameter Formats

### limitTo Parameter

Format: `[category]` or `[category1,category2]`

Examples:
- `limitTo=[animal]` - Only animal jokes
- `limitTo=[animal,nerdy]` - Animal or nerdy jokes

### exclude Parameter

Format: `category1,category2`

Examples:
- `exclude=explicit` - Exclude explicit jokes
- `exclude=explicit,adult` - Exclude multiple categories

## Name Replacement

When using `firstName` and/or `lastName` parameters:

- Replaces "Chuck Norris" with "{firstName} {lastName}"
- Replaces standalone "Chuck" with "{firstName}"
- If only `firstName` provided: replaces with "{firstName} Norris"
- If only `lastName` provided: replaces with "Chuck {lastName}"

Example:
```
Original: "Chuck Norris can divide by zero."
With firstName=John&lastName=Doe: "John Doe can divide by zero."
```

## GraphQL API

The application provides a GraphQL endpoint for flexible querying.

**Endpoint**: `/graphql`

**GraphQL Playground**: Available at `/graphql` (browser interface)

### Schema

```graphql
type Joke {
  id: Int!
  joke: String!
  categories: [String!]!
}

type CategoryStat {
  name: String!
  count: Int!
}

type JokeStats {
  total: Int!
  categories: [CategoryStat!]!
}

type Query {
  joke(id: Int!): Joke
  randomJoke(category: String, exclude: String, limitTo: String): Joke
  randomJokes(count: Int!, exclude: String, limitTo: String): [Joke!]!
  allJokes(limitTo: String, exclude: String): [Joke!]!
  categories: [String!]!
  jokesByCategory(category: String!): [Joke!]!
  stats: JokeStats!
}
```

### Example Queries

**Get Random Joke**
```graphql
query {
  randomJoke {
    id
    joke
    categories
  }
}
```

**Get Random Jokes with Filters**
```graphql
query {
  randomJokes(count: 5, limitTo: "[nerdy,movie]", exclude: "explicit") {
    id
    joke
    categories
  }
}
```

**Get Statistics**
```graphql
query {
  stats {
    total
    categories {
      name
      count
    }
  }
}
```

**Get Joke by ID**
```graphql
query {
  joke(id: 42) {
    id
    joke
    categories
  }
}
```

## Swagger/OpenAPI Documentation

Interactive API documentation is available via Swagger UI.

**Endpoint**: `/swagger`

The Swagger documentation provides:
- Interactive API testing interface
- Complete endpoint documentation
- Request/response schemas
- Example requests and responses
- Try-it-out functionality for all endpoints

**OpenAPI Specification**: OpenAPI 2.0 (Swagger 2.0)

All REST endpoints are fully documented with:
- Path parameters
- Query parameters
- Request bodies (where applicable)
- Response schemas
- Error responses
- Example values

## Web Frontend

The application includes a comprehensive web interface:

**Pages:**
- `/` - Home page with random jokes and quick navigation
- `/browse` - Browse all jokes with pagination
- `/random` - Display random jokes with filters
- `/categories` - Browse jokes by category
- `/api-docs` - API documentation and examples

**Features:**
- Dark theme (Dracula) as default
- Light theme option
- Theme persistence via localStorage
- Responsive mobile design
- PWA support (offline access)
- Custom modals and toast notifications
- No default JavaScript alerts
- Professional UI/UX

**PWA Manifest**: `/manifest.json`
**Service Worker**: `/sw.js`

## Implementation Details

- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Data Storage**: JSON file (src/data/jokes.json), embedded in binary
- **Configuration**: YAML (server.yaml), auto-created on first run
- **Templates**: Go templates with modular structure (header, nav, footer, etc.)
- **Static Assets**: Embedded in binary using Go embed.FS
- **No Database**: All data loaded into memory at startup
- **No Authentication**: Freely accessible API
- **Thread-Safe**: Safe for concurrent requests
- **Port Selection**: Automatic random port in 64000-64999 range
- **Single Binary**: All assets embedded, no external dependencies

## Deployment

The API can be deployed as:

### 1. Standalone Binary
- Single executable file
- All assets embedded
- Auto-creates configuration
- No external dependencies

### 2. Docker Container
```bash
docker run -d -p 64321:80 -v ./data:/data ghcr.io/apimgr/jokes:latest
```

### 3. Service Installation (Built-in)
```bash
sudo ./jokes-api --service --install
```
Supports:
- Linux systemd
- Linux runit
- macOS launchd
- Windows Service Manager
- BSD rc.d

### 4. Installation Scripts
OS-specific installation scripts in `/scripts`:
- `install.sh` - Universal installer (detects OS)
- `linux.sh` - Linux with systemd
- `macos.sh` - macOS with launchd
- `windows.ps1` - Windows with Service Manager

### 5. Cloud/Serverless
Compatible with cloud deployments:
- Cloud VMs (AWS, GCP, Azure)
- Container orchestration (Kubernetes, Docker Swarm)
- PaaS platforms (Heroku, Render, etc.)

Configuration file location depends on the deployment user's permissions (root vs. regular user).

## Platform Support

**Supported Platforms:**
- Linux (AMD64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (AMD64, ARM64)
- FreeBSD (AMD64, ARM64)

**Build System:**
- Makefile with cross-compilation support
- GitHub Actions for releases
- Docker multi-stage builds
- Jenkinsfile for CI/CD

## Rate Limiting

Default: 2000 requests per hour per IP address

Configurable via `server.yaml`:
```yaml
server:
  rate_limit: 2000  # Requests per hour
```

Rate limit headers included in responses:
- `X-RateLimit-Limit` - Maximum requests per hour
- `X-RateLimit-Remaining` - Remaining requests
- `X-RateLimit-Reset` - Time when limit resets

When rate limit is exceeded:
- **Status**: 429 Too Many Requests
- **Response**: `{"type": "error", "value": "Rate limit exceeded"}`
