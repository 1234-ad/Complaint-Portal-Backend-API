# Setup and Installation Guide

## Prerequisites

### Installing Go

1. **Download Go:**
   - Visit https://golang.org/dl/
   - Download the Windows installer for Go 1.21 or later
   - Run the installer and follow the setup wizard

2. **Verify Installation:**
   ```powershell
   go version
   ```
   You should see output like: `go version go1.21.x windows/amd64`

3. **Set Environment Variables (if not done automatically):**
   - Add Go bin directory to PATH: `C:\Program Files\Go\bin`
   - Set GOPATH (optional): `C:\Users\%USERNAME%\go`

## Quick Start

### Option 1: Using Go directly

1. **Navigate to project directory:**
   ```powershell
   cd "C:\Users\ASUS\Downloads\1_Aditya Rastogi"
   ```

2. **Initialize Go modules:**
   ```powershell
   go mod tidy
   ```

3. **Run the application:**
   ```powershell
   go run main.go
   ```

4. **Test the API:**
   Open another terminal window and run:
   ```powershell
   go run client_demo.go
   ```

### Option 2: Using PowerShell test script

1. **Navigate to project directory:**
   ```powershell
   cd "C:\Users\ASUS\Downloads\1_Aditya Rastogi"
   ```

2. **Run the test script:**
   ```powershell
   .\test_api.ps1
   ```

### Option 3: Build executable

1. **Build the application:**
   ```powershell
   go build -o complaint-portal.exe main.go
   ```

2. **Run the executable:**
   ```powershell
   .\complaint-portal.exe
   ```

## Testing with curl (if available)

### Register a user:
```bash
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d "{\"name\": \"John Doe\", \"email\": \"john@example.com\"}"
```

### Login:
```bash
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d "{\"secret_code\": \"YOUR_SECRET_CODE\"}"
```

## Testing with PowerShell (Invoke-RestMethod)

### Register a user:
```powershell
$body = @{
    name = "John Doe"
    email = "john@example.com"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/register" -Method POST -Body $body -ContentType "application/json"
```

### Login:
```powershell
$body = @{
    secret_code = "YOUR_SECRET_CODE"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/login" -Method POST -Body $body -ContentType "application/json"
```

## Default Admin Account

When you start the server, a default admin account is created:
- **Secret Code:** `ADMIN_SECRET_123`
- **Name:** System Administrator
- **Email:** admin@complaintportal.com

Use this admin account to test administrative functions.

## Troubleshooting

### Go not found error:
- Ensure Go is properly installed
- Check that Go bin directory is in your PATH
- Restart your terminal/PowerShell after installation

### Port already in use:
- If port 8080 is busy, modify the port in `main.go`:
  ```go
  port := ":8081"  // Change to different port
  ```

### Permission errors:
- Run PowerShell as Administrator if needed
- Check Windows Defender/Antivirus settings

## Project Structure

```
complaint-portal/
├── main.go              # Main application with all API handlers
├── main_test.go         # Unit tests for the API
├── client_demo.go       # Demo client showing API usage
├── test_api.ps1         # PowerShell test script
├── go.mod               # Go module file
├── Makefile             # Build automation (for systems with make)
├── README.md            # Main documentation
└── SETUP.md             # This setup guide
```

## Next Steps

1. Install Go if not already installed
2. Run `go mod tidy` to initialize modules
3. Start the server with `go run main.go`
4. Test with the demo client: `go run client_demo.go`
5. Explore the API endpoints using the documentation in README.md

## Development

### Running tests:
```powershell
# Start server in background
Start-Job -ScriptBlock { go run main.go }
Start-Sleep 3

# Run tests
go test -v

# Stop background job
Get-Job | Stop-Job
Get-Job | Remove-Job
```

### Code formatting:
```powershell
go fmt ./...
```

### Building for production:
```powershell
go build -ldflags="-s -w" -o complaint-portal.exe main.go
```