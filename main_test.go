package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8080"

// Test API client
func makeRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func TestComplaintPortalAPI(t *testing.T) {
	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Test 1: Health check
	t.Run("Health Check", func(t *testing.T) {
		resp, err := makeRequest("GET", "/health", nil)
		if err != nil {
			t.Fatalf("Health check failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	var userSecretCode string

	// Test 2: Register user
	t.Run("Register User", func(t *testing.T) {
		payload := RegisterRequest{
			Name:  "Test User",
			Email: "test@example.com",
		}

		resp, err := makeRequest("POST", "/register", payload)
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		var response APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("Expected success true, got false")
		}

		// Extract user data
		userData := response.Data.(map[string]interface{})
		userSecretCode = userData["secret_code"].(string)
		fmt.Printf("User registered with secret code: %s\n", userSecretCode)
	})

	// Test 3: Login
	t.Run("Login", func(t *testing.T) {
		payload := LoginRequest{
			SecretCode: userSecretCode,
		}

		resp, err := makeRequest("POST", "/login", payload)
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	var complaintID float64

	// Test 4: Submit complaint
	t.Run("Submit Complaint", func(t *testing.T) {
		payload := SubmitComplaintRequest{
			SecretCode: userSecretCode,
			Title:      "Test Complaint",
			Summary:    "This is a test complaint for API testing",
			Rating:     7,
		}

		resp, err := makeRequest("POST", "/submitComplaint", payload)
		if err != nil {
			t.Fatalf("Submit complaint failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		var response APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Extract complaint ID
		complaintData := response.Data.(map[string]interface{})
		complaintID = complaintData["id"].(float64)
		fmt.Printf("Complaint submitted with ID: %.0f\n", complaintID)
	})

	// Test 5: Get user complaints
	t.Run("Get User Complaints", func(t *testing.T) {
		payload := GetComplaintsRequest{
			SecretCode: userSecretCode,
		}

		resp, err := makeRequest("POST", "/getAllComplaintsForUser", payload)
		if err != nil {
			t.Fatalf("Get user complaints failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test 6: View complaint
	t.Run("View Complaint", func(t *testing.T) {
		payload := ViewComplaintRequest{
			SecretCode:  userSecretCode,
			ComplaintID: int(complaintID),
		}

		resp, err := makeRequest("POST", "/viewComplaint", payload)
		if err != nil {
			t.Fatalf("View complaint failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test 7: Admin operations
	t.Run("Admin Get All Complaints", func(t *testing.T) {
		payload := GetComplaintsRequest{
			SecretCode: "ADMIN_SECRET_123",
		}

		resp, err := makeRequest("POST", "/getAllComplaintsForAdmin", payload)
		if err != nil {
			t.Fatalf("Admin get complaints failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test 8: Resolve complaint (admin only)
	t.Run("Resolve Complaint", func(t *testing.T) {
		payload := ResolveComplaintRequest{
			SecretCode:  "ADMIN_SECRET_123",
			ComplaintID: int(complaintID),
		}

		resp, err := makeRequest("POST", "/resolveComplaint", payload)
		if err != nil {
			t.Fatalf("Resolve complaint failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test 9: Error cases
	t.Run("Invalid Secret Code", func(t *testing.T) {
		payload := LoginRequest{
			SecretCode: "INVALID_SECRET",
		}

		resp, err := makeRequest("POST", "/login", payload)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("Unauthorized Access to Admin Endpoint", func(t *testing.T) {
		payload := GetComplaintsRequest{
			SecretCode: userSecretCode, // Regular user trying to access admin endpoint
		}

		resp, err := makeRequest("POST", "/getAllComplaintsForAdmin", payload)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", resp.StatusCode)
		}
	})
}