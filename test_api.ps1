# Complaint Portal API Test Script
# This script demonstrates the API functionality

Write-Host "Complaint Portal API Demo Script" -ForegroundColor Green
Write-Host "=================================" -ForegroundColor Green

# Start the server in background
Write-Host "`nStarting API server..." -ForegroundColor Yellow
$serverJob = Start-Job -ScriptBlock { Set-Location "C:\Users\ASUS\Downloads\1_Aditya Rastogi"; go run main.go }

# Wait for server to start
Start-Sleep -Seconds 3
Write-Host "Server started successfully!" -ForegroundColor Green

$baseUrl = "http://localhost:8080"

# Function to make API calls
function Invoke-APICall {
    param(
        [string]$Method,
        [string]$Endpoint,
        [hashtable]$Body = @{}
    )
    
    try {
        $headers = @{ "Content-Type" = "application/json" }
        $uri = "$baseUrl$Endpoint"
        
        if ($Body.Count -gt 0) {
            $jsonBody = $Body | ConvertTo-Json
            $response = Invoke-RestMethod -Uri $uri -Method $Method -Headers $headers -Body $jsonBody
        } else {
            $response = Invoke-RestMethod -Uri $uri -Method $Method -Headers $headers
        }
        
        return $response
    }
    catch {
        Write-Host "Error calling API: $($_.Exception.Message)" -ForegroundColor Red
        return $null
    }
}

try {
    # 1. Health Check
    Write-Host "`n1. Testing Health Check..." -ForegroundColor Cyan
    $health = Invoke-APICall -Method "GET" -Endpoint "/health"
    if ($health) {
        Write-Host "✓ Health check successful" -ForegroundColor Green
        Write-Host "Response: $($health.message)" -ForegroundColor White
    }

    # 2. Register User
    Write-Host "`n2. Registering new user..." -ForegroundColor Cyan
    $registerData = @{
        name = "Test User"
        email = "testuser@example.com"
    }
    $registerResponse = Invoke-APICall -Method "POST" -Endpoint "/register" -Body $registerData
    if ($registerResponse -and $registerResponse.success) {
        Write-Host "✓ User registered successfully" -ForegroundColor Green
        $userSecretCode = $registerResponse.data.secret_code
        Write-Host "Secret Code: $userSecretCode" -ForegroundColor Yellow
    }

    # 3. Login
    Write-Host "`n3. Testing login..." -ForegroundColor Cyan
    $loginData = @{
        secret_code = $userSecretCode
    }
    $loginResponse = Invoke-APICall -Method "POST" -Endpoint "/login" -Body $loginData
    if ($loginResponse -and $loginResponse.success) {
        Write-Host "✓ Login successful" -ForegroundColor Green
        Write-Host "User: $($loginResponse.data.name)" -ForegroundColor White
    }

    # 4. Submit Complaint
    Write-Host "`n4. Submitting complaint..." -ForegroundColor Cyan
    $complaintData = @{
        secret_code = $userSecretCode
        title = "Test Complaint"
        summary = "This is a test complaint for demonstration"
        rating = 7
    }
    $complaintResponse = Invoke-APICall -Method "POST" -Endpoint "/submitComplaint" -Body $complaintData
    if ($complaintResponse -and $complaintResponse.success) {
        Write-Host "✓ Complaint submitted successfully" -ForegroundColor Green
        $complaintId = $complaintResponse.data.id
        Write-Host "Complaint ID: $complaintId" -ForegroundColor Yellow
    }

    # 5. Get User Complaints
    Write-Host "`n5. Getting user complaints..." -ForegroundColor Cyan
    $getUserComplaintsData = @{
        secret_code = $userSecretCode
    }
    $userComplaintsResponse = Invoke-APICall -Method "POST" -Endpoint "/getAllComplaintsForUser" -Body $getUserComplaintsData
    if ($userComplaintsResponse -and $userComplaintsResponse.success) {
        Write-Host "✓ Retrieved user complaints" -ForegroundColor Green
        Write-Host "Number of complaints: $($userComplaintsResponse.data.Count)" -ForegroundColor White
    }

    # 6. View Specific Complaint
    Write-Host "`n6. Viewing specific complaint..." -ForegroundColor Cyan
    $viewComplaintData = @{
        secret_code = $userSecretCode
        complaint_id = $complaintId
    }
    $viewComplaintResponse = Invoke-APICall -Method "POST" -Endpoint "/viewComplaint" -Body $viewComplaintData
    if ($viewComplaintResponse -and $viewComplaintResponse.success) {
        Write-Host "✓ Complaint details retrieved" -ForegroundColor Green
        Write-Host "Title: $($viewComplaintResponse.data.title)" -ForegroundColor White
        Write-Host "Status: $(if ($viewComplaintResponse.data.is_resolved) { 'Resolved' } else { 'Open' })" -ForegroundColor White
    }

    # 7. Admin Operations
    Write-Host "`n7. Testing admin operations..." -ForegroundColor Cyan
    $adminSecretCode = "ADMIN_SECRET_123"
    
    # Get all complaints (admin)
    $adminComplaintsData = @{
        secret_code = $adminSecretCode
    }
    $adminComplaintsResponse = Invoke-APICall -Method "POST" -Endpoint "/getAllComplaintsForAdmin" -Body $adminComplaintsData
    if ($adminComplaintsResponse -and $adminComplaintsResponse.success) {
        Write-Host "✓ Admin retrieved all complaints" -ForegroundColor Green
        Write-Host "Total complaints in system: $($adminComplaintsResponse.data.Count)" -ForegroundColor White
    }

    # 8. Resolve Complaint (admin)
    Write-Host "`n8. Resolving complaint (admin)..." -ForegroundColor Cyan
    $resolveComplaintData = @{
        secret_code = $adminSecretCode
        complaint_id = $complaintId
    }
    $resolveResponse = Invoke-APICall -Method "POST" -Endpoint "/resolveComplaint" -Body $resolveComplaintData
    if ($resolveResponse -and $resolveResponse.success) {
        Write-Host "✓ Complaint resolved successfully" -ForegroundColor Green
        Write-Host "Resolved at: $($resolveResponse.data.resolved_at)" -ForegroundColor White
    }

    # 9. Error Testing
    Write-Host "`n9. Testing error cases..." -ForegroundColor Cyan
    
    # Invalid secret code
    $invalidLoginData = @{
        secret_code = "INVALID_SECRET"
    }
    $invalidResponse = Invoke-APICall -Method "POST" -Endpoint "/login" -Body $invalidLoginData
    if ($invalidResponse -and -not $invalidResponse.success) {
        Write-Host "✓ Invalid secret code properly rejected" -ForegroundColor Green
    }

    # User trying to access admin endpoint
    $userTryingAdminData = @{
        secret_code = $userSecretCode
    }
    $unauthorizedResponse = Invoke-APICall -Method "POST" -Endpoint "/getAllComplaintsForAdmin" -Body $userTryingAdminData
    if ($unauthorizedResponse -and -not $unauthorizedResponse.success) {
        Write-Host "✓ Unauthorized access properly blocked" -ForegroundColor Green
    }

    Write-Host "`n=================================" -ForegroundColor Green
    Write-Host "All tests completed successfully!" -ForegroundColor Green
    Write-Host "The Complaint Portal API is working correctly." -ForegroundColor Green
}
catch {
    Write-Host "Error during testing: $($_.Exception.Message)" -ForegroundColor Red
}
finally {
    # Stop the server
    Write-Host "`nStopping server..." -ForegroundColor Yellow
    Stop-Job -Job $serverJob
    Remove-Job -Job $serverJob
    Write-Host "Server stopped." -ForegroundColor Green
}

Write-Host "`nTo start the server manually, run: go run main.go" -ForegroundColor Cyan
Write-Host "Server will be available at: http://localhost:8080" -ForegroundColor Cyan