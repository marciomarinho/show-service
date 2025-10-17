package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/marciomarinho/show-service/internal/domain"
	serviceMocks "github.com/marciomarinho/show-service/internal/service/mocks"
)

func TestShowHTTPHandler_PostShows(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*serviceMocks.MockShowService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful show creation",
			requestBody: `{
				"payload": [
					{
						"slug": "show/testshow",
						"title": "Test Show",
						"seasons": [{"slug": "show/testshow/season/1"}]
					}
				],
				"skip": 0,
				"take": 10,
				"totalRecords": 1
			}`,
			mockSetup: func(m *serviceMocks.MockShowService) {
				m.EXPECT().Create(mock.AnythingOfType("domain.Request")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "Shows created successfully",
			},
		},
		{
			name: "invalid JSON",
			requestBody: `{
				"payload": [
					"invalid": "json"
				]
			}`,
			mockSetup: func(m *serviceMocks.MockShowService) {
				// No service call expected due to JSON parsing error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid character ':' after array element",
			},
		},
		{
			name: "request validation error - empty payload",
			requestBody: `{
				"payload": [],
				"skip": 0,
				"take": 10,
				"totalRecords": 0
			}`,
			mockSetup: func(m *serviceMocks.MockShowService) {
				// No service call expected due to validation error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error":   "Validation failed",
				"details": "payload must contain between 1 and 1000 items",
			},
		},
		{
			name: "request validation error - invalid take value",
			requestBody: `{
				"payload": [
					{
						"slug": "show/testshow",
						"title": "Test Show"
					}
				],
				"skip": 0,
				"take": 0,
				"totalRecords": 1
			}`,
			mockSetup: func(m *serviceMocks.MockShowService) {
				// No service call expected due to validation error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error":   "Validation failed",
				"details": "take must be between 1 and 100",
			},
		},
		{
			name: "show validation error - empty title",
			requestBody: `{
				"payload": [
					{
						"slug": "show/testshow",
						"title": ""
					}
				],
				"skip": 0,
				"take": 10,
				"totalRecords": 1
			}`,
			mockSetup: func(m *serviceMocks.MockShowService) {
				// No service call expected due to show validation error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error":   "Validation failed",
				"details": "payload[0]: title is required",
			},
		},
		{
			name: "show validation error - invalid slug format",
			requestBody: `{
				"payload": [
					{
						"slug": "invalid-slug",
						"title": "Test Show"
					}
				],
				"skip": 0,
				"take": 10,
				"totalRecords": 1
			}`,
			mockSetup: func(m *serviceMocks.MockShowService) {
				// No service call expected due to show validation error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error":   "Validation failed",
				"details": "payload[0]: slug: must be in a valid format.",
			},
		},
		{
			name: "service error",
			requestBody: `{
				"payload": [
					{
						"slug": "show/testshow",
						"title": "Test Show"
					}
				],
				"skip": 0,
				"take": 10,
				"totalRecords": 1
			}`,
			mockSetup: func(m *serviceMocks.MockShowService) {
				m.EXPECT().Create(mock.AnythingOfType("domain.Request")).Return(errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "database connection failed",
			},
		},
		{
			name: "multiple shows with one invalid",
			requestBody: `{
				"payload": [
					{
						"slug": "show/validshow",
						"title": "Valid Show"
					},
					{
						"slug": "show/invalidshow",
						"title": ""
					}
				],
				"skip": 0,
				"take": 10,
				"totalRecords": 2
			}`,
			mockSetup: func(m *serviceMocks.MockShowService) {
				// No service call expected due to second show validation error
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error":   "Validation failed",
				"details": "payload[1]: title is required",
			},
		},
		{
			name:        "large payload - edge case",
			requestBody: createLargePayload(1000),
			mockSetup: func(m *serviceMocks.MockShowService) {
				// Large payload passes validation and should reach service
				m.EXPECT().Create(mock.AnythingOfType("domain.Request")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "Shows created successfully",
			},
		},
		{
			name:        "payload too large - edge case",
			requestBody: createLargePayload(1001),
			mockSetup: func(m *serviceMocks.MockShowService) {
				// No service call expected due to payload size validation
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error":   "Validation failed",
				"details": "payload must contain between 1 and 1000 items",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSvc := serviceMocks.NewMockShowService(t)
			tt.mockSetup(mockSvc)

			handler := NewShowHandler(mockSvc)

			// Create test request
			req, _ := http.NewRequest(http.MethodPost, "/shows", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create test response recorder
			w := httptest.NewRecorder()

			// Create gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.PostShows(c)

			// Assert
			require.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			if tt.expectedBody != nil {
				for key, expectedValue := range tt.expectedBody {
					require.Equal(t, expectedValue, responseBody[key])
				}
			}
		})
	}
}

func TestShowHTTPHandler_GetShows(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(*serviceMocks.MockShowService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful shows retrieval",
			mockSetup: func(m *serviceMocks.MockShowService) {
				m.EXPECT().List().Return(&domain.Response{
					Response: []domain.ShowResponse{
						{
							Slug:  "show/testshow1",
							Title: "Test Show 1",
							Image: "http://example.com/image1.jpg",
						},
						{
							Slug:  "show/testshow2",
							Title: "Test Show 2",
							Image: "http://example.com/image2.jpg",
						},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"response": []interface{}{
					map[string]interface{}{
						"slug":  "show/testshow1",
						"title": "Test Show 1",
						"image": "http://example.com/image1.jpg",
					},
					map[string]interface{}{
						"slug":  "show/testshow2",
						"title": "Test Show 2",
						"image": "http://example.com/image2.jpg",
					},
				},
			},
		},
		{
			name: "empty shows list",
			mockSetup: func(m *serviceMocks.MockShowService) {
				m.EXPECT().List().Return(&domain.Response{
					Response: []domain.ShowResponse{},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"response": []interface{}{},
			},
		},
		{
			name: "service error",
			mockSetup: func(m *serviceMocks.MockShowService) {
				m.EXPECT().List().Return((*domain.Response)(nil), errors.New("database query failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "database query failed",
			},
		},
		{
			name: "service returns nil response",
			mockSetup: func(m *serviceMocks.MockShowService) {
				m.EXPECT().List().Return((*domain.Response)(nil), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil, // Gin returns null for nil JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSvc := serviceMocks.NewMockShowService(t)
			tt.mockSetup(mockSvc)

			handler := NewShowHandler(mockSvc)

			// Create test request
			req, _ := http.NewRequest(http.MethodGet, "/shows", nil)

			// Create test response recorder
			w := httptest.NewRecorder()

			// Create gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Execute
			handler.GetShows(c)

			// Assert
			require.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)
				require.Equal(t, tt.expectedBody, responseBody)
			} else {
				// For nil response, expect JSON null
				require.Equal(t, "null", w.Body.String())
			}
		})
	}
}

// Helper function to create large payloads for testing
func createLargePayload(size int) string {
	shows := make([]map[string]interface{}, size)
	for i := 0; i < size; i++ {
		shows[i] = map[string]interface{}{
			"slug":  "show/testshow" + strconv.Itoa(i),
			"title": "Test Show " + strconv.Itoa(i),
		}
	}

	payload := map[string]interface{}{
		"payload":      shows,
		"skip":         0,
		"take":         10,
		"totalRecords": size,
	}

	jsonBytes, _ := json.Marshal(payload)
	return string(jsonBytes)
}
