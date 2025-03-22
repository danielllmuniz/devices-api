package services

import (
	"context"
	"testing"

	"github.com/danielllmuniz/devices-api/internal/store"
	"github.com/danielllmuniz/devices-api/internal/store/mockstore"
	"github.com/stretchr/testify/assert"
)

// setupTest consolidates the common test setup steps into a single helper.
func setupTest(t *testing.T) (context.Context, store.DeviceStore, *DeviceService) {
	t.Helper()
	ctx := context.Background()
	mock := mockstore.NewMockDeviceStore()
	svc := NewDeviceService(mock)
	return ctx, mock, svc
}

func TestCreateDevice(t *testing.T) {
	ctx, _, svc := setupTest(t)

	t.Run("It_should_be_able_to_create_a_device", func(t *testing.T) {
		device, err := svc.CreateDevice(ctx, "Device A", "BrandX", store.DeviceStateInUse)
		assert.NoError(t, err)
		assert.Equal(t, "Device A", device.Name)
		assert.Equal(t, "BrandX", device.Brand)
		assert.Equal(t, store.DeviceStateInUse, device.State)

		device, err = svc.CreateDevice(ctx, "Device A", "BrandX", "123")
		assert.Error(t, err)
		assert.Empty(t, device)
	})
}

func TestUpdateDevice(t *testing.T) {
	ctx, _, svc := setupTest(t)

	device, err := svc.CreateDevice(ctx, "Device", "BrandY", store.DeviceStateAvailable)
	assert.NoError(t, err)

	t.Run("It_should_be_able_to_update_a_device_available", func(t *testing.T) {
		deviceUpdated, err := svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand", store.DeviceStateInactive)
		assert.NoError(t, err)
		assert.Equal(t, "Device Updated", deviceUpdated.Name)
		assert.Equal(t, store.DeviceStateInactive, deviceUpdated.State)
		assert.Equal(t, "Brand", deviceUpdated.Brand)
	})

	t.Run("It_should_be_able_to_update_a_device_inactive", func(t *testing.T) {
		deviceUpdated, err := svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand2", store.DeviceStateAvailable)
		assert.NoError(t, err)
		assert.Equal(t, "Device Updated", deviceUpdated.Name)
		assert.Equal(t, store.DeviceStateAvailable, deviceUpdated.State)
		assert.Equal(t, "Brand2", deviceUpdated.Brand)
	})

	t.Run("It_should_not_be_able_to_update_name_if_device_in_use_state", func(t *testing.T) {
		_, err := svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInUse)
		assert.NoError(t, err)

		_, err = svc.UpdateDevice(ctx, device.ID, "Device Updated2", "Brand3", store.DeviceStateInUse)
		assert.ErrorIs(t, err, ErrDeviceInUse)
	})

	t.Run("It_should_not_be_able_to_update_brand_if_device_in_use_state", func(t *testing.T) {
		_, err := svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInUse)
		assert.NoError(t, err)

		_, err = svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand4", store.DeviceStateInUse)
		assert.ErrorIs(t, err, ErrDeviceInUse)
	})

	t.Run("It_should_be_able_to_update_state_if_device_in_use_state", func(t *testing.T) {
		_, err := svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInUse)
		assert.NoError(t, err)

		deviceUpdated, err := svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInactive)
		assert.NoError(t, err)
		assert.Equal(t, store.DeviceStateInactive, deviceUpdated.State)
	})

	t.Run("It_should_not_be_able_to_update_a_device_that_does_not_exist", func(t *testing.T) {
		_, err := svc.UpdateDevice(ctx, 999, "Device D", "BrandZ", store.DeviceStateAvailable)
		assert.ErrorIs(t, err, ErrDeviceNotFound)
	})
}

func TestPatchDevice(t *testing.T) {
	ctx, _, svc := setupTest(t)

	device, err := svc.CreateDevice(ctx, "Device", "BrandY", store.DeviceStateAvailable)
	assert.NoError(t, err)

	t.Run("It_should_be_able_to_patch_a_device_available", func(t *testing.T) {
		deviceUpdated, err := svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand", store.DeviceStateInactive)
		assert.NoError(t, err)
		assert.Equal(t, "Device Updated", deviceUpdated.Name)
		assert.Equal(t, store.DeviceStateInactive, deviceUpdated.State)
		assert.Equal(t, "Brand", deviceUpdated.Brand)
	})

	t.Run("It_should_be_able_to_patch_a_device_inactive", func(t *testing.T) {
		deviceUpdated, err := svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand2", store.DeviceStateAvailable)
		assert.NoError(t, err)
		assert.Equal(t, "Device Updated", deviceUpdated.Name)
		assert.Equal(t, store.DeviceStateAvailable, deviceUpdated.State)
		assert.Equal(t, "Brand2", deviceUpdated.Brand)
	})

	t.Run("It_should_not_be_able_to_patch_name_if_device_in_use_state", func(t *testing.T) {
		_, err := svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInUse)
		assert.NoError(t, err)

		_, err = svc.PatchDevice(ctx, device.ID, "Device Updated2", "Brand3", store.DeviceStateInUse)
		assert.ErrorIs(t, err, ErrDeviceInUse)
	})

	t.Run("It_should_not_be_able_to_patch_brand_if_device_in_use_state", func(t *testing.T) {
		_, err := svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInUse)
		assert.NoError(t, err)

		_, err = svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand4", store.DeviceStateInUse)
		assert.ErrorIs(t, err, ErrDeviceInUse)
	})

	t.Run("It_should_be_able_to_patch_state_if_device_in_use_state", func(t *testing.T) {
		_, err := svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInUse)
		assert.NoError(t, err)

		deviceUpdated, err := svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInactive)
		assert.NoError(t, err)
		assert.Equal(t, store.DeviceStateInactive, deviceUpdated.State)
	})

	t.Run("It_should_not_be_able_to_patch_a_device_that_does_not_exist", func(t *testing.T) {
		_, err := svc.PatchDevice(ctx, 999, "Device D", "BrandZ", store.DeviceStateAvailable)
		assert.ErrorIs(t, err, ErrDeviceNotFound)
	})

	t.Run("It_should_be_able_to_update_only_name", func(t *testing.T) {
		deviceTest, err := svc.CreateDevice(ctx, "Name1", "BrandShouldNotChange", store.DeviceStateAvailable)
		assert.NoError(t, err)

		deviceUpdated, err := svc.PatchDevice(ctx, deviceTest.ID, "NameShouldUpdate", "", "")
		assert.NoError(t, err)
		assert.Equal(t, "NameShouldUpdate", deviceUpdated.Name)
		assert.Equal(t, "BrandShouldNotChange", deviceUpdated.Brand)
		assert.Equal(t, store.DeviceStateAvailable, deviceUpdated.State)
	})

	t.Run("It_should_be_able_to_update_only_brand", func(t *testing.T) {
		deviceTest, err := svc.CreateDevice(ctx, "NameShouldNotChange", "Brand1", store.DeviceStateAvailable)
		assert.NoError(t, err)

		deviceUpdated, err := svc.PatchDevice(ctx, deviceTest.ID, "", "BrandShouldUpdate", "")
		assert.NoError(t, err)
		assert.Equal(t, "NameShouldNotChange", deviceUpdated.Name)
		assert.Equal(t, "BrandShouldUpdate", deviceUpdated.Brand)
		assert.Equal(t, store.DeviceStateAvailable, deviceUpdated.State)
	})
}

func TestGetDeviceByID(t *testing.T) {
	ctx, _, svc := setupTest(t)

	deviceCreated, err := svc.CreateDevice(ctx, "Device", "BrandY", store.DeviceStateAvailable)
	assert.NoError(t, err)

	t.Run("It_should_be_able_to_get_a_device", func(t *testing.T) {
		device, err := svc.GetDeviceByID(ctx, deviceCreated.ID)
		assert.NoError(t, err)
		assert.Equal(t, deviceCreated.ID, device.ID)
	})

	t.Run("It_should_not_be_able_to_get_a_device_that_does_not_exist", func(t *testing.T) {
		_, err := svc.GetDeviceByID(ctx, 999)
		assert.ErrorIs(t, err, ErrDeviceNotFound)
	})
}

func TestDeleteDevice(t *testing.T) {
	ctx, _, svc := setupTest(t)

	deviceCreated, err := svc.CreateDevice(ctx, "Device", "BrandY", store.DeviceStateAvailable)
	assert.NoError(t, err)

	t.Run("It_should_be_able_to_delete_a_device", func(t *testing.T) {
		deletedID, err := svc.DeleteDevice(ctx, deviceCreated.ID)
		assert.NoError(t, err)
		assert.Equal(t, deviceCreated.ID, deletedID)
	})

	t.Run("It_should_not_be_able_to_delete_a_device_that_does_not_exist", func(t *testing.T) {
		_, err := svc.DeleteDevice(ctx, 999)
		assert.ErrorIs(t, err, ErrDeviceNotFound)
	})

	t.Run("It_should_not_be_able_to_delete_a_device_in_use", func(t *testing.T) {
		// Create a device in-use
		inUseDevice, err := svc.CreateDevice(ctx, "DeviceInUse", "BrandY", store.DeviceStateInUse)
		assert.NoError(t, err)

		_, err = svc.DeleteDevice(ctx, inUseDevice.ID)
		assert.ErrorIs(t, err, ErrDeviceInUse)
	})
}

func TestGetAllDevices(t *testing.T) {
	ctx, _, svc := setupTest(t)

	_, _ = svc.CreateDevice(ctx, "Device1", "BrandY", store.DeviceStateAvailable)
	_, _ = svc.CreateDevice(ctx, "Device2", "BrandY", store.DeviceStateAvailable)
	_, _ = svc.CreateDevice(ctx, "Device3", "BrandY", store.DeviceStateAvailable)
	_, _ = svc.CreateDevice(ctx, "Device4", "BrandX", store.DeviceStateInactive)
	_, _ = svc.CreateDevice(ctx, "Device5", "BrandX", store.DeviceStateInUse)
	_, _ = svc.CreateDevice(ctx, "Device6", "BrandY", store.DeviceStateInUse)
	_, _ = svc.CreateDevice(ctx, "Device7", "BrandY", store.DeviceStateInactive)

	t.Run("It_should_be_able_to_get_all_devices", func(t *testing.T) {
		devices, err := svc.GetAllDevices(ctx, "", "")
		assert.NoError(t, err)
		assert.Len(t, devices, 7)
	})

	t.Run("It_should_be_able_to_get_all_devices_by_brand", func(t *testing.T) {
		devices, err := svc.GetAllDevices(ctx, "BrandY", "")
		assert.NoError(t, err)
		assert.Len(t, devices, 5)
	})

	t.Run("It_should_be_able_to_get_all_devices_by_state", func(t *testing.T) {
		devices, err := svc.GetAllDevices(ctx, "", store.DeviceStateInactive)
		assert.NoError(t, err)
		assert.Len(t, devices, 2)
	})

	t.Run("It_should_be_able_to_get_all_devices_by_brand_and_state", func(t *testing.T) {
		devices, err := svc.GetAllDevices(ctx, "BrandY", store.DeviceStateAvailable)
		assert.NoError(t, err)
		assert.Len(t, devices, 3)
	})
}
