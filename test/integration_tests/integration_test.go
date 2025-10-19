package integration_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestIntegration runs end-to-end tests for the application
func TestIntegration(t *testing.T) {
	// Wait for the application to start (assumes Makefile has started DynamoDB and app)
	time.Sleep(5 * time.Second) // Adjust based on startup time

	// Run sub-tests
	t.Run("Health Check", testHealthCheck)
	t.Run("Shows Endpoints", testShowsIntegration)
}

func testHealthCheck(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/v1/health")
	if err != nil {
		t.Fatalf("Failed to reach health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expected := `{"message":"ok"}`
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

func testShowsIntegration(t *testing.T) {
	// Read test files
	requestBody, err := os.ReadFile("data/complete_request.json")
	if err != nil {
		t.Fatalf("Failed to read complete_request.json: %v", err)
	}

	// Unmarshal request to get expected payload
	var requestData map[string]interface{}
	if err := json.Unmarshal(requestBody, &requestData); err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	payload, ok := requestData["payload"].([]interface{})
	if !ok {
		t.Fatal("Request payload not found or invalid")
	}

	// Create expected response with only the fields returned by the API
	var expectedResponseList []interface{}
	for _, show := range payload {
		showMap, ok := show.(map[string]interface{})
		if !ok {
			continue
		}
		imageMap, ok := showMap["image"].(map[string]interface{})
		if !ok {
			continue
		}
		responseShow := map[string]interface{}{
			"image": imageMap["showImage"],
			"slug":  showMap["slug"],
			"title": showMap["title"],
		}
		expectedResponseList = append(expectedResponseList, responseShow)
	}
	expectedResponseMap := map[string]interface{}{
		"response": expectedResponseList,
	}

	tests := []struct {
		name           string
		method         string
		body           io.Reader
		expectedStatus int
		expectError    bool
		validate       func(t *testing.T, resp *http.Response)
	}{
		{
			name:           "GET empty shows",
			method:         "GET",
			body:           nil,
			expectedStatus: 200,
			expectError:    false,
			validate: func(t *testing.T, resp *http.Response) {
				body, _ := io.ReadAll(resp.Body)
				expected := `{"response":[]}`
				if string(body) != expected {
					t.Errorf("Expected empty response %q, got %q", expected, string(body))
				}
			},
		},
		{
			name:           "POST complete request",
			method:         "POST",
			body:           bytes.NewReader(requestBody),
			expectedStatus: 201,
			expectError:    false,
			validate: func(t *testing.T, resp *http.Response) {
				// Just check status for success
			},
		},
		{
			name:           "GET after POST",
			method:         "GET",
			body:           nil,
			expectedStatus: 200,
			expectError:    false,
			validate: func(t *testing.T, resp *http.Response) {
				body, _ := io.ReadAll(resp.Body)
				var actualResponse map[string]interface{}
				if err := json.Unmarshal(body, &actualResponse); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				// Check top-level "response" array
				responseList, ok := actualResponse["response"].([]interface{})
				if !ok {
					t.Fatal("Response does not have 'response' array")
				}

				// Create expected map keyed by slug
				expectedMap := make(map[string]interface{})
				if expectedList, ok := expectedResponseMap["response"].([]interface{}); ok {
					for _, show := range expectedList {
						showMap, ok := show.(map[string]interface{})
						if !ok {
							continue
						}
						if slug, ok := showMap["slug"].(string); ok {
							expectedMap[slug] = show
						}
					}
				}

				// Iterate over response items and compare with expected by slug
				for _, show := range responseList {
					showMap, ok := show.(map[string]interface{})
					if !ok {
						t.Error("Invalid show item in response")
						continue
					}
					slug, ok := showMap["slug"].(string)
					if !ok {
						t.Error("Show item missing slug")
						continue
					}
					expectedShow, exists := expectedMap[slug]
					if !exists {
						t.Errorf("Unexpected show with slug: %s", slug)
						continue
					}
					if !reflect.DeepEqual(show, expectedShow) {
						t.Errorf("Mismatch for slug %s.\nExpected: %+v\nGot: %+v", slug, expectedShow, show)
					}
				}
			},
		},
		{
			name:           "POST invalid JSON",
			method:         "POST",
			body:           strings.NewReader("invalid json"),
			expectedStatus: 400,
			expectError:    true,
			validate: func(t *testing.T, resp *http.Response) {
				// Expect error response
			},
		},
		{
			name:           "POST duplicate request",
			method:         "POST",
			body:           bytes.NewReader(requestBody),
			expectedStatus: 500, // Expect error for duplicates
			expectError:    true,
			validate: func(t *testing.T, resp *http.Response) {
				// Expect error response
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			if tt.method == "GET" {
				resp, err = http.Get("http://localhost:8080/v1/shows")
			} else if tt.method == "POST" {
				resp, err = http.Post("http://localhost:8080/v1/shows", "application/json", tt.body)
			}

			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, resp.StatusCode, string(body))
			}

			if tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}
