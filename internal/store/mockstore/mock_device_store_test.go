package mockstore

import (
	"context"
	"testing"

	"github.com/danielllmuniz/devices-api/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestMockDeviceStore(t *testing.T) {
	ctx := context.Background()
	mockStore := NewMockDeviceStore()

	t.Run("CreateDevice", func(t *testing.T) {
		device, err := mockStore.CreateDevice(ctx, "Device A", "BrandX", "available")
		assert.NoError(t, err)
		assert.Equal(t, int32(1), device.ID)
		assert.Equal(t, "Device A", device.Name)
		assert.Equal(t, "BrandX", device.Brand)
		assert.Equal(t, store.DeviceState("available"), device.State)
	})

	t.Run("GetDeviceByID", func(t *testing.T) {
		device, err := mockStore.GetDeviceByID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, "Device A", device.Name)
		assert.Equal(t, "BrandX", device.Brand)
		assert.Equal(t, store.DeviceState("available"), device.State)
	})

	t.Run("UpdateDevice", func(t *testing.T) {
		updated, err := mockStore.UpdateDevice(ctx, 1, "Device A+", "BrandX", "inactive")
		assert.NoError(t, err)
		assert.Equal(t, "Device A+", updated.Name)
		assert.Equal(t, store.DeviceState("inactive"), updated.State)
	})

	t.Run("PatchDevice", func(t *testing.T) {
		patched, err := mockStore.PatchDevice(ctx, 1, "", "BrandY", "in-use")
		assert.NoError(t, err)
		assert.Equal(t, "BrandY", patched.Brand)
		assert.Equal(t, "Device A+", patched.Name)
		assert.Equal(t, store.DeviceState("in-use"), patched.State)
	})

	t.Run("GetAllDevices", func(t *testing.T) {
		_, _ = mockStore.CreateDevice(ctx, "Device B", "BrandX", "inactive")
		devices, err := mockStore.GetAllDevices(ctx)
		assert.NoError(t, err)
		assert.Len(t, devices, 2)
	})

	t.Run("GetDevicesByBrand", func(t *testing.T) {
		devices, err := mockStore.GetDevicesByBrand(ctx, "BrandX")
		assert.NoError(t, err)
		assert.Len(t, devices, 1)
	})

	t.Run("GetDevicesByState", func(t *testing.T) {
		devices, err := mockStore.GetDevicesByState(ctx, "inactive")
		assert.NoError(t, err)
		assert.Len(t, devices, 1)
	})

	t.Run("DeleteDevice", func(t *testing.T) {
		deletedID, err := mockStore.DeleteDevice(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int32(1), deletedID)

		_, err = mockStore.GetDeviceByID(ctx, 1)
		assert.Error(t, err)
	})
}
