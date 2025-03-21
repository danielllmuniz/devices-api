package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/danielllmuniz/devices-api/internal/services"
	"github.com/danielllmuniz/devices-api/internal/store/mockstore"
	"github.com/go-chi/chi/v5"
)

func TestHandleCreateDevice(t *testing.T) {
	mock := mockstore.NewMockDeviceStore()
	api := Api{
		DeviceService: *services.NewDeviceService(mock),
	}

	tests := []struct {
		name         string
		payload      string
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "Happy path in use",
			payload:      `{"name": "Device A", "brand": "BrandX", "state": "in-use"}`,
			wantStatus:   http.StatusCreated,
			wantResponse: ``,
		},
		{
			name:         "Happy path available",
			payload:      `{"name": "Device B", "brand": "BrandY", "state": "available"}`,
			wantStatus:   http.StatusCreated,
			wantResponse: ``,
		},
		{
			name:         "Happy path inactive",
			payload:      `{"name": "Device C", "brand": "BrandZ", "state": "inactive"}`,
			wantStatus:   http.StatusCreated,
			wantResponse: ``,
		},
		{
			name:         "Empty fields",
			payload:      `{"name": "", "brand": "", "state": ""}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"brand":"Brand is required","name":"Name is required","state":"State is required"}`,
		},
		{
			name:         "Empty payload",
			payload:      "",
			wantStatus:   http.StatusBadRequest,
			wantResponse: `{"error":"invalid request"}`,
		},
		{
			name:         "Invalid JSON",
			payload:      `{"name":`,
			wantStatus:   http.StatusBadRequest,
			wantResponse: `{"error":"invalid request"}`,
		},
		{
			name:         "Missing name",
			payload:      `{"brand":"BrandX","state":"in-use"}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"name":"Name is required"}`,
		},
		{
			name:         "Missing brand",
			payload:      `{"name":"Device A","state":"in-use"}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"brand":"Brand is required"}`,
		},
		{
			name:         "Missing state",
			payload:      `{"name":"Device A","brand":"BrandX"}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"state":"State is required"}`,
		},
		{
			name:         "Invalid state value",
			payload:      `{"name":"Device A","brand":"BrandX","state":"unknown-state"}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"state":"State must be 'available', 'in-use' or 'inactive'"}`,
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

			if tt.wantStatus != http.StatusCreated && !strings.Contains(rec.Body.String(), tt.wantResponse) {
				t.Errorf("Expected response to contain '%s', got '%s'", tt.wantResponse, rec.Body)
			}

		})
	}
}

func TestHandleGetDevice(t *testing.T) {
	mock := mockstore.NewMockDeviceStore()
	api := Api{
		DeviceService: *services.NewDeviceService(mock),
	}
	ctx := context.Background()

	mock.CreateDevice(ctx, "Device A", "BrandX", "in-use")
	mock.CreateDevice(ctx, "Device A", "BrandX", "in-use")

	tests := []struct {
		name         string
		deviceID     string
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "Happy path",
			deviceID:     "1",
			wantStatus:   http.StatusOK,
			wantResponse: ``,
		},
		{
			name:         "Invalid device ID",
			deviceID:     "invalid",
			wantStatus:   http.StatusBadRequest,
			wantResponse: `{"error":"invalid device id"}`,
		},
		{
			name:         "Device not found",
			deviceID:     "3",
			wantStatus:   http.StatusNotFound,
			wantResponse: `{"error":"device not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := chi.NewRouter()
			handler.Get("/api/v1/devices/{device_id}", api.handleGetDevice)
			req := httptest.NewRequest("GET", "/api/v1/devices/"+tt.deviceID, nil)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantStatus != http.StatusOK && !strings.Contains(rec.Body.String(), tt.wantResponse) {
				t.Errorf("Expected response to contain '%s', got '%s'", tt.wantResponse, rec.Body)
			}

		})
	}

}

func TestHandleGetAllDevices(t *testing.T) {
	mock := mockstore.NewMockDeviceStore()
	api := Api{
		DeviceService: *services.NewDeviceService(mock),
	}
	ctx := context.Background()

	mock.CreateDevice(ctx, "Device A", "BrandX", "in-use")
	mock.CreateDevice(ctx, "Device B", "BrandY", "available")
	mock.CreateDevice(ctx, "Device C", "BrandZ", "inactive")

	tests := []struct {
		name         string
		paramBrand   string
		paramState   string
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "Happy path",
			paramBrand:   "",
			paramState:   "",
			wantStatus:   http.StatusOK,
			wantResponse: ``,
		},
		{
			name:         "Filter by brand",
			paramBrand:   "BrandX",
			paramState:   "",
			wantStatus:   http.StatusOK,
			wantResponse: ``,
		},
		{
			name:         "Filter by state",
			paramBrand:   "",
			paramState:   "in-use",
			wantStatus:   http.StatusOK,
			wantResponse: `in-use`,
		},
		{
			name:         "Filter by brand and state",
			paramBrand:   "BrandX",
			paramState:   "in-use",
			wantStatus:   http.StatusOK,
			wantResponse: ``,
		},
		{
			name:         "Invalid brand",
			paramBrand:   "BrandInvalid",
			paramState:   "",
			wantStatus:   http.StatusOK,
			wantResponse: ``,
		},
		{
			name:         "Invalid state",
			paramBrand:   "",
			paramState:   "invalid",
			wantStatus:   http.StatusOK,
			wantResponse: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := url.Values{}
			query.Add("brand", tt.paramBrand)
			query.Add("state", tt.paramState)

			req := httptest.NewRequest("GET", "/api/v1/devices?", nil)

			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(api.handleGetAllDevices)
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantStatus != http.StatusOK && !strings.Contains(rec.Body.String(), tt.wantResponse) {
				t.Errorf("Expected response to contain '%s', got '%s'", tt.wantResponse, rec.Body)
			}
		})
	}
}

func TestHandleUpdateDevice(t *testing.T) {
	mock := mockstore.NewMockDeviceStore()
	api := Api{
		DeviceService: *services.NewDeviceService(mock),
	}
	ctx := context.Background()

	mock.CreateDevice(ctx, "Device A", "BrandX", "available")
	mock.CreateDevice(ctx, "Device B", "BrandY", "in-use")
	mock.CreateDevice(ctx, "Device C", "BrandZ", "inactive")

	tests := []struct {
		name         string
		deviceID     string
		payload      string
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "Happy path",
			deviceID:     "1",
			payload:      `{"name": "Device A updated", "brand": "BrandX updated", "state": "in-use"}`,
			wantStatus:   http.StatusOK,
			wantResponse: ``,
		},
		{
			name:         "Empty payload",
			deviceID:     "1",
			payload:      "",
			wantStatus:   http.StatusBadRequest,
			wantResponse: `{"error":"invalid request"}`,
		},
		{
			name:         "Invalid JSON",
			deviceID:     "1",
			payload:      `{"name":`,
			wantStatus:   http.StatusBadRequest,
			wantResponse: `{"error":"invalid request"}`,
		},
		{
			name:         "Invalid state value",
			deviceID:     "1",
			payload:      `{"name": "Device A updated", "brand": "BrandX updated", "state": "unknown-state"}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"state":"State must be 'available', 'in-use' or 'inactive'"}`,
		},
		{
			name:         "Empty fields",
			deviceID:     "1",
			payload:      `{"name": "", "brand": "", "state": ""}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"brand":"Brand is required","name":"Name is required","state":"State is required"}`,
		},
		{
			name:         "Invalid device ID",
			deviceID:     "invalid",
			payload:      `{"name": "Device A updated", "brand": "BrandX updated", "state": "in-use"}`,
			wantStatus:   http.StatusBadRequest,
			wantResponse: `{"error":"invalid device id"}`,
		},
		{
			name:         "Device not found",
			deviceID:     "5",
			payload:      `{"name": "Device A updated", "brand": "BrandX updated", "state": "in-use"}`,
			wantStatus:   http.StatusNotFound,
			wantResponse: `{"error":"device not found"}`,
		},
		{
			name:         "Device in use",
			deviceID:     "2",
			payload:      `{"name": "Device A updated", "brand": "BrandX updated", "state": "available"}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"error":"device is in use, cannot update name or brand"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := chi.NewRouter()
			handler.Put("/api/v1/devices/{device_id}", api.handleUpdateDevice)
			req := httptest.NewRequest("PUT", "/api/v1/devices/"+tt.deviceID, strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantStatus != http.StatusOK && !strings.Contains(rec.Body.String(), tt.wantResponse) {
				t.Errorf("Expected response to contain '%s', got '%s'", tt.wantResponse, rec.Body)
			}

		})
	}
}

func TestHandlePatchDevice(t *testing.T) {
	mock := mockstore.NewMockDeviceStore()
	api := Api{
		DeviceService: *services.NewDeviceService(mock),
	}
	ctx := context.Background()

	mock.CreateDevice(ctx, "Device A", "BrandX", "available")
	mock.CreateDevice(ctx, "Device B", "BrandY", "in-use")
	mock.CreateDevice(ctx, "Device C", "BrandZ", "inactive")

	tests := []struct {
		name         string
		deviceID     string
		payload      string
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "Happy path",
			deviceID:     "1",
			payload:      `{"name": "Device A updated", "brand": "BrandX updated", "state": "in-use"}`,
			wantStatus:   http.StatusOK,
			wantResponse: ``,
		},
		{
			name:         "Empty payload",
			deviceID:     "1",
			payload:      "",
			wantStatus:   http.StatusBadRequest,
			wantResponse: `{"error":"invalid request"}`,
		},
		{
			name:         "Invalid JSON",
			deviceID:     "1",
			payload:      `{"name":`,
			wantStatus:   http.StatusBadRequest,
			wantResponse: `{"error":"invalid request"}`,
		},
		{
			name:         "Invalid state value",
			deviceID:     "1",
			payload:      `{"name": "Device A updated", "brand": "BrandX updated", "state": "unknown-state"}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"state":"State must be 'available', 'in-use' or 'inactive"}`,
		},
		{
			name:         "Empty fields",
			deviceID:     "1",
			payload:      `{"name": "", "brand": "", "state": ""}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"brand":"At least one field must be informed","name":"At least one field must be informed","state":"At least one field must be informed"}`,
		},
		{
			name:         "Invalid device ID",
			deviceID:     "invalid",
			payload:      `{"name": "Device A updated", "brand": "BrandX updated", "state": "in-use"}`,
			wantStatus:   http.StatusBadRequest,
			wantResponse: `{"error":"invalid device id"}`,
		},
		{
			name:         "Device not found",
			deviceID:     "4",
			payload:      `{"name": "Device A updated", "brand": "BrandX updated", "state": "in-use"}`,
			wantStatus:   http.StatusNotFound,
			wantResponse: `{"error":"device not found"}`,
		},
		{
			name:         "Device in use",
			deviceID:     "2",
			payload:      `{"name": "Device A updated", "brand": "BrandX updated", "state": "available"}`,
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"error":"device is in use, cannot update name or brand"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := chi.NewRouter()
			handler.Patch("/api/v1/devices/{device_id}", api.handlePatchDevice)
			req := httptest.NewRequest("PATCH", "/api/v1/devices/"+tt.deviceID, strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantStatus != http.StatusOK && !strings.Contains(rec.Body.String(), tt.wantResponse) {
				t.Errorf("Expected response to contain '%s', got '%s'", tt.wantResponse, rec.Body)
			}

		})
	}
}

func TestHandleDeleteDevice(t *testing.T) {
	mock := mockstore.NewMockDeviceStore()
	api := Api{
		DeviceService: *services.NewDeviceService(mock),
	}
	ctx := context.Background()

	mock.CreateDevice(ctx, "Device A", "BrandX", "available")
	mock.CreateDevice(ctx, "Device B", "BrandY", "in-use")
	mock.CreateDevice(ctx, "Device C", "BrandZ", "inactive")

	tests := []struct {
		name         string
		deviceID     string
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "Happy path",
			deviceID:     "1",
			wantStatus:   http.StatusOK,
			wantResponse: `{"device_id":1,"message":"device deleted successfully"}`,
		},
		{
			name:         "Invalid device ID",
			deviceID:     "invalid",
			wantStatus:   http.StatusBadRequest,
			wantResponse: `{"error":"invalid device id"}`,
		},
		{
			name:         "Device not found",
			deviceID:     "4",
			wantStatus:   http.StatusNotFound,
			wantResponse: `{"error":"device not found"}`,
		},
		{
			name:         "Device in use",
			deviceID:     "2",
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"error":"device is in use, cannot update name or brand"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handle := chi.NewRouter()
			handle.Delete("/api/v1/devices/{device_id}", api.handleDeleteDevice)
			req := httptest.NewRequest("DELETE", "/api/v1/devices/"+tt.deviceID, nil)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handle.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Expected status code %d, got %d", tt.wantStatus, rec.Code)
			}

			if !strings.Contains(rec.Body.String(), tt.wantResponse) {
				t.Errorf("Expected response to contain '%s', got '%s'", tt.wantResponse, rec.Body)
			}
		})
	}
}
