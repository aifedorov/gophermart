# Gophermart - Loyalty System

An educational project demonstrating real-world Go development patterns and libraries used in production services, implementing clean architecture principles and modern software engineering practices.

## 🚀 Overview

Gophermart is a microservices-based loyalty system that enables users to:
- Register and authenticate securely
- Submit order numbers for loyalty point accrual
- Track order processing status in real-time
- Manage loyalty balance and transactions
- Withdraw points for order payments

## 🛠 Technology Stack

### Core Technologies
- **Go 1.24.1** - Primary programming language
- **Chi Router** - High-performance HTTP router for REST API
- **PostgreSQL** - Primary database with ACID compliance
- **JWT** - Stateless authentication and authorization
- **Docker & Docker Compose** - Containerization and orchestration

### Architecture & Design Patterns
- **Clean Architecture** - Layered separation of concerns (handler, service, repository)
- **Domain-Driven Design** - Business logic encapsulated in domain models
- **Repository Pattern** - Data access abstraction layer
- **Dependency Injection** - Loose coupling through interface-based design
- **Middleware Pattern** - Cross-cutting concerns via request interceptors

### Development Tools
- **SQLC** - Type-safe SQL code generation
- **Golang Migrate** - Database schema versioning
- **Testify** - Comprehensive testing framework
- **Uber Mock** - Interface mocking for unit tests
- **Zap** - High-performance structured logging

### DevOps & CI/CD
- **GitHub Actions** - Automated CI/CD pipeline
- **Docker** - Application containerization
- **Makefile** - Build and test automation

## 🏗 System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │    │   Gophermart    │    │   Accrual       │
│                 │◄──►│   (Main API)    │◄──►│   (External)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │   Database      │
                       └─────────────────┘
```

### Core Components
- **HTTP Handlers** - REST API request processing
- **Domain Services** - Business logic implementation
- **Repository Layer** - Data persistence abstraction
- **External Client** - Third-party service integration
- **Background Workers** - Asynchronous order processing

## 📋 API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/user/register` | User registration | ❌ No |
| POST | `/api/user/login` | User authentication | ❌ No |
| POST | `/api/user/orders` | Order number submission | ✅ Yes |
| GET | `/api/user/orders` | User order history | ✅ Yes |
| GET | `/api/user/balance` | Current balance | ✅ Yes |
| POST | `/api/user/balance/withdraw` | Points withdrawal | ✅ Yes |
| GET | `/api/user/withdrawals` | Withdrawal history | ✅ Yes |

## 🧪 Testing Strategy

Comprehensive testing approach ensuring code quality and reliability:
- **Unit Tests** - Individual component testing with mocks
- **Integration Tests** - Cross-layer interaction validation
- **Mock Tests** - Dependency isolation and behavior verification
- **Automated Tests** - API contract testing and validation

```bash
# Run all tests
make test

# Run automated test suite
make run-autotests
```

## 🚀 Getting Started

### Local Development
```bash
# Clone repository
git clone <repository-url>
cd gophermart

# Start with Docker Compose
make docker-up

# Or run locally
make run
```

### Environment Configuration
```bash
RUN_ADDRESS=localhost:8080
DATABASE_URI=postgres://user:password@localhost:5432/gophermart
ACCRUAL_SYSTEM_ADDRESS=http://localhost:8084
SECRET_KEY=your-secret-key
```

## 📊 Engineering Excellence

### Backend Development
- ✅ RESTful API design and implementation
- ✅ PostgreSQL integration with migrations
- ✅ JWT-based authentication and authorization
- ✅ Data validation and error handling
- ✅ Structured logging and observability

### Software Architecture
- ✅ Clean Architecture principles
- ✅ Domain-Driven Design patterns
- ✅ Repository pattern implementation
- ✅ Dependency injection framework
- ✅ Middleware-based request processing

### Quality Assurance
- ✅ Unit testing with comprehensive coverage
- ✅ Integration testing for system components
- ✅ Mock-based testing for isolated validation
- ✅ Automated testing pipeline

### DevOps & Infrastructure
- ✅ Docker containerization
- ✅ CI/CD with GitHub Actions
- ✅ Build and deployment automation
- ✅ Infrastructure as Code practices

### Go Ecosystem Mastery
- ✅ Production-grade Go libraries (Chi, JWT, Zap)
- ✅ SQLC for type-safe database operations
- ✅ Go best practices and idiomatic code
- ✅ Concurrent programming with goroutines

## 📁 Project Structure

```
gophermart/
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── user/              # User domain module
│   ├── order/             # Order domain module
│   ├── client/            # HTTP client implementations
│   ├── pkg/               # Shared packages
│   └── server/            # HTTP server configuration
├── migrations/            # Database migrations
├── docker-compose.yml     # Docker environment
└── Makefile              # Build automation
```

## 🔗 Documentation

- [Technical Specification](SPECIFICATION.md)
- [API Documentation](SPECIFICATION.md#http-api-summary)
