# Go Template

A modern, production-ready Golang project template featuring clean architecture, database migrations, testing utilities, and development tooling.

## Features

- 🏗️ **Clean Architecture**: Well-organized project structure with domain-driven design
- 🗃️ **Database Support**: PostgreSQL integration with migrations and connection pooling
- 🛠️ **Code Generation**: SQLC for type-safe SQL queries
- 🧪 **Testing**: Comprehensive testing setup with Docker integration
- 🔍 **Quality Tools**: Linting, security scanning, and code coverage
- 🐳 **Docker**: Optimized multi-stage builds with security best practices
- ⚙️ **Configuration**: Environment-based configuration management
- 📦 **Multiple Binaries**: Service, CLI, and Worker applications

## Project Structure

```
├── cmd/                    # Application entry points
│   ├── service/           # Main API service
│   ├── cli/               # Command-line interface
│   └── worker/            # Background worker
├── internal/              # Private application code
│   ├── api/               # HTTP handlers and routing
│   ├── config/            # Configuration management
│   └── repository/        # Data access layer
├── domain/                # Business logic and entities
│   ├── entities/          # Domain models
│   ├── example/           # Example implementations
│   └── errors.go          # Domain errors
├── build/                 # Compiled binaries
├── docker-compose.yaml    # Development environment
├── Makefile              # Development commands
├── sqlc.yaml             # SQLC configuration
└── change_repo.sh        # Repository name change script
```

## Quick Start

### 1. Use This Template

Click "Use this template" on GitHub or clone the repository:

```bash
git clone https://github.com/guilhermebr/go-template.git your-project-name
cd your-project-name
```

### 2. Change Repository Name

Run the provided script to update all references to the new repository name:

```bash
./change_repo.sh your-new-repo-name
```

This script will:
- Update the module name in `go.mod`
- Replace all import paths throughout the codebase
- Update any references in configuration files

### 3. Install Dependencies

Set up the development environment:

```bash
make setup
```

This will install all required tools:
- `golangci-lint` - Linting
- `migrate` - Database migrations
- `sqlc` - SQL code generation
- `gotestfmt` - Test output formatting
- `gosec` - Security analysis
- `moq` - Mock generation

### 4. Start Development Environment

Start the PostgreSQL database:

```bash
docker-compose up -d db
```

Run database migrations:

```bash
make migration/up
```

### 5. Build and Run

Compile the service:

```bash
make compile
```

Run the service:

```bash
./build/service
```

## Development

### Environment Configuration

Copy and modify the environment file:

```bash
cp .env-dev .env
```

Key configuration options:
- `DATABASE_*`: PostgreSQL connection settings
- `API_ADDRESS`: Service bind address (default: `0.0.0.0:3000`)
- `ENVIRONMENT`: Runtime environment (`development`/`production`)

### Available Make Commands

#### Setup & Installation
- `make setup` - Install all development tools
- `make install-*` - Install specific tools

#### Code Generation
- `make generate` - Generate all code (SQLC, mocks, etc.)
- `make sqlc-generate` - Generate SQLC code only

#### Building
- `make compile` - Build the service binary

#### Testing
- `make test` - Run tests (short mode)
- `make test-full` - Run all tests including integration tests
- `make coverage` - Generate test coverage report

#### Quality Assurance
- `make lint` - Run linters
- `make gosec` - Run security analysis

#### Database Migrations
- `make migration/create` - Create new migration files
- `make migration/up` - Apply all pending migrations
- `make migration/down` - Rollback all migrations

### Database Migrations

Create a new migration:

```bash
make migration/create
# Enter migration name when prompted
```

Apply migrations:

```bash
# Set environment variables first
export DATABASE_HOST=localhost
export DATABASE_USER=postgres
export DATABASE_PASSWORD=postgres
export DATABASE_NAME=app

make migration/up
```

### Code Generation

This template uses SQLC for type-safe database queries. After modifying SQL files:

```bash
make generate
```

## Testing

### Unit Tests

Run quick tests:

```bash
make test
```

### Integration Tests

Run full test suite including integration tests:

```bash
make test-full
```

### Test Coverage

Generate coverage report:

```bash
make coverage
```

## Dependencies

This template includes several carefully selected dependencies:

### Core Dependencies
- **Chi Router** (`github.com/go-chi/chi/v5`) - HTTP router
- **pgx** (`github.com/jackc/pgx/v5`) - PostgreSQL driver
- **migrate** (`github.com/golang-migrate/migrate/v4`) - Database migrations
- **conf** (`github.com/ardanlabs/conf/v3`) - Configuration management
- **uuid** (`github.com/gofrs/uuid/v5`) - UUID generation

### Development & Testing
- **testify** (`github.com/stretchr/testify`) - Testing toolkit
- **dockertest** (`github.com/ory/dockertest/v3`) - Docker integration testing
- **godotenv** (`github.com/joho/godotenv`) - Environment file loading

## Docker

This template includes optimized Docker configurations for both development and production use.

### Docker Files

- **`Dockerfile`**: Standard production build with Alpine Linux base
- **`Dockerfile.prod`**: Ultra-secure production build using Google's distroless image
- **`Dockerfile.migrations`**: Separate container for running database migrations

### Building Docker Images

#### Standard Production Build
```bash
# Build the image
docker build -t go-template:latest .

# Run the container
docker run -p 3000:3000 \
  -e DATABASE_HOST=your-db-host \
  -e DATABASE_USER=your-db-user \
  -e DATABASE_PASSWORD=your-db-password \
  -e DATABASE_NAME=your-db-name \
  go-template:latest
```

#### Ultra-Secure Production Build
```bash
# Build with distroless base (smallest, most secure)
docker build -f Dockerfile.prod -t go-template:prod .

# Run the container
docker run -p 3000:3000 \
  -e DATABASE_HOST=your-db-host \
  -e DATABASE_USER=your-db-user \
  -e DATABASE_PASSWORD=your-db-password \
  -e DATABASE_NAME=your-db-name \
  go-template:prod
```

### Docker Compose

#### Development Environment
```bash
# Start all services (database + migrations)
docker-compose up

# Start only the database
docker-compose up -d db

# Run migrations separately
docker-compose run migrations
```

#### Production Deployment
```bash
# Create production docker-compose.prod.yaml and run
docker-compose -f docker-compose.prod.yaml up -d
```

### Image Optimization Features

The Docker builds include several optimizations:

- **Multi-stage builds**: Separate build and runtime stages for smaller final images
- **Layer caching**: Go modules are downloaded in a separate layer for better cache utilization
- **Security**: Non-root user execution and minimal attack surface
- **Static binaries**: Fully static Go binaries for distroless compatibility
- **Health checks**: Built-in health monitoring for the standard Dockerfile
- **Build optimization**: Stripped binaries with `-ldflags="-w -s"` for smaller size

### Image Sizes

Approximate final image sizes:
- **Dockerfile**: ~20MB (Alpine-based)
- **Dockerfile.prod**: ~15MB (Distroless-based)
- **Development**: Full Go toolchain (~800MB)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

