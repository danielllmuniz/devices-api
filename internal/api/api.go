package api

import (
	"github.com/danielllmuniz/devices-api/internal/services"
	"github.com/go-chi/chi/v5"
)

type Api struct {
	Router        *chi.Mux
	DeviceService *services.DeviceService
}
