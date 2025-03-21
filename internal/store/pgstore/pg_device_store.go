package pgstore

import (
	"context"

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

func (s *PGDeviceStore) CreateDevice(ctx context.Context, name, brand string, state DeviceState) (Device, error) {
	device, err := s.Queries.CreateDevice(ctx, CreateDeviceParams{
		Name:  name,
		Brand: brand,
		State: state,
	})
	if err != nil {
		return Device{}, err
	}
	return device, nil
}

func (s *PGDeviceStore) UpdateDevice(ctx context.Context, id int32, name, brand string, state DeviceState) (Device, error) {
	device, err := s.Queries.UpdateDevice(ctx, UpdateDeviceParams{
		ID:    id,
		Name:  name,
		Brand: brand,
		State: state,
	})
	if err != nil {
		return Device{}, err
	}
	return device, nil
}

func (s *PGDeviceStore) PatchDevice(ctx context.Context, id int32, name, brand string, state DeviceState) (Device, error) {
	device, err := s.Queries.PatchDevice(ctx, PatchDeviceParams{
		ID:      id,
		Column2: name,
		Column3: brand,
		Column4: state,
	})
	if err != nil {
		return Device{}, err
	}
	return device, nil
}

func (s *PGDeviceStore) GetDeviceByID(ctx context.Context, id int32) (Device, error) {
	device, err := s.Queries.GetDeviceById(ctx, id)
	if err != nil {
		return Device{}, err
	}
	return device, nil
}

func (s *PGDeviceStore) GetAllDevices(ctx context.Context) ([]Device, error) {
	devices, err := s.Queries.GetAllDevices(ctx)
	if err != nil {
		return nil, err
	}

	var result []Device
	for _, d := range devices {
		result = append(result, d)
	}
	return result, nil
}

func (s *PGDeviceStore) GetDevicesByBrand(ctx context.Context, brand string) ([]Device, error) {
	devices, err := s.Queries.GetDevicesByBrand(ctx, brand)
	if err != nil {
		return nil, err
	}

	var result []Device
	for _, d := range devices {
		result = append(result, d)
	}
	return result, nil
}

func (s *PGDeviceStore) GetDevicesByState(ctx context.Context, state DeviceState) ([]Device, error) {
	devices, err := s.Queries.GetDevicesByState(ctx, state)
	if err != nil {
		return nil, err
	}

	var result []Device
	for _, d := range devices {
		result = append(result, d)
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
