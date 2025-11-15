ser Service API Examples

## New Clean Handler Pattern

The user-service now uses a cleaner handler pattern with `ginext` wrapper that provides:

- **Structured Error Handling**: Errors are handled with proper HTTP status codes at the point where they occur
- **Clean Response Format**: Consistent JSON response structure across all endpoints
- **Better Logging**: Contextual logging with handler names
- **Panic Recovery**: Automatic panic recovery with proper error responses

## API Endpoints

### Authentication Endpoints

#### 1. User Signup

```bash
POST /api/v1/auth/signup
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "username": "john_doe",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "role": 1
}
```

**Success Response (201):**

```json
{
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "username": "john_doe",
      "first_name": "John",
      "last_name": "Doe",
      "role": 1,
      "status": "active"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  },
  "message": "User registered successfully",
  "success": true
}
```

**Error Response (409 - Conflict):**

```json
{
  "success": false,
  "error_message": "email already exists",
  "error_code": "HTTP_409"
}
```

#### 2. User Signin

```bash
POST /api/v1/auth/signin
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### 3. OAuth2 Firebase Signin

```bash
POST /api/v1/auth/oauth2/signin
Content-Type: application/json

{
  "provider": "firebase",
  "id_token": "firebase_id_token_here"
}
```

#### 4. Token Verification (for Gateway)

```bash
POST /api/v1/auth/verify-token
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response:**

```json
{
  "data": {
    "valid": true,
    "user_id": "1",
    "role": 1
  },
  "message": "Token verified successfully",
  "success": true
}
```

#### 5. Token Refresh

```bash
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 6. User Signout

```bash
POST /api/v1/auth/signout
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Features

### ✅ Bitwise Role System

- **Passenger Role**: `1` (binary: 001)
- **Admin Role**: `2` (binary: 010)
- **Operator Role**: `4` (binary: 100)
- **Support Role**: `8` (binary: 1000)

### ✅ Firebase OAuth2 Integration

- Firebase Admin SDK v4.14.1
- ID token verification
- Automatic user creation for OAuth2 users

### ✅ JWT Token Management

- Access tokens (short-lived)
- Refresh tokens (long-lived)
- Token verification endpoint for API Gateway

### ✅ Clean Error Handling

- Structured error responses
- Proper HTTP status codes
- Detailed error messages
- Panic recovery

### ✅ Enhanced Logging

- Contextual logging with handler names
- Request tracing
- Error details in logs

## Router Pattern

The new router uses `ginext.WrapHandler()` for clean error handling:

```go
auth.POST("/signup", ginext.WrapHandler(config.AuthHandler.Signup))
auth.POST("/signin", ginext.WrapHandler(config.AuthHandler.Signin))
auth.POST("/oauth2/signin", ginext.WrapHandler(config.AuthHandler.OAuth2Signin))
```

## Handler Pattern

Handlers now return structured responses and errors:

```go
func (h *AuthHandler) Signup(r *ginext.Request) (*ginext.Response, error) {
    req := model.SignupRequest{}
    r.MustBind(&req)

    if err := ginext.ValidateRequest(&req); err != nil {
        return nil, err
    }

    authResp, err := h.authService.Signup(r.Context(), &req)
    if err != nil {
        if err.Error() == "email already exists" {
            return nil, ginext.NewConflictError(err.Error())
        }
        return nil, ginext.NewInternalServerError("Registration failed")
    }

    return ginext.NewCreatedResponse(authResp, "User registered successfully"), nil
}
```

This pattern provides:

- ✅ **Explicit error handling** with proper status codes
- ✅ **Clean separation** between business logic and HTTP concerns
- ✅ **Consistent response format** across all endpoints
- ✅ **Better testability** with structured responses
- ✅ **Automatic panic recovery** and error formatting
