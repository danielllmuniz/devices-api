package device

import (
	"context"

	"github.com/danielllmuniz/devices-api/internal/store"
	"github.com/danielllmuniz/devices-api/internal/validator"
)

type PatchDeviceReq struct {
	Name  string `json:"name"`
	Brand string `json:"brand"`
	State string `json:"state"`
}

func (req PatchDeviceReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	eval.CheckField(validator.MinChars(req.Name, 3) && validator.MaxChars(req.Name, 255), "name", "Name must be between 3 and 255 characters")
	eval.CheckField(validator.MinChars(req.Brand, 3) && validator.MaxChars(req.Brand, 255), "brand", "Brand must be between 3 and 255 characters")
	eval.CheckField(validator.InEnum(req.State, []any{store.DeviceStateAvailable, store.DeviceStateInUse, store.DeviceStateInactive}), "state", "State must be 'available', 'in-use' or 'inactive'")
	return eval
}
