package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type Iot struct {
	ID                 int       `json:"iot_id"`
	User               *User     `json:"user"`
	Group              *IotGroup `json:"group"`
	Longitude          float64   `json:"longitude"`
	Latitude           float64   `json:"latitude"`
	LastUpdateTimeUnix int64     `json:"last_update_time_unix"`
	IotState           string    `json:"iot_state"`
	IotType            string    `json:"iot_type"`
	Password           string    `json:"password,omitempty"`
	EncryptedPassword  string    `json:"-"`
}

func (i *Iot) Validate() error {
	if err := validation.ValidateStruct(
		i,
		validation.Field(&i.LastUpdateTimeUnix, validation.Min(0)),
		validation.Field(&i.IotState, validation.Required),
		validation.Field(&i.IotType, validation.Required),
	); err != nil {
		return err
	}

	if err := i.ValidateLongitudeAndLatitude(); err != nil {
		return err
	}

	if err := i.ValidateState(); err != nil {
		return err
	}

	return i.ValidateType()
}

func (i *Iot) ValidateLongitudeAndLatitude() error {
	sLongitude := strconv.FormatFloat(i.Longitude, 'f', 6, 64)
	fLongitude, err := strconv.ParseFloat(sLongitude, 64)

	if err != nil {
		return err
	}

	if fLongitude < -180 && fLongitude > 180 {
		return ErrIncorrectLongOrLang
	}

	sLatitude := strconv.FormatFloat(i.Latitude, 'f', 6, 64)
	fLatitude, err := strconv.ParseFloat(sLatitude, 64)

	if fLatitude < -90 && fLatitude > 90 {
		return ErrIncorrectLongOrLang
	}

	return nil
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

func (i *Iot) BeforeCreate() error {
	if len(i.Password) > 0 {
		enc, err := encryptString(i.Password)
		if err != nil {
			return err
		}

		i.EncryptedPassword = enc
	}

	sLongitude := strconv.FormatFloat(i.Longitude, 'f', 6, 64)
	fLongitude, err := strconv.ParseFloat(sLongitude, 64)

	if err != nil {
		return err
	}

	i.Longitude = fLongitude

	sLatitude := strconv.FormatFloat(i.Latitude, 'f', 6, 64)
	fLatitude, err := strconv.ParseFloat(sLatitude, 64)
	if err != nil {
		return err
	}

	i.Latitude = fLatitude

	return nil
}

func (i *Iot) Sanitize() {
	i.Password = ""
}

func (i *Iot) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(i.EncryptedPassword), []byte(password)) == nil
}
