package store

import (
	"context"
	"time"
)

type DeviceState string

const (
	DeviceStateAvailable DeviceState = "available"
	DeviceStateInUse     DeviceState = "in-use"
	DeviceStateInactive  DeviceState = "inactive"
)

type Device struct {
	ID        int32       `json:"id"`
	Name      string      `json:"name"`
	Brand     string      `json:"brand"`
	State     DeviceState `json:"state"`
	CreatedAt time.Time   `json:"created_at"`
}

type DeviceStore interface {
	CreateDevice(ctx context.Context, name, brand string, state DeviceState) (Device, error)
	UpdateDevice(ctx context.Context, id int32, name, brand string, state DeviceState) (Device, error)
	PatchDevice(ctx context.Context, id int32, name, brand, state string) (Device, error)
	GetDeviceByID(ctx context.Context, id int32) (Device, error)
	GetAllDevices(ctx context.Context) ([]Device, error)
	GetDevicesByBrand(ctx context.Context, brand string) ([]Device, error)
	GetDevicesByState(ctx context.Context, state DeviceState) ([]Device, error)
	DeleteDevice(ctx context.Context, id int32) (int32, error)
}
