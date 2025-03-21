package services

import (
	"context"
	"errors"

	"github.com/danielllmuniz/devices-api/internal/store"
)

var (
	ErrDeviceInUse         = errors.New("device is currently in use and cannot be modified or deleted")
	ErrCannotUpdateCreated = errors.New("creation time cannot be updated")
	ErrDeviceNotFound      = errors.New("device not found")
)

type DeviceService struct {
	Store store.DeviceStore
}

func NewDeviceService(store store.DeviceStore) *DeviceService {
	return &DeviceService{Store: store}
}

func (s *DeviceService) CreateDevice(ctx context.Context, name, brand string, state store.DeviceState) (store.Device, error) {
	device, err := s.Store.CreateDevice(ctx, name, brand, state)
	if err != nil {
		return store.Device{}, err
	}
	return device, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, id int32, name, brand string, state store.DeviceState) (store.Device, error) {
	device, err := s.Store.GetDeviceByID(ctx, id)
	if err != nil {
		return store.Device{}, ErrDeviceNotFound
	}

	if device.State == store.DeviceStateInUse && (name != device.Name || brand != device.Brand) {
		return store.Device{}, ErrDeviceInUse
	}

	deviceUpdated, err := s.Store.UpdateDevice(ctx, id, name, brand, state)
	if err != nil {
		return store.Device{}, err
	}
	return deviceUpdated, nil
}

func (s *DeviceService) PatchDevice(ctx context.Context, id int32, name, brand string, state store.DeviceState) (store.Device, error) {
	device, err := s.Store.GetDeviceByID(ctx, id)
	if err != nil {
		return store.Device{}, ErrDeviceNotFound
	}

	if device.State == store.DeviceStateInUse && (name != device.Name || brand != device.Brand) {
		return store.Device{}, ErrDeviceInUse
	}

	deviceUpdated, err := s.Store.PatchDevice(ctx, id, name, brand, state)
	if err != nil {
		return store.Device{}, err
	}
	return deviceUpdated, nil
}

func (s *DeviceService) GetDeviceByID(ctx context.Context, id int32) (store.Device, error) {
	device, err := s.Store.GetDeviceByID(ctx, id)
	if err != nil {
		return store.Device{}, ErrDeviceNotFound
	}
	return device, nil
}

func (s *DeviceService) GetAllDevices(ctx context.Context, brand string, state store.DeviceState) ([]store.Device, error) {

	if brand != "" && state != "" {
		devices, err := s.Store.GetDevicesByBrandAndState(ctx, brand, state)
		if err != nil {
			return nil, err
		}
		return devices, nil
	} else if brand != "" {
		devices, err := s.Store.GetDevicesByBrand(ctx, brand)
		if err != nil {
			return nil, err
		}
		return devices, nil
	} else if state != "" {
		devices, err := s.Store.GetDevicesByState(ctx, state)
		if err != nil {
			return nil, err
		}
		return devices, nil
	}
	devices, err := s.Store.GetAllDevices(ctx)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, id int32) (int32, error) {
	device, err := s.Store.GetDeviceByID(ctx, id)
	if err != nil {
		return 0, ErrDeviceNotFound
	}

	if device.State == store.DeviceStateInUse {
		return 0, ErrDeviceInUse
	}

	return s.Store.DeleteDevice(ctx, id)
}
