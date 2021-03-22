package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Iot struct {
	IotID              int       `json:"iot_id"`
	User               *User     `json:"user"`
	Group              *IotGroup `json:"group"`
	Longitude          float32   `json:"longitude"`
	Latitude           float32   `json:"latitude"`
	LastUpdateTimeUnix int64     `json:"last_update_time_unix"`
	IotState           string    `json:"iot_state"`
	IotType            string    `json:"iot_type"`
}

func (i *Iot) Validate() error {
	if err := validation.ValidateStruct(
		i,
		validation.Field(&i.Longitude, validation.Required, is.Longitude),
		validation.Field(&i.Latitude, validation.Required, is.Latitude),
		validation.Field(&i.LastUpdateTimeUnix, validation.Required, validation.Min(0)),
		validation.Field(&i.IotState, validation.Required),
		validation.Field(&i.IotType, validation.Required),
	); err != nil {
		return err
	}

	if err := i.ValidateState(); err != nil {
		return err
	}

	return i.ValidateType()
}

func (i *Iot) ValidateState() error {
	if i.IotState == "nothing" || i.IotState == "active" || i.IotState == "lost" {
		return nil
	}

	return ErrIncorrectState
}

func (i *Iot) ValidateType() error {
	if i.IotType == "gyroscope" || i.IotType == "microphone" {
		return nil
	}

	return ErrIncorrectType
}
