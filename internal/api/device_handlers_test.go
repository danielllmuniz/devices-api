package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielllmuniz/devices-api/internal/services"
	"github.com/danielllmuniz/devices-api/internal/store/mockstore"
)

// Existing happy-path test
func TestHandleCreateDevice(t *testing.T) {
	mock := mockstore.NewMockDeviceStore()
	api := Api{
		DeviceService: *services.NewDeviceService(mock),
	}

	payload := `{"name": "Device A", "brand": "BrandX", "state": "in-use"}`

	req := httptest.NewRequest("POST", "/api/v1/devices", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(api.handleCreateDevice)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", rec.Code)
	}

	var resBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %s", err)
	}

	if resBody["message"] != "device created successfully" {
		t.Errorf("Expected message 'device created successfully', got '%v'", resBody["message"])
	}

	device, ok := resBody["device"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'device' key in response body")
	}

	if device["name"] != "Device A" {
		t.Errorf("Expected device name 'Device A', got '%v'", device["name"])
	}

	if device["brand"] != "BrandX" {
		t.Errorf("Expected device brand 'BrandX', got '%v'", device["brand"])
	}

	if device["state"] != "in-use" {
		t.Errorf("Expected device state 'in-use', got '%v'", device["state"])
	}
}

func TestHandleCreateDeviceValidation(t *testing.T) {
	mock := mockstore.NewMockDeviceStore()
	api := Api{
		DeviceService: *services.NewDeviceService(mock),
	}

	tests := []struct {
		name       string
		payload    string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "Empty payload",
			payload:    "",
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    "invalid payload", // or whatever your handler returns
		},
		{
			name:       "Invalid JSON",
			payload:    `{"name":`,
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    "invalid payload",
		},
		{
			name:       "Missing name",
			payload:    `{"brand":"BrandX","state":"in-use"}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    "missing field: name", // adapt to match your actual message
		},
		{
			name:       "Missing brand",
			payload:    `{"name":"Device A","state":"in-use"}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    "missing field: brand",
		},
		{
			name:       "Missing state",
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    "missing field: state",
		},
		{
			name:       "Invalid state value",
			payload:    `{"name":"Device A","brand":"BrandX","state":"unknown-state"}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    "invalid state", // or your domain-specific message
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/devices", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(api.handleCreateDevice)
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, rec.Code)
			}

			// If your endpoint always returns a JSON with an "error" key for failures:
			var resBody map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal response body: %s", err)
			}
		})
	}
}
