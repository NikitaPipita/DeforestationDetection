package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type IotGroup struct {
	ID                     int   `json:"group_id"`
	User                   *User `json:"user,omitempty"`
	UpdateDurationSeconds  int   `json:"update_duration_seconds"`
	LastIotChangesTimeUnix int   `json:"last_iot_changes_time_unix"`
}

func (g *IotGroup) Validate() error {
	return validation.ValidateStruct(
		g,
		validation.Field(&g.UpdateDurationSeconds, validation.Min(60)),
		validation.Field(&g.LastIotChangesTimeUnix, validation.Min(0)),
	)
}
