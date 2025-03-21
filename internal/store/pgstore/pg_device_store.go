package pgstore

import (
	"context"

	"github.com/danielllmuniz/devices-api/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGDeviceStore struct {
	Queries *Queries
	db      *pgxpool.Pool
}

func NewPGDeviceStore(db *pgxpool.Pool) *PGDeviceStore {
	return &PGDeviceStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *PGDeviceStore) CreateDevice(ctx context.Context, name, brand string, state store.DeviceState) (store.Device, error) {
	device, err := s.Queries.CreateDevice(ctx, CreateDeviceParams{
		Name:  name,
		Brand: brand,
		State: DeviceState(state),
	})
	if err != nil {
		return store.Device{}, err
	}
	return store.Device{
		ID:        device.ID,
		Name:      device.Name,
		Brand:     device.Brand,
		State:     store.DeviceState(device.State),
		CreatedAt: device.CreatedAt,
	}, nil
}

func (s *PGDeviceStore) UpdateDevice(ctx context.Context, id int32, name, brand string, state store.DeviceState) (store.Device, error) {
	device, err := s.Queries.UpdateDevice(ctx, UpdateDeviceParams{
		ID:    id,
		Name:  name,
		Brand: brand,
		State: DeviceState(state),
	})
	if err != nil {
		return store.Device{}, err
	}
	return store.Device{
		ID:        device.ID,
		Name:      device.Name,
		Brand:     device.Brand,
		State:     store.DeviceState(device.State),
		CreatedAt: device.CreatedAt,
	}, nil
}

func (s *PGDeviceStore) PatchDevice(ctx context.Context, id int32, name, brand string, state store.DeviceState) (store.Device, error) {
	device, err := s.Queries.PatchDevice(ctx, PatchDeviceParams{
		ID:      id,
		Column2: name,
		Column3: brand,
		Column4: state,
	})
	if err != nil {
		return store.Device{}, err
	}
	return store.Device{
		ID:        device.ID,
		Name:      device.Name,
		Brand:     device.Brand,
		State:     store.DeviceState(device.State),
		CreatedAt: device.CreatedAt,
	}, nil
}

func (s *PGDeviceStore) GetDeviceByID(ctx context.Context, id int32) (store.Device, error) {
	device, err := s.Queries.GetDeviceById(ctx, id)
	if err != nil {
		return store.Device{}, err
	}
	return store.Device{
		ID:        device.ID,
		Name:      device.Name,
		Brand:     device.Brand,
		State:     store.DeviceState(device.State),
		CreatedAt: device.CreatedAt,
	}, nil
}

func (s *PGDeviceStore) GetAllDevices(ctx context.Context) ([]store.Device, error) {
	devices, err := s.Queries.GetAllDevices(ctx)
	if err != nil {
		return nil, err
	}

	var result []store.Device
	for _, d := range devices {
		result = append(result, store.Device{
			ID:        d.ID,
			Name:      d.Name,
			Brand:     d.Brand,
			State:     store.DeviceState(d.State),
			CreatedAt: d.CreatedAt,
		})
	}
	return result, nil
}

func (s *PGDeviceStore) GetDevicesByBrand(ctx context.Context, brand string) ([]store.Device, error) {
	devices, err := s.Queries.GetDevicesByBrand(ctx, brand)
	if err != nil {
		return nil, err
	}

	var result []store.Device
	for _, d := range devices {
		result = append(result, store.Device{
			ID:        d.ID,
			Name:      d.Name,
			Brand:     d.Brand,
			State:     store.DeviceState(d.State),
			CreatedAt: d.CreatedAt,
		})
	}
	return result, nil
}

func (s *PGDeviceStore) GetDevicesByState(ctx context.Context, state store.DeviceState) ([]store.Device, error) {
	devices, err := s.Queries.GetDevicesByState(ctx, DeviceState(state))
	if err != nil {
		return nil, err
	}

	var result []store.Device
	for _, d := range devices {
		result = append(result, store.Device{
			ID:        d.ID,
			Name:      d.Name,
			Brand:     d.Brand,
			State:     store.DeviceState(d.State),
			CreatedAt: d.CreatedAt,
		})
	}
	return result, nil
}

func (s *PGDeviceStore) GetDevicesByBrandAndState(ctx context.Context, brand string, state store.DeviceState) ([]store.Device, error) {
	devices, err := s.Queries.GetDevicesByBrandAndState(ctx, GetDevicesByBrandAndStateParams{
		Brand: brand,
		State: DeviceState(state),
	})
	if err != nil {
		return nil, err
	}

	var result []store.Device
	for _, d := range devices {
		result = append(result, store.Device{
			ID:        d.ID,
			Name:      d.Name,
			Brand:     d.Brand,
			State:     store.DeviceState(d.State),
			CreatedAt: d.CreatedAt,
		})
	}
	return result, nil
}

func (s *PGDeviceStore) DeleteDevice(ctx context.Context, id int32) (int32, error) {
	deletedID, err := s.Queries.DeleteDevice(ctx, id)
	if err != nil {
		return 0, err
	}
	return deletedID, nil
}
