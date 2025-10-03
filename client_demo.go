package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const apiBaseURL = "http://localhost:8080"

// API client helper
func callAPI(method, endpoint string, payload interface{}) (map[string]interface{}, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, apiBaseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	fmt.Println("Complaint Portal API Client Demo")
	fmt.Println("=================================")

	// Wait for server to be ready
	fmt.Println("Waiting for server to start...")
	time.Sleep(2 * time.Second)

	// 1. Health check
	fmt.Println("\n1. Health Check:")
	result, err := callAPI("GET", "/health", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// 2. Register a new user
	fmt.Println("\n2. Registering a new user:")
	registerPayload := map[string]interface{}{
		"name":  "Alice Johnson",
		"email": "alice.johnson@example.com",
	}
	result, err = callAPI("POST", "/register", registerPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// Extract user secret code
	userData := result["data"].(map[string]interface{})
	userSecretCode := userData["secret_code"].(string)
	fmt.Printf("User Secret Code: %s\n", userSecretCode)

	// 3. Login with secret code
	fmt.Println("\n3. Logging in:")
	loginPayload := map[string]interface{}{
		"secret_code": userSecretCode,
	}
	result, err = callAPI("POST", "/login", loginPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// 4. Submit a complaint
	fmt.Println("\n4. Submitting a complaint:")
	complaintPayload := map[string]interface{}{
		"secret_code": userSecretCode,
		"title":       "Broken Air Conditioning",
		"summary":     "The AC in conference room B is not working. It's too hot for meetings.",
		"rating":      8,
	}
	result, err = callAPI("POST", "/submitComplaint", complaintPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// Extract complaint ID
	complaintData := result["data"].(map[string]interface{})
	complaintID := complaintData["id"].(float64)
	fmt.Printf("Complaint ID: %.0f\n", complaintID)

	// 5. Submit another complaint
	fmt.Println("\n5. Submitting another complaint:")
	complaint2Payload := map[string]interface{}{
		"secret_code": userSecretCode,
		"title":       "Parking Issue",
		"summary":     "Not enough parking spaces for employees. Need more spots.",
		"rating":      6,
	}
	result, err = callAPI("POST", "/submitComplaint", complaint2Payload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// 6. Get all complaints for user
	fmt.Println("\n6. Getting all complaints for user:")
	getUserComplaintsPayload := map[string]interface{}{
		"secret_code": userSecretCode,
	}
	result, err = callAPI("POST", "/getAllComplaintsForUser", getUserComplaintsPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// 7. View specific complaint
	fmt.Println("\n7. Viewing specific complaint:")
	viewComplaintPayload := map[string]interface{}{
		"secret_code":  userSecretCode,
		"complaint_id": int(complaintID),
	}
	result, err = callAPI("POST", "/viewComplaint", viewComplaintPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// 8. Admin operations - Get all complaints
	fmt.Println("\n8. Admin: Getting all complaints:")
	adminSecretCode := "ADMIN_SECRET_123"
	getAllComplaintsPayload := map[string]interface{}{
		"secret_code": adminSecretCode,
	}
	result, err = callAPI("POST", "/getAllComplaintsForAdmin", getAllComplaintsPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// 9. Admin resolves complaint
	fmt.Println("\n9. Admin: Resolving complaint:")
	resolveComplaintPayload := map[string]interface{}{
		"secret_code":  adminSecretCode,
		"complaint_id": int(complaintID),
	}
	result, err = callAPI("POST", "/resolveComplaint", resolveComplaintPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// 10. View resolved complaint
	fmt.Println("\n10. Viewing resolved complaint:")
	result, err = callAPI("POST", "/viewComplaint", viewComplaintPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// 11. Error case - Try to access admin endpoint with user credentials
	fmt.Println("\n11. Error case - User trying to access admin endpoint:")
	result, err = callAPI("POST", "/getAllComplaintsForAdmin", getUserComplaintsPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	// 12. Error case - Invalid secret code
	fmt.Println("\n12. Error case - Invalid secret code:")
	invalidLoginPayload := map[string]interface{}{
		"secret_code": "INVALID_SECRET_CODE",
	}
	result, err = callAPI("POST", "/login", invalidLoginPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Response: %v\n", result)

	fmt.Println("\n=================================")
	fmt.Println("Demo completed successfully!")
}