package store

import "deforestation.detection.com/server/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	GetAll() ([]model.User, error)
}
