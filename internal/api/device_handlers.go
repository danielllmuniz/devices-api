package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danielllmuniz/devices-api/internal/jsonutils"
	"github.com/danielllmuniz/devices-api/internal/services"
	"github.com/danielllmuniz/devices-api/internal/store"
	deviceValidator "github.com/danielllmuniz/devices-api/internal/validator/device"
	"github.com/go-chi/chi/v5"
)

func (api *Api) handleCreateDevice(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[deviceValidator.CreateDeviceReq](r)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(problems)
		if problems == nil {
			jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": "invalid request",
			})
			return
		}

		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	device, err := api.DeviceService.CreateDevice(
		r.Context(),
		data.Name,
		data.Brand,
		data.State,
	)
	if err != nil {
		fmt.Println(err.Error())
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create device, try again later",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"message": "device created successfully",
		"device": map[string]any{
			"id":         device.ID,
			"name":       device.Name,
			"brand":      device.Brand,
			"state":      device.State,
			"created_at": device.CreatedAt,
		},
	})
}

func (api *Api) handleGetDevice(w http.ResponseWriter, r *http.Request) {
	strDeviceID := chi.URLParam(r, "device_id")

	intDeviceID, err := strconv.Atoi(strDeviceID)
	if err != nil {
		fmt.Println(err.Error())
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid device id",
		})
		return
	}

	device, err := api.DeviceService.GetDeviceByID(r.Context(), int32(intDeviceID))
	if err != nil {
		fmt.Println(err.Error())
		jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
			"error": "device not found",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"device": map[string]any{
			"id":         device.ID,
			"name":       device.Name,
			"brand":      device.Brand,
			"state":      device.State,
			"created_at": device.CreatedAt,
		},
	})
}

func (api *Api) handleGetAllDevices(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	brand := queryParams.Get("brand")
	state := queryParams.Get("state")

	devices, err := api.DeviceService.GetAllDevices(r.Context(), brand, store.DeviceState(state))
	if err != nil {
		fmt.Println(err.Error())
		jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
			"error": "no devices found",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"devices": devices,
	})
}

func (api *Api) handleUpdateDevice(w http.ResponseWriter, r *http.Request) {
	strDeviceID := chi.URLParam(r, "device_id")

	intDeviceID, err := strconv.Atoi(strDeviceID)
	if err != nil {
		fmt.Println(err.Error())
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid device id",
		})
		return
	}

	data, problems, err := jsonutils.DecodeValidJson[deviceValidator.UpdateDeviceReq](r)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(problems)
		if problems == nil {
			jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": "invalid request",
			})
			return
		}

		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	device, err := api.DeviceService.UpdateDevice(
		r.Context(),
		int32(intDeviceID),
		data.Name,
		data.Brand,
		store.DeviceState(data.State),
	)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, services.ErrDeviceInUse) {
			jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "device is in use, cannot update name or brand",
			})
			return
		}

		if errors.Is(err, services.ErrDeviceNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"error": "device not found",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to update device, try again later",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "device updated successfully",
		"device": map[string]any{
			"id":         device.ID,
			"name":       device.Name,
			"brand":      device.Brand,
			"state":      device.State,
			"created_at": device.CreatedAt,
		},
	})
}

func (api *Api) handlePatchDevice(w http.ResponseWriter, r *http.Request) {
	strDeviceID := chi.URLParam(r, "device_id")

	intDeviceID, err := strconv.Atoi(strDeviceID)
	if err != nil {
		fmt.Println(err.Error())
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid device id",
		})
		return
	}

	data, problems, err := jsonutils.DecodeValidJson[deviceValidator.PatchDeviceReq](r)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(problems)
		if problems == nil {
			jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": "invalid request",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	device, err := api.DeviceService.PatchDevice(
		r.Context(),
		int32(intDeviceID),
		data.Name,
		data.Brand,
		store.DeviceState(data.State),
	)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, services.ErrDeviceInUse) {
			jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "device is in use, cannot update name or brand",
			})
			return
		}

		if errors.Is(err, services.ErrDeviceNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"error": "device not found",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to update device, try again later",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "device patched successfully",
		"device": map[string]any{
			"id":         device.ID,
			"name":       device.Name,
			"brand":      device.Brand,
			"state":      device.State,
			"created_at": device.CreatedAt,
		},
	})
}

func (api *Api) handleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	strDeviceID := chi.URLParam(r, "device_id")

	intDeviceID, err := strconv.Atoi(strDeviceID)
	if err != nil {
		fmt.Println(err.Error())
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid device id",
		})
		return
	}

	id, err := api.DeviceService.DeleteDevice(r.Context(), int32(intDeviceID))
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, services.ErrDeviceInUse) {
			jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "device is in use, cannot update name or brand",
			})
			return
		}

		if errors.Is(err, services.ErrDeviceNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"error": "device not found",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to update device, try again later",
		})
		return
	}

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message":   "device deleted successfully",
		"device_id": id,
	})
}
