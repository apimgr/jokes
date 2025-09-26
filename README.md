# JOKES API

A comprehensive REST API serving over 5,000 jokes across 16 categories. Originally inspired by the Internet Chuck Norris Database (ICNDB).

## Features

- 🎯 **5,160+ jokes** across 16 categories
- 🔄 **RESTful API** with proper HTTP status codes
- 🚀 **API versioning** (`/api/v1/*`)
- 🛡️ **Security** with Helmet and CORS
- ⚡ **Rate limiting** (2000 requests/hour)
- 🏥 **Health checks** (`/healthz`)
- 📚 **Auto-generated documentation** (`/docs`)
- 🐳 **Docker support** with compose
- 🔄 **Legacy compatibility** with redirects

## Quick Start

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Start production server
npm start

# Run with Docker
npm run docker:compose
```

## API Endpoints

### Core Endpoints
- `GET /healthz` - Health check
- `GET /docs` - API documentation
- `GET /api/v1/jokes/random` - Random joke
- `GET /api/v1/jokes/{id}` - Specific joke
- `GET /api/v1/jokes/all` - All jokes with filtering

### Query Parameters
- `limitTo=[category1,category2]` - Limit to specific categories
- `exclude=category1,category2` - Exclude categories
- `firstName=Name` - Replace "Chuck" 
- `lastName=Name` - Replace "Norris"

### Examples
```bash
# Random animal joke
curl "http://localhost:3009/api/v1/jokes/random?limitTo=[animal]"

# 5 jokes from nerdy or movie categories, exclude explicit
curl "http://localhost:3009/api/v1/jokes/random/5?limitTo=[nerdy,movie]&exclude=explicit"

# All lawyer jokes with custom name
curl "http://localhost:3009/api/v1/jokes/all?limitTo=[lawyer]&firstName=John&lastName=Doe"

# Using environment variables for dynamic URLs
export JOKES_API_URL="https://your-domain.com:3009"
curl "${JOKES_API_URL}/api/v1/jokes/random"
```

## Categories (5,160 total jokes)

| Category | Jokes | Description |
|----------|-------|-------------|
| animal | 670 | Animal humor |
| chuck-norris | 593 | Classic Chuck Norris |
| nerdy | 554 | Programming/tech |
| food | 515 | Food & cooking |
| movie | 510 | Film & entertainment |
| sports | 380 | Athletics & sports |
| explicit | 360 | Adult humor |
| history | 340 | Historical references |
| work | 340 | Office & workplace |
| travel | 180 | Travel & tourism |
| general | 158 | General humor |
| medical | 120 | Healthcare |
| lawyer | 120 | Legal profession |
| school | 120 | Education |
| science | 120 | Scientific humor |
| music | 80 | Musical references |

## NPM Scripts

### Development
- `npm start` - Production server
- `npm run dev` - Development with auto-reload
- `npm run prod` - Production with NODE_ENV

### Testing & Validation
- `npm test` - Quick health check

### Docker
- `npm run build` - Build image

## Project Structure

```
icndb/
├── src/
│   ├── app.js              # Express app configuration
│   ├── server.js           # Server startup
│   ├── models/
│   │   ├── jokes.js        # Joke data access layer
│   │   └── jokes.json      # Joke database (5,160 jokes)
│   └── routes/
│       ├── jokes.js        # Joke API endpoints
│       ├── health.js       # Health check endpoint
│       ├── docs.js         # Documentation endpoint
│       └── legacy.js       # Legacy redirects
├── scripts/               # Utility scripts
│   └── show-stats.js     # Database statistics
├── docs/                 # Documentation
├── tests/                # Test files
├── public/               # Static files
├── Dockerfile            # Container definition
├── docker-compose.yml    # Orchestration
└── package.json          # Dependencies & scripts
```

## Response Format

All responses follow this format:

```json
{
  "type": "success|error",
  "value": "joke object, array, or error message"
}
```

## Rate Limiting

- **Limit**: 2000 requests per hour per IP
- **Headers**: `X-RateLimit-*` headers included
- **Status**: 429 when exceeded

## Legacy Support

Original ICNDB endpoints (`/jokes/*`) automatically redirect to versioned endpoints (`/api/v1/jokes/*`) with 301 status codes.

## License

MIT License

