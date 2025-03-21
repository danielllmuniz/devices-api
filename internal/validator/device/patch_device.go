package device

import (
	"context"

	"github.com/danielllmuniz/devices-api/internal/validator"
)

type PatchDeviceReq struct {
	Name  string `json:"name"`
	Brand string `json:"brand"`
	State string `json:"state"`
}

func (req PatchDeviceReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	if req.Name == "" && req.Brand == "" && req.State == "" {
		eval.AddFieldError("name", "At least one field must be informed")
		eval.AddFieldError("brand", "At least one field must be informed")
		eval.AddFieldError("state", "At least one field must be informed")
		return eval
	}
	if req.Name != "" {
		eval.CheckField(validator.MinChars(req.Name, 3) && validator.MaxChars(req.Name, 255), "name", "Name must be between 3 and 255 characters")
	}
	if req.Brand != "" {
		eval.CheckField(validator.MinChars(req.Brand, 3) && validator.MaxChars(req.Brand, 255), "brand", "Brand must be between 3 and 255 characters")
	}
	if req.State != "" {
		eval.CheckField(validator.InEnum(req.State, []interface{}{"available", "in-use", "inactive"}), "state", "State must be 'available', 'in-use' or 'inactive")
	}
	return eval
}
