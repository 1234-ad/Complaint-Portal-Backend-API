# Complaint Portal API Documentation

## Table of Contents
1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Data Models](#data-models)
4. [API Endpoints](#api-endpoints)
5. [Error Handling](#error-handling)
6. [Examples](#examples)
7. [Testing](#testing)

## Overview

The Complaint Portal API is a RESTful HTTP JSON API built in Go that allows users to submit complaints and administrators to manage them. The API follows best practices for security, error handling, and concurrency safety.

### Key Features
- User registration and authentication via secret codes
- Role-based access control (Users vs Administrators)
- Complaint submission and management
- Thread-safe operations
- Comprehensive error handling
- No external dependencies (pure Go standard library)

### Base URL
```
http://localhost:8080
```

## Authentication

All API operations require authentication using a unique secret code. There are two types of users:

1. **Regular Users**: Can submit complaints and view their own complaints
2. **Administrators**: Can view all complaints and resolve them

### Secret Codes
- Generated automatically during user registration
- Format: `SEC_{timestamp}_{user_id}`
- Admin default: `ADMIN_SECRET_123`

## Data Models

### User
```json
{
    "id": 1,
    "secret_code": "SEC_1696348800_1",
    "name": "John Doe",
    "email": "john@example.com",
    "complaints": [],
    "is_admin": false
}
```

**Fields:**
- `id` (int): Unique user identifier (auto-generated)
- `secret_code` (string): Unique authentication code (auto-generated)
- `name` (string): User's full name (required)
- `email` (string): User's email address (required, unique)
- `complaints` (array): List of user's complaints
- `is_admin` (boolean): Admin privilege flag

### Complaint
```json
{
    "id": 1,
    "title": "Network Issue",
    "summary": "WiFi connectivity problems in office",
    "rating": 8,
    "user_id": 2,
    "user_name": "John Doe",
    "is_resolved": false,
    "created_at": "2023-10-03 14:30:15",
    "resolved_at": ""
}
```

**Fields:**
- `id` (int): Unique complaint identifier (auto-generated)
- `title` (string): Complaint title (required)
- `summary` (string): Detailed description (required)
- `rating` (int): Severity rating 1-10 (required)
- `user_id` (int): ID of user who submitted complaint
- `user_name` (string): Name of user who submitted complaint
- `is_resolved` (boolean): Resolution status
- `created_at` (string): Timestamp when complaint was created
- `resolved_at` (string): Timestamp when complaint was resolved (if applicable)

## API Endpoints

### 1. Health Check
**GET** `/health`

Check if the API server is running.

**Response:**
```json
{
    "success": true,
    "message": "Complaint Portal API is running"
}
```

---

### 2. Register User
**POST** `/register`

Create a new user account.

**Request Body:**
```json
{
    "name": "John Doe",
    "email": "john@example.com"
}
```

**Validation:**
- `name`: Required, non-empty string
- `email`: Required, unique, non-empty string

**Response (201 Created):**
```json
{
    "success": true,
    "message": "User registered successfully",
    "data": {
        "id": 2,
        "secret_code": "SEC_1696348800_2",
        "name": "John Doe",
        "email": "john@example.com",
        "complaints": [],
        "is_admin": false
    }
}
```

**Errors:**
- `400`: Missing or invalid fields
- `409`: Email already exists

---

### 3. Login
**POST** `/login`

Authenticate user with secret code.

**Request Body:**
```json
{
    "secret_code": "SEC_1696348800_2"
}
```

**Validation:**
- `secret_code`: Required, must exist in system

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Login successful",
    "data": {
        "id": 2,
        "secret_code": "SEC_1696348800_2",
        "name": "John Doe",
        "email": "john@example.com",
        "complaints": [],
        "is_admin": false
    }
}
```

**Errors:**
- `400`: Missing secret code
- `401`: Invalid secret code

---

### 4. Submit Complaint
**POST** `/submitComplaint`

Submit a new complaint.

**Request Body:**
```json
{
    "secret_code": "SEC_1696348800_2",
    "title": "Network Issue",
    "summary": "WiFi connectivity problems in conference room",
    "rating": 8
}
```

**Validation:**
- `secret_code`: Required, must be valid
- `title`: Required, non-empty string
- `summary`: Required, non-empty string
- `rating`: Required, integer between 1-10

**Response (201 Created):**
```json
{
    "success": true,
    "message": "Complaint submitted successfully",
    "data": {
        "id": 1,
        "title": "Network Issue",
        "summary": "WiFi connectivity problems in conference room",
        "rating": 8,
        "user_id": 2,
        "user_name": "John Doe",
        "is_resolved": false,
        "created_at": "2023-10-03 14:30:15"
    }
}
```

**Errors:**
- `400`: Missing or invalid fields
- `401`: Invalid secret code

---

### 5. Get User Complaints
**POST** `/getAllComplaintsForUser`

Get all complaints submitted by the authenticated user.

**Request Body:**
```json
{
    "secret_code": "SEC_1696348800_2"
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "User complaints retrieved successfully",
    "data": [
        {
            "id": 1,
            "title": "Network Issue",
            "summary": "WiFi connectivity problems",
            "rating": 8,
            "user_id": 2,
            "user_name": "John Doe",
            "is_resolved": false,
            "created_at": "2023-10-03 14:30:15"
        }
    ]
}
```

**Errors:**
- `400`: Missing secret code
- `401`: Invalid secret code

---

### 6. Get All Complaints (Admin)
**POST** `/getAllComplaintsForAdmin`

Get all complaints in the system. **Admin only**.

**Request Body:**
```json
{
    "secret_code": "ADMIN_SECRET_123"
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "All complaints retrieved successfully",
    "data": [
        {
            "id": 1,
            "title": "Network Issue",
            "summary": "WiFi connectivity problems",
            "rating": 8,
            "user_id": 2,
            "user_name": "John Doe",
            "is_resolved": false,
            "created_at": "2023-10-03 14:30:15"
        }
    ]
}
```

**Errors:**
- `400`: Missing secret code
- `401`: Invalid secret code
- `403`: Not an administrator

---

### 7. View Complaint
**POST** `/viewComplaint`

View details of a specific complaint. Users can only view their own complaints; admins can view any complaint.

**Request Body:**
```json
{
    "secret_code": "SEC_1696348800_2",
    "complaint_id": 1
}
```

**Validation:**
- `secret_code`: Required, must be valid
- `complaint_id`: Required, must be positive integer

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Complaint retrieved successfully",
    "data": {
        "id": 1,
        "title": "Network Issue",
        "summary": "WiFi connectivity problems in conference room",
        "rating": 8,
        "user_id": 2,
        "user_name": "John Doe",
        "is_resolved": false,
        "created_at": "2023-10-03 14:30:15"
    }
}
```

**Errors:**
- `400`: Missing or invalid fields
- `401`: Invalid secret code
- `403`: No permission to view this complaint
- `404`: Complaint not found

---

### 8. Resolve Complaint
**POST** `/resolveComplaint`

Mark a complaint as resolved. **Admin only**.

**Request Body:**
```json
{
    "secret_code": "ADMIN_SECRET_123",
    "complaint_id": 1
}
```

**Validation:**
- `secret_code`: Required, must be valid admin
- `complaint_id`: Required, must exist and not already resolved

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Complaint resolved successfully",
    "data": {
        "id": 1,
        "title": "Network Issue",
        "summary": "WiFi connectivity problems in conference room",
        "rating": 8,
        "user_id": 2,
        "user_name": "John Doe",
        "is_resolved": true,
        "created_at": "2023-10-03 14:30:15",
        "resolved_at": "2023-10-03 16:45:30"
    }
}
```

**Errors:**
- `400`: Missing fields or complaint already resolved
- `401`: Invalid secret code
- `403`: Not an administrator
- `404`: Complaint not found

## Error Handling

All errors return a consistent format:

```json
{
    "success": false,
    "error": "Error message describing what went wrong"
}
```

### HTTP Status Codes

| Code | Description | When Used |
|------|-------------|-----------|
| 200 | OK | Successful GET/view operations |
| 201 | Created | Successful POST/create operations |
| 400 | Bad Request | Invalid input, validation errors |
| 401 | Unauthorized | Invalid secret code |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 405 | Method Not Allowed | Wrong HTTP method |
| 409 | Conflict | Duplicate resource (email exists) |

## Examples

### Complete User Flow

1. **Register:**
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Johnson", "email": "alice@company.com"}'
```

2. **Login (get user details):**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"secret_code": "SEC_1696348800_2"}'
```

3. **Submit complaint:**
```bash
curl -X POST http://localhost:8080/submitComplaint \
  -H "Content-Type: application/json" \
  -d '{"secret_code": "SEC_1696348800_2", "title": "Parking Issue", "summary": "Not enough parking spaces", "rating": 6}'
```

4. **View own complaints:**
```bash
curl -X POST http://localhost:8080/getAllComplaintsForUser \
  -H "Content-Type: application/json" \
  -d '{"secret_code": "SEC_1696348800_2"}'
```

### Admin Flow

1. **View all complaints:**
```bash
curl -X POST http://localhost:8080/getAllComplaintsForAdmin \
  -H "Content-Type: application/json" \
  -d '{"secret_code": "ADMIN_SECRET_123"}'
```

2. **Resolve complaint:**
```bash
curl -X POST http://localhost:8080/resolveComplaint \
  -H "Content-Type: application/json" \
  -d '{"secret_code": "ADMIN_SECRET_123", "complaint_id": 1}'
```

## Testing

### Automated Testing

Run the built-in tests:
```bash
go test -v
```

### Manual Testing

Use the provided client demo:
```bash
go run client_demo.go
```

Or use the PowerShell test script:
```powershell
.\test_api.ps1
```

### Testing with Postman

Import the following collection:
- Base URL: `http://localhost:8080`
- Set Content-Type header: `application/json`
- Use the request examples above

## Security Considerations

1. **Authentication**: All operations require valid secret codes
2. **Authorization**: Role-based access control for admin operations
3. **Input Validation**: All inputs are validated and sanitized
4. **Concurrency**: Thread-safe operations using mutexes
5. **Error Handling**: No sensitive information leaked in error messages

## Performance

- **Concurrency**: Uses `sync.RWMutex` for optimal read/write performance
- **Memory**: In-memory storage for fast access
- **Scalability**: Stateless design allows for horizontal scaling

## Limitations

1. **Persistence**: Data is stored in memory (lost on restart)
2. **Authentication**: Simple secret code system (no JWT/OAuth)
3. **File Upload**: No support for file attachments
4. **Search**: No advanced search/filtering capabilities
5. **Notifications**: No email/SMS notifications

## Future Enhancements

1. Database integration (PostgreSQL/MySQL)
2. JWT-based authentication
3. File upload support
4. Email notifications
5. Advanced search and filtering
6. Audit logging
7. Rate limiting
8. API versioning