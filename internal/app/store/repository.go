package store

import "deforestation.detection.com/server/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	GetAll() ([]model.User, error)
	FindByID(int) (*model.User, error)
	FindByIDWithPassword(id int) (*model.User, error)
	Update(int, *model.User) error
	Delete(int) error
}

type IotGroupRepository interface {
	GetAll() ([]model.IotGroup, error)
	FindByID(int) (*model.IotGroup, error)
	Create(*model.IotGroup) error
	CreateByUser(*model.IotGroup) error
	Update(int, *model.IotGroup) error
	Delete(int) error
}