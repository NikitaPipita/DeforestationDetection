package model

type User struct {
	ID                int
	Email             string
	EncryptedPassword string
	Role              string
	FullName          string
}
