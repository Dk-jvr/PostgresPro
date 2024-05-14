package Models

type (
	ScriptData struct {
		Id     int64
		Script string `json:"script" validate:"required"`
		Type   string `json:"type" validate:"required"`
	}
)
