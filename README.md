# Complaint Portal Backend API

A comprehensive HTTP JSON Go API for managing complaints in an organization. This API allows users to submit complaints and administrators to review and resolve them.

## Features

- **User Management**: Registration and login with unique secret codes
- **Complaint Management**: Submit, view, and resolve complaints
- **Role-based Access Control**: Separate permissions for users and administrators
- **Concurrency Safe**: Thread-safe operations using mutexes
- **Error Handling**: Comprehensive error handling and validation
- **No Third-party Dependencies**: Built using only Go standard library

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

1. Clone or download the project
2. Navigate to the project directory
3. Run the application:

```bash
go run main.go
```

The server will start on port 8080.

## API Endpoints

### Base URL
```
http://localhost:8080
```

### 1. Register User
**Endpoint:** `POST /register`

**Description:** Create a new user account

**Request Body:**
```json
{
    "name": "John Doe",
    "email": "john.doe@example.com"
}
```

**Response:**
```json
{
    "success": true,
    "message": "User registered successfully",
    "data": {
        "id": 2,
        "secret_code": "SEC_1696348800_2",
        "name": "John Doe",
        "email": "john.doe@example.com",
        "complaints": [],
        "is_admin": false
    }
}
```

### 2. Login
**Endpoint:** `POST /login`

**Description:** Login with secret code

**Request Body:**
```json
{
    "secret_code": "SEC_1696348800_2"
}
```

**Response:**
```json
{
    "success": true,
    "message": "Login successful",
    "data": {
        "id": 2,
        "secret_code": "SEC_1696348800_2",
        "name": "John Doe",
        "email": "john.doe@example.com",
        "complaints": [],
        "is_admin": false
    }
}
```

### 3. Submit Complaint
**Endpoint:** `POST /submitComplaint`

**Description:** Submit a new complaint

**Request Body:**
```json
{
    "secret_code": "SEC_1696348800_2",
    "title": "Network Issue",
    "summary": "The office WiFi is not working properly",
    "rating": 8
}
```

**Response:**
```json
{
    "success": true,
    "message": "Complaint submitted successfully",
    "data": {
        "id": 1,
        "title": "Network Issue",
        "summary": "The office WiFi is not working properly",
        "rating": 8,
        "user_id": 2,
        "user_name": "John Doe",
        "is_resolved": false,
        "created_at": "2023-10-03 14:30:15"
    }
}
```

### 4. Get All Complaints for User
**Endpoint:** `POST /getAllComplaintsForUser`

**Description:** Get all complaints submitted by the authenticated user

**Request Body:**
```json
{
    "secret_code": "SEC_1696348800_2"
}
```

**Response:**
```json
{
    "success": true,
    "message": "User complaints retrieved successfully",
    "data": [
        {
            "id": 1,
            "title": "Network Issue",
            "summary": "The office WiFi is not working properly",
            "rating": 8,
            "user_id": 2,
            "user_name": "John Doe",
            "is_resolved": false,
            "created_at": "2023-10-03 14:30:15"
        }
    ]
}
```

### 5. Get All Complaints for Admin
**Endpoint:** `POST /getAllComplaintsForAdmin`

**Description:** Get all complaints in the system (Admin only)

**Request Body:**
```json
{
    "secret_code": "ADMIN_SECRET_123"
}
```

**Response:**
```json
{
    "success": true,
    "message": "All complaints retrieved successfully",
    "data": [
        {
            "id": 1,
            "title": "Network Issue",
            "summary": "The office WiFi is not working properly",
            "rating": 8,
            "user_id": 2,
            "user_name": "John Doe",
            "is_resolved": false,
            "created_at": "2023-10-03 14:30:15"
        }
    ]
}
```

### 6. View Complaint
**Endpoint:** `POST /viewComplaint`

**Description:** View details of a specific complaint

**Request Body:**
```json
{
    "secret_code": "SEC_1696348800_2",
    "complaint_id": 1
}
```

**Response:**
```json
{
    "success": true,
    "message": "Complaint retrieved successfully",
    "data": {
        "id": 1,
        "title": "Network Issue",
        "summary": "The office WiFi is not working properly",
        "rating": 8,
        "user_id": 2,
        "user_name": "John Doe",
        "is_resolved": false,
        "created_at": "2023-10-03 14:30:15"
    }
}
```

### 7. Resolve Complaint
**Endpoint:** `POST /resolveComplaint`

**Description:** Mark a complaint as resolved (Admin only)

**Request Body:**
```json
{
    "secret_code": "ADMIN_SECRET_123",
    "complaint_id": 1
}
```

**Response:**
```json
{
    "success": true,
    "message": "Complaint resolved successfully",
    "data": {
        "id": 1,
        "title": "Network Issue",
        "summary": "The office WiFi is not working properly",
        "rating": 8,
        "user_id": 2,
        "user_name": "John Doe",
        "is_resolved": true,
        "created_at": "2023-10-03 14:30:15",
        "resolved_at": "2023-10-03 16:45:30"
    }
}
```

### 8. Health Check
**Endpoint:** `GET /health`

**Description:** Check if the API is running

**Response:**
```json
{
    "success": true,
    "message": "Complaint Portal API is running"
}
```

## Data Models

### User
```json
{
    "id": 1,
    "secret_code": "SEC_1696348800_1",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "complaints": [],
    "is_admin": false
}
```

### Complaint
```json
{
    "id": 1,
    "title": "Network Issue",
    "summary": "The office WiFi is not working properly",
    "rating": 8,
    "user_id": 2,
    "user_name": "John Doe",
    "is_resolved": false,
    "created_at": "2023-10-03 14:30:15",
    "resolved_at": ""
}
```

## Error Handling

All endpoints return consistent error responses:

```json
{
    "success": false,
    "error": "Error message describing what went wrong"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created (for registration and complaint submission)
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid secret code)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found (resource doesn't exist)
- `405` - Method Not Allowed
- `409` - Conflict (duplicate email)

## Default Admin Account

The system automatically creates a default admin account:
- **Secret Code:** `ADMIN_SECRET_123`
- **Name:** System Administrator
- **Email:** admin@complaintportal.com

## Security Features

1. **Secret Code Authentication**: All operations require a valid secret code
2. **Role-based Access Control**: Admin-only operations are protected
3. **Input Validation**: All inputs are validated and sanitized
4. **Concurrency Safety**: Thread-safe operations using mutexes

## Testing with curl

### Register a new user:
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john.doe@example.com"}'
```

### Login:
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"secret_code": "YOUR_SECRET_CODE"}'
```

### Submit a complaint:
```bash
curl -X POST http://localhost:8080/submitComplaint \
  -H "Content-Type: application/json" \
  -d '{"secret_code": "YOUR_SECRET_CODE", "title": "Network Issue", "summary": "WiFi not working", "rating": 8}'
```

## Architecture

- **Concurrency Safe**: Uses `sync.RWMutex` for thread-safe operations
- **In-Memory Storage**: Data is stored in memory using Go maps
- **RESTful Design**: Follows REST principles for API design
- **No External Dependencies**: Uses only Go standard library
- **Error Handling**: Comprehensive error handling and validation
- **JSON Communication**: All requests and responses use JSON format

## Project Structure

```
complaint-portal/
├── main.go          # Main application file with all handlers
├── go.mod           # Go module file
└── README.md        # This documentation
```

## Future Enhancements

- Database integration (PostgreSQL, MySQL)
- JWT-based authentication
- Email notifications
- File attachments for complaints
- Complaint categories and priorities
- Advanced search and filtering
- Rate limiting
- Logging and monitoring
