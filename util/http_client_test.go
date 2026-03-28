package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPGetJSON(t *testing.T) {
	// Create a test response
	expectedData := map[string]string{"message": "success"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedData)
	}))
	defer server.Close()

	var result map[string]string
	err := HTTPGetJSON(server.URL, &result)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["message"] != "success" {
		t.Errorf("Expected message 'success', got %q", result["message"])
	}
}

func TestHTTPGetJSONError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	var result map[string]string
	err := HTTPGetJSON(server.URL, &result)

	if err == nil {
		t.Fatal("Expected error for non-200 status code")
	}
}

func TestHTTPGetJSONWithHeaders(t *testing.T) {
	expectedData := map[string]string{"status": "ok"}
	expectedAuth := "Bearer test-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the header was sent correctly
		auth := r.Header.Get("Authorization")
		if auth != expectedAuth {
			t.Errorf("Expected Authorization header %q, got %q", expectedAuth, auth)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedData)
	}))
	defer server.Close()

	headers := map[string]string{
		"Authorization": expectedAuth,
	}

	var result map[string]string
	err := HTTPGetJSONWithHeaders(server.URL, headers, &result)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %q", result["status"])
	}
}

func TestHTTPGetJSONWithHeadersError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	}))
	defer server.Close()

	headers := map[string]string{
		"Authorization": "Bearer invalid",
	}

	var result map[string]string
	err := HTTPGetJSONWithHeaders(server.URL, headers, &result)

	if err == nil {
		t.Fatal("Expected error for 403 status code")
	}
}
