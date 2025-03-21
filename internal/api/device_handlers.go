package api

import (
	"fmt"
	"net/http"

	"github.com/danielllmuniz/devices-api/internal/jsonutils"
	deviceValidator "github.com/danielllmuniz/devices-api/internal/validator/device"
)

func (api *Api) handleCreateDevice(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[deviceValidator.CreateDeviceReq](r)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(problems)

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
