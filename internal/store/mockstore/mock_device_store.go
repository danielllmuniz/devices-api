package mockstore

import (
	"context"
	"errors"
	"sync"

	"github.com/danielllmuniz/devices-api/internal/store"
)

type MockDeviceStore struct {
	mu      sync.Mutex
	devices map[int32]store.Device
	nextID  int32
}

func NewMockDeviceStore() *MockDeviceStore {
	return &MockDeviceStore{
		devices: make(map[int32]store.Device),
		nextID:  1,
	}
}

func (m *MockDeviceStore) CreateDevice(ctx context.Context, name, brand string, state store.DeviceState) (store.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if state != store.DeviceStateAvailable && state != store.DeviceStateInUse && state != store.DeviceStateInactive {
		return store.Device{}, errors.New("invalid state")
	}

	device := store.Device{
		ID:    m.nextID,
		Name:  name,
		Brand: brand,
		State: state,
	}
	m.devices[m.nextID] = device
	m.nextID++

	return device, nil
}

func (m *MockDeviceStore) UpdateDevice(ctx context.Context, id int32, name, brand string, state store.DeviceState) (store.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.devices[id]; !ok {
		return store.Device{}, errors.New("device not found")
	}

	device := store.Device{
		ID:    id,
		Name:  name,
		Brand: brand,
		State: state,
	}
	m.devices[id] = device
	return device, nil
}

func (m *MockDeviceStore) PatchDevice(ctx context.Context, id int32, name, brand string, state store.DeviceState) (store.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	device, ok := m.devices[id]
	if !ok {
		return store.Device{}, errors.New("device not found")
	}

	if name != "" {
		device.Name = name
	}
	if brand != "" {
		device.Brand = brand
	}
	if state != "" {
		device.State = state
	}
	m.devices[id] = device

	return device, nil
}

func (m *MockDeviceStore) GetDeviceByID(ctx context.Context, id int32) (store.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	device, ok := m.devices[id]
	if !ok {
		return store.Device{}, errors.New("device not found")
	}
	return device, nil
}

func (m *MockDeviceStore) GetAllDevices(ctx context.Context) ([]store.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []store.Device
	for _, device := range m.devices {
		result = append(result, device)
	}
	return result, nil
}

func (m *MockDeviceStore) GetDevicesByBrand(ctx context.Context, brand string) ([]store.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []store.Device
	for _, device := range m.devices {
		if device.Brand == brand {
			result = append(result, device)
		}
	}
	return result, nil
}

func (m *MockDeviceStore) GetDevicesByState(ctx context.Context, state store.DeviceState) ([]store.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []store.Device
	for _, device := range m.devices {
		if device.State == state {
			result = append(result, device)
		}
	}
	return result, nil
}

func (m *MockDeviceStore) GetDevicesByBrandAndState(ctx context.Context, brand string, state store.DeviceState) ([]store.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []store.Device
	for _, device := range m.devices {
		if device.Brand == brand && device.State == state {
			result = append(result, device)
		}
	}
	return result, nil
}

func (m *MockDeviceStore) DeleteDevice(ctx context.Context, id int32) (int32, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.devices[id]; !ok {
		return 0, errors.New("device not found")
	}
	delete(m.devices, id)
	return id, nil
}
