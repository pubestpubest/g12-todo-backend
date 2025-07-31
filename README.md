# G12 Todo Backend

A production-ready Todo application backend built with Go using Clean Architecture principles.

## 🚀 Features

- Clean Architecture implementation
- RESTful API with Gin framework
- PostgreSQL database integration with GORM
- Environment-based configuration
- Structured logging with Logrus
- CORS middleware
- Health check endpoint
- Versioned API routes
- Feature-based organization
- Domain-driven design
- Standardized error handling
- Docker containerization
- Task management (CRUD operations)
- Pagination support

## 📁 Project Structure

```
.
├── configs/              # Configuration files
├── configs.example/      # Example configuration files
├── constant/            # Global constants
├── database/            # Database connection and migrations
├── domain/              # Core business logic and entities
├── feature/             # Feature modules
│   └── task/            # Task feature
│       ├── delivery/    # HTTP handlers
│       ├── repository/  # Data access layer
│       └── usecase/     # Business logic
├── middlewares/         # HTTP middlewares
├── models/              # Data models
├── request/             # Request DTOs
├── response/            # Response DTOs
├── routes/              # API route definitions
├── utils/               # Utility functions
│   └── error.go         # Error handling utilities
├── main.go              # Application entry point
├── Dockerfile           # Docker configuration
├── docker-compose.yaml  # Docker Compose configuration
└── go.mod               # Go module definition
```

## 🛠️ Prerequisites

- Go 1.24 or higher
- PostgreSQL
- Docker (optional, for containerized deployment)

## 🏁 Getting Started

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/pubestpubest/g12-todo-backend.git
   cd g12-todo-backend
   ```

2. Copy example configuration:
   ```bash
   cp -r configs.example/* configs/
   ```

3. Update the configuration files in `configs/` with your environment-specific settings.

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Run the application:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:3000` (or the port specified in your configuration).

### Docker Deployment

1. Build and run with Docker Compose (requires env file):
   ```bash
   docker-compose --env-file ./configs/.env up --force-recreate --build
   ```

2. Or build and run manually (requires environment variables):
   ```bash
   docker build -t g12-todo-backend .
   docker run --env-file ./configs/.env -p 3000:3000 g12-todo-backend
   ```

**Note:** All Docker commands require the environment file (`./configs/.env`) to be present and properly configured with the required variables.

## 🔧 Configuration

The application uses environment variables for configuration. Create a `.env` file in the `configs/` directory with the following variables:

```env
PROJECT_NAME=g12-todo
RUN_ENV=development # development || production
DEPLOY_ENV=local # local || container

DATABASE_HOST=localhost # g12-todo-db || localhost
DATABASE_PORT=5432
DATABASE_USERNAME=appuser
DATABASE_PASSWORD=apppass
DATABASE_NAME=todo_app

BACKEND_PORT=3000
```

### Environment Variables

- `PROJECT_NAME`: Project name for Docker services
- `RUN_ENV`: Application environment (development/production)
- `DEPLOY_ENV`: Deployment environment (local/container)
- `DATABASE_HOST`: PostgreSQL host address
- `DATABASE_PORT`: PostgreSQL port
- `DATABASE_USERNAME`: Database username
- `DATABASE_PASSWORD`: Database password
- `DATABASE_NAME`: Database name
- `BACKEND_PORT`: Backend server port

## 📚 API Documentation

### Health Check
- `GET /healthz` - Health check endpoint

### Task Management API (v1)

Base URL: `/v1/tasks`

#### Get Task List
- **GET** `/v1/tasks`
- **Query Parameters:**
  - `page` (optional): Page number (default: 1)
  - `limit` (optional): Items per page (default: 10)
- **Response:**
  ```json
  {
    "status": "SUCCESS",
    "message": "List tasks successfully",
    "data": [
      {
        "taskId": 1,
        "title": "Complete project",
        "description": "Finish the todo backend project",
        "status": false,
        "createdAt": "2024-01-01T00:00:00Z",
        "updateAt": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "totalPages": 1
    }
  }
  ```

#### Get Task by ID
- **GET** `/v1/tasks/{id}`
- **Response:**
  ```json
  {
    "status": "SUCCESS",
    "message": "Task retrieved successfully",
    "data": {
      "taskId": 1,
      "title": "Complete project",
      "description": "Finish the todo backend project",
      "status": false,
      "createdAt": "2024-01-01T00:00:00Z",
      "updateAt": "2024-01-01T00:00:00Z"
    }
  }
  ```

#### Create Task
- **POST** `/v1/tasks`
- **Request Body:**
  ```json
  {
    "title": "New task",
    "description": "Task description",
    "status": false
  }
  ```
- **Response:**
  ```json
  {
    "status": "SUCCESS",
    "message": "Task created successfully",
    "data": {
      "taskId": 1,
      "title": "New task",
      "description": "Task description",
      "status": false,
      "createdAt": "2024-01-01T00:00:00Z",
      "updateAt": "2024-01-01T00:00:00Z"
    }
  }
  ```

#### Update Task
- **PUT** `/v1/tasks/{id}`
- **Request Body:**
  ```json
  {
    "title": "Updated task",
    "description": "Updated description",
    "status": true
  }
  ```
- **Response:**
  ```json
  {
    "status": "SUCCESS",
    "message": "Task updated successfully",
    "data": {
      "taskId": 1,
      "title": "Updated task",
      "description": "Updated description",
      "status": true,
      "createdAt": "2024-01-01T00:00:00Z",
      "updateAt": "2024-01-01T00:00:00Z"
    }
  }
  ```

#### Delete Task
- **DELETE** `/v1/tasks/{id}`
- **Response:**
  ```json
  {
    "status": "SUCCESS",
    "message": "Task deleted successfully",
    "data": null
  }
  ```

## 🔍 Logging

The application uses Logrus for structured logging with the following features:

### Log Levels
- `Info` - General operational entries about what's happening inside the application
- `Warn` - Warning messages that don't necessarily affect the application's operation
- `Error` - Error events that might still allow the application to continue running
- `Fatal` - Critical errors that force the application to exit

### Log Format
```go
log.SetFormatter(&log.TextFormatter{
    ForceColors:   true,
    FullTimestamp: true,
})
```

### Example Usage
```go
// Info level logging
log.Info("Application started successfully")

// Warning level logging
log.Warn("Resource usage is high")

// Error level logging
log.Error("Failed to connect to database")

// Fatal level logging (will exit the application)
log.Fatal("Critical error occurred")
```

### Environment-based Logging
- Development mode: Colorized output with full timestamps
- Production mode: Plain text output with essential information only

## 🚨 Error Handling

The application implements a standardized error handling approach using the `utils/error.go` utility.

### Error Wrapping
Errors are wrapped with context using the `errors.Wrap` function from the `github.com/pkg/errors` package:

```go
err = errors.Wrap(err, "[TaskRepository.GetTaskByID]: Error getting task")
```

### Standard Error Formatting
The `StandardError` function in `utils/error.go` is used to extract clean error messages:

```go
func StandardError(err error) string {
    errorMessages := strings.Split(err.Error(), "]: ")
    return errorMessages[len(errorMessages)-1]
}
```

### Usage in Handlers
In HTTP handlers, errors are processed and returned in a standardized format:

```go
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
    task, err := h.taskUsecase.GetTaskByID(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": utils.StandardError(err),
        })
        return
    }
    c.JSON(http.StatusOK, task)
}
```

This approach provides:
- Consistent error formatting across the application
- Clean error messages for API responses
- Detailed error context in logs
- Easy error tracking through the call stack

## 🏗️ Architecture

This project follows Clean Architecture principles with clear separation of concerns:

### Domain Layer
- Contains business logic and interfaces
- Independent of external frameworks
- Defines repository interfaces
- Contains core domain models

### Application Layer
- Implements use cases
- Orchestrates domain objects
- Contains business rules specific to application
- Depends only on the domain layer

### Interface Adapters
- HTTP handlers in feature modules
- Repository implementations
- DTOs in request and response packages
- Converts data between layers

### Infrastructure Layer
- Database configuration
- External services integration
- Framework configuration
- Middleware implementation

## 📚 Libraries and Dependencies

### Core Libraries
- [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP web framework
- [Logrus](https://github.com/sirupsen/logrus) - Structured logging
- [godotenv](https://github.com/joho/godotenv) - Environment variable management
- [GORM](https://gorm.io/) - ORM for database operations
- [PostgreSQL](https://www.postgresql.org/) - Database system

### Development Tools
- [Go](https://golang.org/) - Programming language
- [Docker](https://www.docker.com/) - Containerization platform

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin) - For providing a robust and efficient web framework
- [Logrus](https://github.com/sirupsen/logrus) - For structured logging capabilities
- [godotenv](https://github.com/joho/godotenv) - For environment variable management
- [GORM](https://gorm.io/) - For database operations and migrations
- [PostgreSQL](https://www.postgresql.org/) - For reliable database management
