# 🎭 Jokes API

A comprehensive REST API serving over 5,000 jokes across 16 categories. Built with Go, featuring a full web interface, GraphQL support, and Swagger documentation.

**Official Site**: [jokes.apimgr.us](https://jokes.apimgr.us)

## About

The Jokes API provides free, open access to 5,160+ jokes across 16 different categories. No authentication required, no API keys needed - just pure comedy at your fingertips.

### Features

- 🎯 **5,160+ Jokes** across 16 categories
- 🌐 **Full Web Interface** with dark/light themes
- 🚀 **REST API** with versioning (`/api/v1`)
- 🔮 **GraphQL** support for flexible queries
- ⚡ **Swagger/OpenAPI** documentation
- 📱 **PWA Support** for offline access
- 🔓 **No Authentication** - freely available
- ⚙️ **File-based Configuration** (YAML)
- 🐳 **Docker** ready with compose
- 📦 **Single Static Binary** with embedded assets

## Production

### Quick Install

#### Linux
```bash
curl -fsSL https://raw.githubusercontent.com/apimgr/jokes/main/scripts/install.sh | sudo bash
```

#### macOS
```bash
curl -fsSL https://raw.githubusercontent.com/apimgr/jokes/main/scripts/install.sh | bash
```

#### Windows (PowerShell as Administrator)
```powershell
irm https://raw.githubusercontent.com/apimgr/jokes/main/scripts/windows.ps1 | iex
```

### Docker Deployment

```bash
# Using docker-compose
docker-compose up -d

# Using docker directly
docker run -d \
  -p 64321:80 \
  -v ./data:/data \
  -v ./config:/config \
  --name jokes-api \
  ghcr.io/apimgr/jokes:latest
```

### Binary Releases

Download pre-built binaries from the [releases page](https://github.com/apimgr/jokes/releases):

- **Linux**: `jokes-api-linux-amd64`, `jokes-api-linux-arm64`
- **macOS**: `jokes-api-darwin-amd64`, `jokes-api-darwin-arm64`
- **Windows**: `jokes-api-windows-amd64.exe`, `jokes-api-windows-arm64.exe`
- **BSD**: `jokes-api-freebsd-amd64`, `jokes-api-freebsd-arm64`

### Configuration

The application auto-creates configuration on first run:

- **Root users**: `/etc/apimgr/jokes/server.yaml`
- **Regular users**: `~/.config/apimgr/jokes/server.yaml`

#### Configuration Options

```yaml
server:
  host: "0.0.0.0"       # Bind address
  port: 64321           # Port (random 64xxx by default)
  rate_limit: 2000      # Requests per hour per IP
```

### CLI Usage

```bash
# Start the server
./jokes-api

# Show help
./jokes-api --help

# Show version
./jokes-api --version

# Check status (no root required)
./jokes-api --status

# Custom port
./jokes-api --port 8080

# Custom address
./jokes-api --address 127.0.0.1

# Custom data directory
./jokes-api --data /path/to/data
```

### Service Management

#### Linux (systemd)
```bash
sudo systemctl status jokes-api
sudo systemctl start jokes-api
sudo systemctl stop jokes-api
sudo systemctl restart jokes-api
sudo journalctl -u jokes-api -f
```

#### macOS (launchd)
```bash
sudo launchctl list | grep jokes
sudo launchctl unload /Library/LaunchDaemons/com.apimgr.jokes.plist
sudo launchctl load /Library/LaunchDaemons/com.apimgr.jokes.plist
tail -f /usr/local/var/log/apimgr/jokes/stdout.log
```

#### Windows (Service)
```powershell
Get-Service -Name JokesAPI
Start-Service -Name JokesAPI
Stop-Service -Name JokesAPI
Restart-Service -Name JokesAPI
```

## API Usage

### REST API

```bash
# Get a random joke
curl http://localhost:64321/api/v1/jokes/random

# Get 5 random jokes
curl http://localhost:64321/api/v1/jokes/random/5

# Get joke by ID
curl http://localhost:64321/api/v1/jokes/1

# Get jokes from specific category
curl "http://localhost:64321/api/v1/jokes/random?limitTo=[nerdy]"

# Exclude explicit content
curl "http://localhost:64321/api/v1/jokes/random?exclude=explicit"

# Replace names
curl "http://localhost:64321/api/v1/jokes/random?firstName=John&lastName=Doe"

# Get all categories
curl http://localhost:64321/api/v1/jokes/categories

# Get statistics
curl http://localhost:64321/api/v1/jokes/count

# Health check
curl http://localhost:64321/healthz
```

### GraphQL

Access GraphQL Playground at: `http://localhost:64321/graphql`

Example query:
```graphql
query {
  randomJoke(limitTo: "[nerdy]", exclude: "explicit") {
    id
    joke
    categories
  }

  stats {
    total
    categories {
      name
      count
    }
  }
}
```

### Swagger/OpenAPI

Interactive API documentation: `http://localhost:64321/swagger`

## Categories

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

## Response Format

All API responses use a consistent wrapper:

```json
{
  "type": "success|error",
  "value": <data or error message>
}
```

## Development

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, for build automation)

### Build from Source

```bash
# Clone repository
git clone https://github.com/apimgr/jokes.git
cd jokes

# Download dependencies
go mod download

# Build for your platform
go build -o jokes-api .

# Or use Make to build for all platforms
make build
```

### Development Server

```bash
# Run with live reload (requires air)
go install github.com/cosmtrek/air@latest
air

# Or run directly
go run main.go
```

### Testing

```bash
# Run tests
make test

# Or use go test directly
go test -v ./...
```

### Building Docker Image

```bash
# Build image
make docker

# Or use docker directly
docker build -t jokes-api .
```

### Project Structure

```
jokes/
├── main.go                 # Application entry point
├── go.mod                  # Go dependencies
├── Makefile                # Build automation
├── Dockerfile              # Docker image definition
├── docker-compose.yml      # Docker compose config
├── Jenkinsfile             # CI/CD pipeline
├── src/
│   ├── config/             # Configuration management
│   ├── models/             # Data models
│   ├── handlers/           # HTTP handlers
│   ├── routes/             # Route definitions
│   ├── web/                # Web frontend
│   │   ├── templates/      # HTML templates
│   │   └── static/         # CSS, JS, assets
│   ├── graphql/            # GraphQL schema & resolvers
│   ├── swagger/            # Swagger/OpenAPI spec
│   └── data/               # Jokes data (JSON)
├── scripts/                # Installation scripts
│   ├── install.sh          # OS-agnostic installer
│   ├── linux.sh            # Linux-specific
│   ├── macos.sh            # macOS-specific
│   └── windows.ps1         # Windows-specific
├── README.md               # This file
├── SPEC.md                 # API specification
└── LICENSE.md              # MIT License
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

### Code Standards

- Use `go fmt` for formatting
- Run `go vet` for static analysis
- Write tests for new features
- Follow existing code structure

## Troubleshooting

### Port Already in Use

Edit `server.yaml` and change the port, or specify a different port via CLI:
```bash
./jokes-api --port 8080
```

### Service Won't Start

Check logs:
- **Linux**: `sudo journalctl -u jokes-api -f`
- **macOS**: `tail -f /usr/local/var/log/apimgr/jokes/stdout.log`
- **Windows**: Check Event Viewer or service logs

### Binary Won't Run

Make sure the binary is executable:
```bash
chmod +x jokes-api
```

### Configuration Not Found

The app will create default configuration automatically. If you need a custom location:
```bash
./jokes-api --config /path/to/config/dir
```

## License

MIT License - see [LICENSE.md](LICENSE.md) for details.

## Support

- **GitHub Issues**: [github.com/apimgr/jokes/issues](https://github.com/apimgr/jokes/issues)
- **Website**: [jokes.apimgr.us](https://jokes.apimgr.us)
- **Documentation**: Full API documentation available at `/docs` endpoint

---

**Made with ❤️ by APIMGR** | Part of the [APIMGR Project](https://github.com/apimgr)
