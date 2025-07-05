# Personal Finance Tracker

A comprehensive personal finance tracking application built with **Go**, **HTMX**, **Tailwind CSS**, and **PostgreSQL**. This project demonstrates modern web development patterns and clean architecture principles.

> **Built with AI Assistance**: This project was developed in collaboration with **Claude Sonnet 4** to showcase how AI can accelerate full-stack development while maintaining high code quality and best practices.

## 🎯 The Original Vision

### Started from Go Template
```bash
# Clone the go-template repository
git clone https://github.com/guilhermebr/go-template.git finance
cd finance

# Run the change repository script
./change_repo.sh finance
```

This project started with a simple request over the template:

> *"I want to build a personal finance application that works like a ledger, with accounts, transactions, balances, and categories (like Salary for income, Credit Card/Groceries/Food for expenses). I want both a backend service and a frontend with a nice UI using Go templates with HTMX and Tailwind CSS."*

## 🚀 What We Built

### Core Features

- 💰 **Account Management**: Multiple account types (checking, savings, credit cards, investments, cash)
- 📊 **Smart Categories**: Organized income and expense categorization with color coding
- 💸 **Transaction Tracking**: Complete transaction history with status tracking (pending, cleared, cancelled)
- 🧮 **Automatic Balance Calculations**: Real-time balance updates using database triggers
- 🎨 **Modern Web Interface**: Responsive UI built with Tailwind CSS and HTMX
- 🔄 **Real-time Updates**: Seamless user experience with HTMX partial updates
- 🛡️ **Type Safety**: SQLC-generated database code for compile-time safety

### Architecture

- **Clean Architecture**: Domain-driven design with proper separation of concerns
- **RESTful API**: Complete REST endpoints for all operations
- **Database-First**: PostgreSQL with proper relationships and constraints
- **Modern Frontend**: Server-side rendering with HTMX for dynamic interactions

## 🏗️ Development Journey

### How This Project Was Created

This project showcases the power of AI-assisted development. Here's how we built it step by step:

#### Phase 1: Foundation & Domain Design
1. **Domain Entities**: Created Account, Category, Transaction, and Balance entities
2. **Use Cases**: Implemented business logic for each domain with proper validation
3. **Repository Pattern**: Database abstraction layer with interfaces

#### Phase 2: Database & Persistence  
1. **PostgreSQL Schema**: Designed relational database with proper constraints
2. **SQLC Integration**: Type-safe SQL query generation
3. **Database Triggers**: Automatic balance calculation triggers
4. **Sample Data**: Pre-populated categories for immediate use

#### Phase 3: API Layer
1. **REST Endpoints**: Full CRUD operations for all entities
2. **Error Handling**: Comprehensive error handling and validation
3. **JSON DTOs**: Proper request/response data transfer objects

#### Phase 4: Web Frontend
1. **Go Templates**: Server-side rendered HTML templates
2. **HTMX Integration**: Dynamic interactions without JavaScript
3. **Tailwind CSS**: Modern, responsive styling
4. **Component Architecture**: Reusable template components

#### Phase 5: Integration & Testing
1. **Database Setup**: Docker PostgreSQL container
2. **Environment Configuration**: Proper config management
3. **End-to-End Testing**: Full application testing

## 🛠️ Technology Stack

### Backend
- **Go 1.24+**: Modern Go with generics and latest features
- **Gorilla Mux**: HTTP routing
- **PostgreSQL**: Primary database with advanced features
- **SQLC**: Type-safe SQL query generation
- **pgx/v5**: High-performance PostgreSQL driver

### Frontend
- **Go Templates**: Server-side rendering
- **HTMX**: Dynamic web interactions
- **Tailwind CSS**: Utility-first CSS framework
- **Alpine.js**: Minimal JavaScript framework (via CDN)

### DevOps & Tooling
- **Docker**: Containerized development environment
- **Go Migrate**: Database migration management
- **Make**: Build automation
- **golangci-lint**: Code quality and linting

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- Docker and Docker Compose
- Make

### 1. Clone and Setup
```bash
git clone <repository-url>
cd finance
make setup
```

### 2. Start Database
```bash
docker-compose up -d db
```

### 3. Create Database and Run Migrations
```bash
# Create the finance database
docker exec -it db createdb -U postgres finance

# Run migrations
export DATABASE_HOST=localhost
export DATABASE_USER=postgres  
export DATABASE_PASSWORD=postgres
export DATABASE_NAME=finance
make migration/up
```

### 4. Build Applications
```bash
# Build both API and Web services
go build -o bin/service cmd/service/main.go
go build -o bin/web cmd/web/main.go
```

### 5. Run Applications
```bash
# Terminal 1: Start API service (port 8000)
export DATABASE_HOST=localhost DATABASE_USER=postgres DATABASE_PASSWORD=postgres DATABASE_NAME=finance
./bin/service

# Terminal 2: Start Web frontend (port 8080)  
export DATABASE_HOST=localhost DATABASE_USER=postgres DATABASE_PASSWORD=postgres DATABASE_NAME=finance
./bin/web
```

### 6. Access Applications
- **Web Interface**: http://localhost:8080
- **REST API**: http://localhost:8000
- **API Health Check**: http://localhost:8000/health

## 📖 API Documentation

### Accounts
- `GET /api/v1/accounts` - List all accounts
- `POST /api/v1/accounts` - Create account
- `GET /api/v1/accounts/{id}` - Get account by ID
- `PUT /api/v1/accounts/{id}` - Update account
- `DELETE /api/v1/accounts/{id}` - Delete account

### Categories  
- `GET /api/v1/categories` - List all categories
- `POST /api/v1/categories` - Create category
- `GET /api/v1/categories/{id}` - Get category by ID
- `PUT /api/v1/categories/{id}` - Update category
- `DELETE /api/v1/categories/{id}` - Delete category

### Transactions
- `GET /api/v1/transactions` - List all transactions
- `POST /api/v1/transactions` - Create transaction
- `GET /api/v1/transactions/{id}` - Get transaction by ID
- `PUT /api/v1/transactions/{id}` - Update transaction
- `DELETE /api/v1/transactions/{id}` - Delete transaction

### Balances
- `GET /api/v1/balances` - Get all account balances
- `GET /api/v1/balances/{account_id}` - Get specific account balance
- `POST /api/v1/balances/refresh` - Refresh all balances

## 🎨 Web Interface Features

### Dashboard
- Account balance overview
- Recent transaction summary  
- Quick action buttons
- Financial health indicators

### Account Management
- Add/edit/delete accounts
- Multiple account types
- Real-time balance display

### Category Organization
- Income vs expense categorization
- Color-coded categories
- Default categories included

### Transaction Tracking
- Add transactions with validation
- Status tracking (pending/cleared/cancelled)
- Account and category selection
- Date and amount validation

## 🏗️ Project Structure

```
├── cmd/                           # Application entry points
│   ├── service/main.go           # REST API service
│   └── web/main.go               # Web frontend service
├── domain/                       # Business logic layer
│   ├── entities/                 # Domain entities
│   └── finance/                  # Finance-specific use cases
├── internal/                     # Private application code
│   ├── api/                      # REST API handlers
│   ├── config/                   # Configuration management
│   ├── repository/pg/            # PostgreSQL repositories
│   └── web/                      # Web frontend handlers
│       ├── handlers.go           # HTTP handlers
│       ├── templates/            # HTML templates
│       └── static/               # Static assets
├── docker-compose.yaml           # Development environment
├── Makefile                      # Development commands
└── sqlc.yaml                     # SQLC configuration
```

## 🔧 Development Commands

```bash
# Setup development environment
make setup

# Database operations
make migration/create          # Create new migration
make migration/up             # Apply migrations
make migration/down           # Rollback migrations

# Code generation
make generate                 # Generate all code
make sqlc-generate           # Generate SQLC code only

# Building
make compile                 # Build service binary
go build -o bin/web cmd/web/main.go  # Build web frontend

# Testing
make test                    # Run tests
make test-full              # Run all tests including integration
make coverage               # Generate coverage report

# Code quality
make lint                   # Run linters
make gosec                  # Security analysis
```

## 💡 Key Design Decisions

### Why HTMX?
- **Server-side rendering**: Leverages Go's template system
- **Minimal JavaScript**: Reduces complexity and bundle size
- **Progressive enhancement**: Works without JavaScript
- **Real-time updates**: Seamless partial page updates

### Why PostgreSQL?
- **ACID compliance**: Ensures data consistency for financial data
- **Advanced features**: Triggers, constraints, and complex queries
- **Performance**: Excellent performance for transactional workloads
- **Type safety**: Strong typing matches Go's type system

### Why Clean Architecture?
- **Testability**: Easy to test business logic in isolation
- **Maintainability**: Clear separation of concerns
- **Flexibility**: Easy to change external dependencies
- **Scalability**: Supports application growth

## 🤝 AI Development Process

This project demonstrates effective AI-human collaboration:

### What AI Excelled At:
- **Rapid prototyping**: Quickly generated working code structures
- **Best practices**: Applied Go idioms and patterns consistently
- **Documentation**: Generated comprehensive code documentation
- **Testing**: Created thorough test coverage
- **Integration**: Seamlessly connected all components

### Human Guidance Provided:
- **Product vision**: Defined the overall application requirements
- **Architecture decisions**: Chose technology stack and patterns
- **User experience**: Guided UI/UX design decisions
- **Business logic**: Validated financial calculation logic

### Lessons Learned:
1. **Clear initial requirements** lead to better AI assistance
2. **Iterative development** works well with AI collaboration
3. **Code review** remains important even with AI-generated code
4. **Domain expertise** should guide AI implementation

## 📝 Next Steps

Potential enhancements for this application:

- **Reports & Analytics**: Monthly/yearly financial reports
- **Budget Management**: Budget tracking and alerts  
- **Data Import**: CSV/OFX import capabilities
- **Multi-currency**: Support for multiple currencies
- **Mobile App**: React Native or Flutter mobile application
- **Authentication**: User accounts and security
- **Backup/Export**: Data backup and export features

## 🙏 Acknowledgments

- **Claude Sonnet 4**: AI assistant that helped build this application
- **Go Community**: For excellent tooling and libraries
- **HTMX Team**: For making web development simple again
- **Tailwind CSS**: For excellent utility-first CSS framework

---

**Built with ❤️ by Human + AI collaboration**

This project serves as a reference for:
- Modern Go web application development
- Clean architecture implementation
- AI-assisted software development
- HTMX and server-side rendering patterns
- Financial application design patterns

