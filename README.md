# Show Service

A REST API service for managing TV shows, built with Go, Gin, and DynamoDB. This service provides endpoints for creating and retrieving show information with comprehensive validation and testing.

## 🏗️ Project Structure

```
show-service/
├── configs/                    # Configuration files
│   └── config.yaml            # Application configuration
├── curl_and_postman/          # API testing examples
├── internal/                  # Private application code
│   ├── config/               # Configuration management
│   │   └── config.go         # Config loading and validation
│   ├── database/             # Database layer
│   │   ├── dynamo.go         # DynamoDB wrapper
│   │   ├── dynamo_test.go    # Database tests
│   │   └── mock_dynamoapi.go # Database mocks
│   ├── domain/               # Business logic models
│   │   ├── model.go          # Data structures and validation
│   │   └── model_test.go     # Domain tests
│   ├── handlers/             # HTTP handlers
│   │   ├── health.go         # Health check endpoint
│   │   ├── shows.go          # Show CRUD endpoints
│   │   └── shows_test.go     # Handler tests
│   ├── repository/           # Data access layer
│   │   ├── show_repo.go      # Show repository implementation
│   │   └── show_repo_test.go # Repository tests
│   └── service/              # Business logic layer
│       ├── show_service.go   # Show service implementation
│       └── show_service_test.go # Service tests
├── scripts/                  # Docker and deployment scripts
├── .mockery.yml             # Mockery configuration
├── Dockerfile               # Application container
├── docker-compose.yml       # Local development environment
├── Makefile                 # Development commands
├── swagger.yaml             # API specification
└── main.go                  # Application entry point
```

## 🎯 Design Decisions

### Architecture Pattern
The application follows a **Clean Architecture** pattern with clear separation of concerns:

- **Handlers**: HTTP request/response handling
- **Service**: Business logic orchestration
- **Repository**: Data access abstraction
- **Domain**: Core business models and validation
- **Database**: External data store integration

### Data Model

#### Show Entity Structure
```go
type Show struct {
    Slug          string       `json:"slug"`                       // Primary Key
    Title         string       `json:"title"`                      // Required
    Country       *string      `json:"country,omitempty"`          // Optional
    Description   *string      `json:"description,omitempty"`      // Optional
    DRM           *bool        `json:"drm,omitempty"`              // Optional
    EpisodeCount  *int         `json:"episodeCount,omitempty"`     // Optional
    Genre         *string      `json:"genre,omitempty"`            // Optional
    Image         *Image       `json:"image,omitempty"`            // Optional
    Language      *string      `json:"language,omitempty"`         // Optional
    NextEpisode   *NextEpisode `json:"nextEpisode,omitempty"`      // Optional
    PrimaryColour *string      `json:"primaryColour,omitempty"`    // Optional
    Seasons       *[]Season    `json:"seasons,omitempty"`          // Optional
    TVChannel     *string      `json:"tvChannel,omitempty"`        // Optional
}
```

#### Primary Key Strategy
- **Slug as Primary Key**: The `slug` field serves as the primary key in DynamoDB
- **Format**: `show/{handle}` where handle contains letters, digits, and dashes
- **Uniqueness**: Attempting to insert a show with an existing slug results in an error
- **Validation**: Enforced via regex pattern matching

#### Alternative Data Stores Considered
While DynamoDB was chosen for this implementation, the architecture supports alternative data stores:
- **PostgreSQL/MySQL**: Traditional relational databases
- **Amazon S3**: JSON storage for simple use cases
- **Redis**: For caching or session storage

### Technology Stack

#### Core Technologies
- **Go 1.25**: Modern Go version with enhanced features
- **Gin Framework**: HTTP web framework for REST APIs
- **DynamoDB**: NoSQL database for scalable data storage
- **AWS SDK v2**: Official AWS SDK for Go

#### Development Tools
- **Mockery**: Mock generation for unit testing
- **Testify**: Testing assertions and utilities
- **Docker**: Containerization for development and deployment

#### Validation & Quality
- **ozzo-validation**: Declarative validation framework
- **go vet**: Static analysis for suspicious code
- **go fmt**: Code formatting enforcement

## 🚀 Quick Start

### Prerequisites
- Go 1.25+
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd show-service
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Start DynamoDB Local**
   ```bash
   make start-dynamo
   # or
   docker-compose up dynamodb-local -d
   ```

4. **Run the application**
   ```bash
   make start
   # or
   docker-compose up --build
   ```

5. **Verify it's running**
   - API: http://localhost:8080
   - Health check: http://localhost:8080/health
   - DynamoDB: http://localhost:8000

## 📷 Sample Screenshots

### make dev

<img src="./screenshots/make_dev.png" alt="Make Dev">

### make start

<img src="./screenshots/make_start1.png" alt="Make Start">

<img src="./screenshots/make_start2.png" alt="Make Start">

<img src="./screenshots/make_start3.png" alt="Make Start">

### DynamoDB running on Docker locally

<img src="./screenshots/dynamodb_docker_local.png" alt="Make Start">

### Manually querying DynamoDB locally

  *** You will need the AWS CLI installed and configured to run this command ***
  ```bash
    export AWS_ACCESS_KEY_ID=dummy
    export AWS_SECRET_ACCESS_KEY=dummy
    export AWS_DEFAULT_REGION=ap-southeast-2

    aws dynamodb scan       --table-name shows-local       --endpoint-url http://localhost:8000
    aws dynamodb list-tables --endpoint-url http://localhost:8000
  ```

<img src="./screenshots/dynamodb_docker_local2.png" alt="Make Start">

### go run main.go

<img src="./screenshots/go_run_main.go.png" alt="Go Run Main Go">

## 📋 API Endpoints

### Health Check
```http
GET /health
```

### Shows Management
```http
GET    /shows     # List all shows
POST   /shows     # Create new shows (batch)
```

### Example Requests

#### Create Shows
```bash
cd curl_for_manual_tests

curl -X POST http://localhost:8080/shows \
      -H "Content-Type: application/json" \
      -d @shows_request.json

{"message":"Shows created successfully"}       
```

<img src="./screenshots/localhost_request1.png" alt="Post Shows">

#### List Shows
```bash
curl http://localhost:8080/shows

{
  "response": [
    {
      "image": "http://catchup.ninemsn.com.au/img/jump-in/shows/Worlds1280.jpg",
      "slug": "show/worlds",
      "title": "World's..."
    },
    {
      "image": "http://catchup.ninemsn.com.au/img/jump-in/shows/TheOriginals1280.jpg",
      "slug": "show/theoriginals",
      "title": "The Originals"
    },
    {
      "image": "http://catchup.ninemsn.com.au/img/jump-in/shows/ToyHunter1280.jpg",
      "slug": "show/toyhunter",
      "title": "Toy Hunter"
    },
    {
      "image": "http://catchup.ninemsn.com.au/img/jump-in/shows/TheTaste1280.jpg",
      "slug": "show/thetaste",
      "title": "The Taste (Le Goût)"
    },
    {
      "image": "http://catchup.ninemsn.com.au/img/jump-in/shows/16KidsandCounting1280.jpg",
      "slug": "show/16kidsandcounting",
      "title": "16 Kids and Counting"
    },
    {
      "image": "http://catchup.ninemsn.com.au/img/jump-in/shows/ScoobyDoo1280.jpg",
      "slug": "show/scoobydoomysteryincorporated",
      "title": "Scooby-Doo! Mystery Incorporated"
    },
    {
      "image": "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg",
      "slug": "show/thunderbirds",
      "title": "Thunderbirds"
    }
  ]
}
```

<img src="./screenshots/localhost_request2.png" alt="Get Shows">

## 🛠️ Development Workflow

### Using Make (Recommended)

```bash
# Development checks (format, vet, test, build)
make dev

# Individual commands
make tidy    # Clean dependencies
make fmt     # Format code
make vet     # Static analysis
make test    # Run tests
make build   # Build application

# Docker commands
make start-dynamo    # Start DynamoDB Local
make start          # Start full stack
make stop           # Stop all services
make logs           # View logs
make logs-app       # View app logs only

# Quick start
make quick-start    # Build and start everything
make reset         # Clean rebuild
```

### Manual Commands

```bash
# Dependencies
go mod tidy

# Code quality
go fmt ./...
go vet ./...

# Testing
go test ./...

# Building
go build -o show-service .

# Running
./show-service

# Docker
docker-compose up --build
```

## 🧪 Testing

### Running Tests
```bash
# Run all tests
make test
# or
go test ./...

# Run specific package tests
go test ./internal/handlers/...
go test ./internal/service/...
go test ./internal/repository/...
```

### Test Structure
- **Unit Tests**: Each package includes comprehensive tests
- **Mock Testing**: Uses Mockery for dependency isolation
- **Table-Driven Tests**: Consistent test patterns across packages
- **Coverage**: Positive, negative, and edge case scenarios

## 🎭 Mock Generation

### Setting up Mockery

Mockery is used for generating mocks in unit tests. The Docker image provides a consistent environment across platforms.

#### Docker Setup

```bash
# Test installation
docker run vektra/mockery --version
```

#### Generating Mocks

**Linux/macOS:**
```bash
# Generate all mocks
docker run -v "$PWD":/src -w /src vektra/mockery --all

# Generate specific mocks
docker run -v "$PWD":/src -w /src vektra/mockery --name ShowService
```

**Windows:**
```bash
# Generate all mocks
docker run -v "%cd%":/src -w /src vektra/mockery --all

# Generate specific mocks
docker run -v "%cd%":/src -w /src vektra/mockery --name ShowService
```

**PowerShell:**
```powershell
# Generate all mocks
docker run -v "${PWD}":/src -w /src vektra/mockery --all
```

### Mockery Configuration

The `.mockery.yml` file controls mock generation:

```yaml
with-expecter: true
packages:
  github.com/marciomarinho/show-service/internal/database:
    interfaces:
      DynamoAPI:
  github.com/marciomarinho/show-service/internal/repository:
    interfaces:
      ShowRepository:
  github.com/marciomarinho/show-service/internal/service:
    interfaces:
      ShowService:
  github.com/marciomarinho/show-service/internal/handlers:
    interfaces:
      ShowHandler:
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENV` | Environment (local/dev/prod) | local |
| `DYNAMODB_REGION` | AWS region | ap-southeast-2 |
| `DYNAMODB_ENDPOINT` | DynamoDB endpoint | http://localhost:8000 |
| `SHOWS_TABLE` | DynamoDB table name | shows-local |
| `APP_COGNITO_USER_POOL_ID` | Cognito User Pool ID | - |
| `APP_COGNITO_CLIENT_ID` | Cognito Client ID | - |
| `APP_COGNITO_REGION` | Cognito region | - |
| `APP_COGNITO_JWKS_URL` | Cognito JWKS URL | Auto-constructed |

### Configuration File

Application configuration is managed through `configs/config.yaml`:

```yaml
env: local
log:
  level: info
dynamodb:
  region: ap-southeast-2
  endpoint_override: http://localhost:8000
  shows_table: shows-local
cognito:
  user_pool_id: your-user-pool-id
  client_id: your-client-id
  region: ap-southeast-2
  jwks_url: ""  # Auto-constructed if empty
## 🔐 Authentication

### JWT Token Validation

The application includes JWT token validation middleware for non-local environments (dev/prod). Authentication is handled via AWS Cognito JWT tokens.

### Authentication Flow

1. **Token Extraction**: Bearer token extracted from `Authorization` header
2. **Environment Check**: Authentication skipped for `local` environment
3. **Token Validation**: JWT signature and claims validation
4. **Scope Validation**: Verify token has required scope for the endpoint **and** that the required scope is in the configured valid scopes list
5. **User Context**: Authenticated user info added to request context

### Protected Endpoints

All endpoints except `/health` require valid JWT authentication:

```http
GET    /shows     # Requires authentication
POST   /shows     # Requires authentication
GET    /health    # Public endpoint
```

### Scope Requirements

Different endpoints require different scopes from the configured `valid_scopes` list:

| Endpoint | Method | Required Scope Pattern |
|----------|--------|------------------------|
| `/shows` | GET | `*.shows.read` |
| `/shows` | POST | `*.shows.write` |

**Note**: The required scope is dynamically determined by finding the first configured scope that contains the pattern (e.g., `shows.read` or `shows.write`).

### Example Authentication

* You will need cognito_id and cognito_secret ( provided via email )

```bash
# Get token from Cognito (example)
BEARER_TOKEN=$(
  curl --location --request POST "https://s2bnkh07ae.execute-api.ap-southeast-2.amazonaws.com/oauth/token" \
    --header "Authorization: Basic <id:secret encoded in base64>" \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --data-urlencode 'grant_type=client_credentials' \
    --data-urlencode 'scope=https://show-service-dev.api/shows.read https://show-service-dev.api/shows.write' \
    | jq -r '.access_token'
)

# Use token in the header for GET requests
curl -X GET https://s2bnkh07ae.execute-api.ap-southeast-2.amazonaws.com/shows \
  -H "Authorization: Bearer $BEARER_TOKEN"

# Use token in the header for POST requests
# Make sure you have the shows_request.json file (payload) in the same directory
# Alternatively, you can also use :
## Posman - https://www.postman.com/
## Insomnia - https://github.com/Kong/insomnia
## etc.
curl --location --request POST https://s2bnkh07ae.execute-api.ap-southeast-2.amazonaws.com/shows \
  --header "Authorization: Bearer $BEARER_TOKEN" \
  --header "Content-Type: application/json" \      
  -d @shows_request.json

```

### Configuration

Configure Cognito settings in `configs/config.yaml`:

```yaml
cognito:
  user_pool_id: ap-southeast-2_XXXXXXXXX
  client_id: 1a2b3c4d5e6f7g8h9i0j
  region: ap-southeast-2
  jwks_url: ""  # Auto-constructed if empty
  valid_scopes:
    - "https://show-service-prod.api/shows.read"
    - "https://show-service-prod.api/shows.write"
    - "https://show-service-prod.api/admin.read"
```

### Development Notes

- **Local Environment**: Authentication is bypassed for development
- **Token Validation**: Currently implements basic format validation
- **Scope Validation**: Validates against configured `valid_scopes` list
- **Configuration**: Required scopes must be present in config file's `valid_scopes` array
- **Production**: Implement full JWT validation with Cognito JWKS
- **Error Responses**: Returns `401 Unauthorized` for invalid/missing tokens, `403 Forbidden` for insufficient scope

### Future Enhancements

For production deployment, implement complete JWT validation:

1. **JWKS Fetching**: Dynamic retrieval of Cognito public keys
2. **Token Caching**: Cache public keys to reduce latency
3. **Claim Validation**: Verify issuer, audience, expiration
4. **Group Authorization**: Role-based access control
5. **Token Refresh**: Handle token expiration gracefully

### Docker Deployment

```bash
# Build production image
make build-prod

# Or manually
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o show-service .

# Run with Docker Compose
docker-compose up -d
```

### AWS Deployment

The application is designed for AWS deployment with:
- **DynamoDB**: Data storage
- **API Gateway**: HTTP API management
- **Lambda**: Serverless compute (if needed)
- **CloudWatch**: Logging and monitoring

## 🤝 Contributing

1. **Code Style**: Follow Go conventions and run `make fmt`
2. **Testing**: Write tests for new features, run `make test`
3. **Mocks**: Update mocks when interfaces change
4. **Documentation**: Update README for significant changes

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Note**: This service assumes `slug` as the primary key for shows. If you need different partitioning strategies or data models, consider the repository pattern which allows for easy data store replacement.
