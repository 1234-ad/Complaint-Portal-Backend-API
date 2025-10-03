package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// User represents a user in the system
type User struct {
	ID         int         `json:"id"`
	SecretCode string      `json:"secret_code"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Complaints []Complaint `json:"complaints"`
	IsAdmin    bool        `json:"is_admin"`
}

// Complaint represents a complaint in the system
type Complaint struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Summary      string `json:"summary"`
	Rating       int    `json:"rating"`
	UserID       int    `json:"user_id"`
	UserName     string `json:"user_name,omitempty"`
	IsResolved   bool   `json:"is_resolved"`
	CreatedAt    string `json:"created_at"`
	ResolvedAt   string `json:"resolved_at,omitempty"`
}

// Request/Response structures
type LoginRequest struct {
	SecretCode string `json:"secret_code"`
}

type RegisterRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SubmitComplaintRequest struct {
	SecretCode string `json:"secret_code"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Rating     int    `json:"rating"`
}

type ViewComplaintRequest struct {
	SecretCode  string `json:"secret_code"`
	ComplaintID int    `json:"complaint_id"`
}

type ResolveComplaintRequest struct {
	SecretCode  string `json:"secret_code"`
	ComplaintID int    `json:"complaint_id"`
}

type GetComplaintsRequest struct {
	SecretCode string `json:"secret_code"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Global storage with mutex for concurrency safety
type Storage struct {
	users      map[int]*User
	complaints map[int]*Complaint
	userIDGen  int
	compIDGen  int
	mutex      sync.RWMutex
}

var storage = &Storage{
	users:      make(map[int]*User),
	complaints: make(map[int]*Complaint),
	userIDGen:  0,
	compIDGen:  0,
}

// Helper functions
func generateSecretCode() string {
	return fmt.Sprintf("SEC_%d_%d", time.Now().Unix(), storage.userIDGen+1)
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func findUserBySecretCode(secretCode string) *User {
	storage.mutex.RLock()
	defer storage.mutex.RUnlock()
	
	for _, user := range storage.users {
		if user.SecretCode == secretCode {
			return user
		}
	}
	return nil
}

func findUserByEmail(email string) *User {
	storage.mutex.RLock()
	defer storage.mutex.RUnlock()
	
	for _, user := range storage.users {
		if user.Email == email {
			return user
		}
	}
	return nil
}

func respondWithJSON(w http.ResponseWriter, statusCode int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

// API Handlers

// /register - Create a new user
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate input
	if strings.TrimSpace(req.Name) == "" {
		respondWithError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// Check if email already exists
	if findUserByEmail(req.Email) != nil {
		respondWithError(w, http.StatusConflict, "User with this email already exists")
		return
	}

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	storage.userIDGen++
	newUser := &User{
		ID:         storage.userIDGen,
		SecretCode: generateSecretCode(),
		Name:       strings.TrimSpace(req.Name),
		Email:      strings.TrimSpace(req.Email),
		Complaints: []Complaint{},
		IsAdmin:    false, // Default users are not admin
	}

	storage.users[newUser.ID] = newUser

	respondWithJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    newUser,
	})
}

// /login - User login with secret code
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if strings.TrimSpace(req.SecretCode) == "" {
		respondWithError(w, http.StatusBadRequest, "Secret code is required")
		return
	}

	user := findUserBySecretCode(req.SecretCode)
	if user == nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid secret code")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    user,
	})
}

// /submitComplaint - Submit a new complaint
func submitComplaintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req SubmitComplaintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate input
	if strings.TrimSpace(req.SecretCode) == "" {
		respondWithError(w, http.StatusBadRequest, "Secret code is required")
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if strings.TrimSpace(req.Summary) == "" {
		respondWithError(w, http.StatusBadRequest, "Summary is required")
		return
	}
	if req.Rating < 1 || req.Rating > 10 {
		respondWithError(w, http.StatusBadRequest, "Rating must be between 1 and 10")
		return
	}

	user := findUserBySecretCode(req.SecretCode)
	if user == nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid secret code")
		return
	}

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	storage.compIDGen++
	newComplaint := &Complaint{
		ID:         storage.compIDGen,
		Title:      strings.TrimSpace(req.Title),
		Summary:    strings.TrimSpace(req.Summary),
		Rating:     req.Rating,
		UserID:     user.ID,
		UserName:   user.Name,
		IsResolved: false,
		CreatedAt:  getCurrentTime(),
	}

	storage.complaints[newComplaint.ID] = newComplaint
	user.Complaints = append(user.Complaints, *newComplaint)

	respondWithJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "Complaint submitted successfully",
		Data:    newComplaint,
	})
}

// /getAllComplaintsForUser - Get all complaints for a specific user
func getAllComplaintsForUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req GetComplaintsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if strings.TrimSpace(req.SecretCode) == "" {
		respondWithError(w, http.StatusBadRequest, "Secret code is required")
		return
	}

	user := findUserBySecretCode(req.SecretCode)
	if user == nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid secret code")
		return
	}

	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	var userComplaints []Complaint
	for _, complaint := range storage.complaints {
		if complaint.UserID == user.ID {
			userComplaints = append(userComplaints, *complaint)
		}
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "User complaints retrieved successfully",
		Data:    userComplaints,
	})
}

// /getAllComplaintsForAdmin - Get all complaints (admin only)
func getAllComplaintsForAdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req GetComplaintsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if strings.TrimSpace(req.SecretCode) == "" {
		respondWithError(w, http.StatusBadRequest, "Secret code is required")
		return
	}

	user := findUserBySecretCode(req.SecretCode)
	if user == nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid secret code")
		return
	}

	if !user.IsAdmin {
		respondWithError(w, http.StatusForbidden, "Access denied. Admin privileges required")
		return
	}

	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	var allComplaints []Complaint
	for _, complaint := range storage.complaints {
		allComplaints = append(allComplaints, *complaint)
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "All complaints retrieved successfully",
		Data:    allComplaints,
	})
}

// /viewComplaint - View a specific complaint
func viewComplaintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ViewComplaintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if strings.TrimSpace(req.SecretCode) == "" {
		respondWithError(w, http.StatusBadRequest, "Secret code is required")
		return
	}
	if req.ComplaintID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Valid complaint ID is required")
		return
	}

	user := findUserBySecretCode(req.SecretCode)
	if user == nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid secret code")
		return
	}

	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	complaint, exists := storage.complaints[req.ComplaintID]
	if !exists {
		respondWithError(w, http.StatusNotFound, "Complaint not found")
		return
	}

	// Check if user has permission to view this complaint
	if !user.IsAdmin && complaint.UserID != user.ID {
		respondWithError(w, http.StatusForbidden, "Access denied. You can only view your own complaints")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Complaint retrieved successfully",
		Data:    complaint,
	})
}

// /resolveComplaint - Mark a complaint as resolved (admin only)
func resolveComplaintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ResolveComplaintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if strings.TrimSpace(req.SecretCode) == "" {
		respondWithError(w, http.StatusBadRequest, "Secret code is required")
		return
	}
	if req.ComplaintID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Valid complaint ID is required")
		return
	}

	user := findUserBySecretCode(req.SecretCode)
	if user == nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid secret code")
		return
	}

	if !user.IsAdmin {
		respondWithError(w, http.StatusForbidden, "Access denied. Admin privileges required")
		return
	}

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	complaint, exists := storage.complaints[req.ComplaintID]
	if !exists {
		respondWithError(w, http.StatusNotFound, "Complaint not found")
		return
	}

	if complaint.IsResolved {
		respondWithError(w, http.StatusBadRequest, "Complaint is already resolved")
		return
	}

	complaint.IsResolved = true
	complaint.ResolvedAt = getCurrentTime()

	// Update the complaint in user's list as well
	if userOwner, exists := storage.users[complaint.UserID]; exists {
		for i := range userOwner.Complaints {
			if userOwner.Complaints[i].ID == complaint.ID {
				userOwner.Complaints[i].IsResolved = true
				userOwner.Complaints[i].ResolvedAt = complaint.ResolvedAt
				break
			}
		}
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Complaint resolved successfully",
		Data:    complaint,
	})
}

// Create default admin user
func createDefaultAdmin() {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	storage.userIDGen++
	adminUser := &User{
		ID:         storage.userIDGen,
		SecretCode: "ADMIN_SECRET_123",
		Name:       "System Administrator",
		Email:      "admin@complaintportal.com",
		Complaints: []Complaint{},
		IsAdmin:    true,
	}

	storage.users[adminUser.ID] = adminUser
	fmt.Println("Default admin created with secret code:", adminUser.SecretCode)
}

func main() {
	// Create default admin user
	createDefaultAdmin()

	// Setup routes
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/submitComplaint", submitComplaintHandler)
	http.HandleFunc("/getAllComplaintsForUser", getAllComplaintsForUserHandler)
	http.HandleFunc("/getAllComplaintsForAdmin", getAllComplaintsForAdminHandler)
	http.HandleFunc("/viewComplaint", viewComplaintHandler)
	http.HandleFunc("/resolveComplaint", resolveComplaintHandler)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Message: "Complaint Portal API is running",
		})
	})

	port := ":8080"
	fmt.Printf("Complaint Portal API server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST /register")
	fmt.Println("  POST /login")
	fmt.Println("  POST /submitComplaint")
	fmt.Println("  POST /getAllComplaintsForUser")
	fmt.Println("  POST /getAllComplaintsForAdmin")
	fmt.Println("  POST /viewComplaint")
	fmt.Println("  POST /resolveComplaint")
	fmt.Println("  GET  /health")
	fmt.Println("\nDefault Admin Secret Code: ADMIN_SECRET_123")

	log.Fatal(http.ListenAndServe(port, nil))
}