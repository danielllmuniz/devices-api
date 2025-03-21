package services

import (
	"context"
	"testing"

	"github.com/danielllmuniz/devices-api/internal/store"
	"github.com/danielllmuniz/devices-api/internal/store/mockstore"
	"github.com/stretchr/testify/assert"
)

func TestCreateDevice(t *testing.T) {
	ctx := context.Background()
	Store := mockstore.NewMockDeviceStore()
	svc := NewDeviceService(Store)

	t.Run("It_should_be_able_to_create_a_device", func(t *testing.T) {
		device, err := svc.CreateDevice(ctx, "Device A", "BrandX", store.DeviceStateInUse)
		assert.NoError(t, err)
		assert.Equal(t, "Device A", device.Name)
	})
}

func TestUpdateDevice(t *testing.T) {
	ctx := context.Background()
	Store := mockstore.NewMockDeviceStore()
	svc := NewDeviceService(Store)

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
		_, err = svc.UpdateDevice(ctx, device.ID, "Device Updated2", "Brand3", store.DeviceStateInUse)
		assert.ErrorIs(t, err, ErrDeviceInUse)
	})

	t.Run("It_should_not_be_able_to_update_brand_if_device_in_use_state", func(t *testing.T) {
		_, err := svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInUse)
		_, err = svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand4", store.DeviceStateInUse)
		assert.ErrorIs(t, err, ErrDeviceInUse)
	})

	t.Run("It_should_be_able_to_update_state_if_device_in_use_state", func(t *testing.T) {
		_, err := svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInUse)
		deviceUpdated, err := svc.UpdateDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInactive)
		assert.NoError(t, err)
		assert.Equal(t, store.DeviceStateInactive, deviceUpdated.State)
	})

	t.Run("It_should_not_be_able_to_update_a_device_that_does_not_exist", func(t *testing.T) {
		_, err := svc.UpdateDevice(ctx, 4, "Device D", "BrandZ", store.DeviceStateAvailable)
		assert.ErrorIs(t, err, ErrDeviceNotFound)
	})
}

func TestPatchDevice(t *testing.T) {
	ctx := context.Background()
	Store := mockstore.NewMockDeviceStore()
	svc := NewDeviceService(Store)

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
		_, err = svc.PatchDevice(ctx, device.ID, "Device Updated2", "Brand3", store.DeviceStateInUse)
		assert.ErrorIs(t, err, ErrDeviceInUse)
	})

	t.Run("It_should_not_be_able_to_patch_brand_if_device_in_use_state", func(t *testing.T) {
		_, err := svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInUse)
		_, err = svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand4", store.DeviceStateInUse)
		assert.ErrorIs(t, err, ErrDeviceInUse)
	})

	t.Run("It_should_be_able_to_patch_state_if_device_in_use_state", func(t *testing.T) {
		_, err := svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInUse)
		deviceUpdated, err := svc.PatchDevice(ctx, device.ID, "Device Updated", "Brand3", store.DeviceStateInactive)
		assert.NoError(t, err)
		assert.Equal(t, store.DeviceStateInactive, deviceUpdated.State)
	})
	t.Run("It_should_not_be_able_to_patch_a_device_that_does_not_exist", func(t *testing.T) {
		_, err := svc.PatchDevice(ctx, 4, "Device D", "BrandZ", store.DeviceStateAvailable)
		assert.ErrorIs(t, err, ErrDeviceNotFound)
	})
}

func TestGetDeviceByID(t *testing.T) {
	ctx := context.Background()
	Store := mockstore.NewMockDeviceStore()
	svc := NewDeviceService(Store)

	deviceCreated, err := svc.CreateDevice(ctx, "Device", "BrandY", store.DeviceStateAvailable)
	assert.NoError(t, err)
	t.Run("It_should_be_able_to_get_a_device", func(t *testing.T) {
		device, err := svc.GetDeviceByID(ctx, deviceCreated.ID)
		assert.NoError(t, err)
		assert.Equal(t, deviceCreated.ID, device.ID)
	})

	t.Run("It_should_not_be_able_to_get_a_device_that_does_not_exist", func(t *testing.T) {
		_, err := svc.GetDeviceByID(ctx, 4)
		assert.ErrorIs(t, err, ErrDeviceNotFound)
	})
}

func TestDeleteDevice(t *testing.T) {
	ctx := context.Background()
	Store := mockstore.NewMockDeviceStore()
	svc := NewDeviceService(Store)

	deviceCreated, err := svc.CreateDevice(ctx, "Device", "BrandY", store.DeviceStateAvailable)
	assert.NoError(t, err)
	t.Run("It_should_be_able_to_delete_a_device", func(t *testing.T) {
		deletedID, err := svc.DeleteDevice(ctx, deviceCreated.ID)
		assert.NoError(t, err)
		assert.Equal(t, deviceCreated.ID, deletedID)
	})

	t.Run("It_should_not_be_able_to_delete_a_device_that_does_not_exist", func(t *testing.T) {
		_, err := svc.DeleteDevice(ctx, 4)
		assert.ErrorIs(t, err, ErrDeviceNotFound)
	})

	t.Run("It_should_not_be_able_to_delete_a_device_in_use", func(t *testing.T) {
		_, err := svc.CreateDevice(ctx, "Device", "BrandY", store.DeviceStateInUse)
		_, err = svc.DeleteDevice(ctx, 2)
		assert.ErrorIs(t, err, ErrDeviceInUse)
	})
}

func TestGetAllDevices(t *testing.T) {
	ctx := context.Background()
	Store := mockstore.NewMockDeviceStore()
	svc := NewDeviceService(Store)

	t.Run("It_should_be_able_to_get_all_devices", func(t *testing.T) {
		_, _ = svc.CreateDevice(ctx, "Device", "BrandY", store.DeviceStateAvailable)
		_, _ = svc.CreateDevice(ctx, "Device2", "BrandY", store.DeviceStateAvailable)
		_, _ = svc.CreateDevice(ctx, "Device3", "BrandY", store.DeviceStateAvailable)

		devices, err := svc.GetAllDevices(ctx, "", "")
		assert.NoError(t, err)
		assert.Len(t, devices, 3)
	})
}

func TestGetDevicesByBrand(t *testing.T) {
	ctx := context.Background()
	Store := mockstore.NewMockDeviceStore()
	svc := NewDeviceService(Store)

	t.Run("It_should_be_able_to_get_devices_by_brand", func(t *testing.T) {
		_, _ = svc.CreateDevice(ctx, "Device", "BrandY", store.DeviceStateAvailable)
		_, _ = svc.CreateDevice(ctx, "Device2", "BrandY", store.DeviceStateAvailable)
		_, _ = svc.CreateDevice(ctx, "Device3", "BrandZ", store.DeviceStateAvailable)

		devices, err := svc.GetDevicesByBrand(ctx, "BrandY")
		assert.NoError(t, err)
		assert.Len(t, devices, 2)
	})
}

func TestGetDevicesByState(t *testing.T) {
	ctx := context.Background()
	Store := mockstore.NewMockDeviceStore()
	svc := NewDeviceService(Store)

	t.Run("It_should_be_able_to_get_devices_by_state", func(t *testing.T) {
		_, _ = svc.CreateDevice(ctx, "Device", "BrandY", store.DeviceStateAvailable)
		_, _ = svc.CreateDevice(ctx, "Device2", "BrandY", store.DeviceStateInUse)
		_, _ = svc.CreateDevice(ctx, "Device3", "BrandZ", store.DeviceStateAvailable)

		devices, err := svc.GetDevicesByState(ctx, store.DeviceStateAvailable)
		assert.NoError(t, err)
		assert.Len(t, devices, 2)
	})
}
