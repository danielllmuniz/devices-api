package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *Api) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/devices", api.handleCreateDevice)
			r.Get("/devices", api.handleGetAllDevices)
			r.Get("/devices/{device_id}", api.handleGetDevice)
			r.Patch("/devices/{device_id}", api.handlePatchDevice)
			r.Delete("/devices/{device_id}", api.handleDeleteDevice)
			r.Put("/devices/{device_id}", api.handleUpdateDevice)
		})
	})
}
