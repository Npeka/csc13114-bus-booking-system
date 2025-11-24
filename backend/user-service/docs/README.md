# User Service API Documentation

This directory contains the auto-generated Swagger/OpenAPI documentation for the User Service API.

## Generating Documentation

### Prerequisites

Install the `swag` tool:

```bash
make install-swag
```

Or manually:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Swagger Docs

Run the following command from the `user-service` directory:

```bash
make swagger
```

Or manually:

```bash
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

This will generate:
- `docs/swagger.json` - OpenAPI 3.0 specification in JSON format
- `docs/swagger.yaml` - OpenAPI 3.0 specification in YAML format
- `docs/docs.go` - Go code for embedding Swagger UI

## Viewing Documentation

### Option 1: Swagger UI (Recommended)

The Swagger UI is automatically served when you run the user-service:

1. Start the service:
   ```bash
   make run
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:8081/swagger/index.html
   ```

### Option 2: External Tools

You can also view the generated `swagger.json` or `swagger.yaml` files using:

- [Swagger Editor](https://editor.swagger.io/) - Paste the content of `swagger.yaml`
- [Postman](https://www.postman.com/) - Import `swagger.json` to create a collection
- [Insomnia](https://insomnia.rest/) - Import `swagger.json`

## API Overview

The User Service API provides endpoints for:

### Authentication (`/api/v1/auth`)
- `POST /auth/verify` - Verify access token
- `POST /auth/firebase` - Authenticate with Firebase
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout user

### Users (`/api/v1/users`)
- `GET /users/profile` - Get current user profile
- `GET /users` - List all users (Admin only)
- `POST /users` - Create a new user (Admin only)
- `GET /users/{id}` - Get user by ID (Admin only)
- `PUT /users/{id}` - Update user information (Admin only)
- `DELETE /users/{id}` - Delete a user (Admin only)
- `PATCH /users/{id}/status` - Update user status (Admin only)

## Authentication

Most endpoints require authentication using Bearer tokens:

```
Authorization: Bearer <your-access-token>
```

You can obtain an access token by:
1. Authenticating with Firebase (`POST /auth/firebase`)
2. Using the refresh token endpoint (`POST /auth/refresh`)

## Notes

- The documentation is auto-generated from code comments
- Update the Swagger annotations in handler files when making API changes
- Regenerate docs after any API modifications using `make swagger`
