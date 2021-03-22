package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int    `json:"user_id"`
	Email             string `json:"email,omitempty"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
	Role              string `json:"user_role,omitempty"`
	FullName          string `json:"full_name,omitempty"`
}

func (u *User) Validate() error {
	if err := validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
		validation.Field(&u.Role, validation.Required),
		validation.Field(&u.FullName, validation.Required),
	); err != nil {
		return err
	}

	return u.ValidateRole()
}

func (u *User) UpdateValidate() error {
	if err := validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Role, validation.Required),
		validation.Field(&u.FullName, validation.Required),
	); err != nil {
		return err
	}

	return u.ValidateRole()
}

func (u *User) ValidateRole() error {
	if u.Role == "admin" || u.Role == "manager" || u.Role == "employee" || u.Role == "observer" || u.Role == "locked" {
		return nil
	}

	return ErrIncorrectRole
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
