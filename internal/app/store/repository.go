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
	GetRole(id int) (string, error)
}

type IotGroupRepository interface {
	GetAll() ([]model.IotGroup, error)
	FindByID(int) (*model.IotGroup, error)
	Create(*model.IotGroup) error
	CreateByUser(*model.IotGroup) error
	Update(int, *model.IotGroup) error
	Delete(int) error
}

type IotRepository interface {
	GetAll() ([]model.Iot, error)
	FindAllInGroup(int) ([]model.Iot, error)
	FindByID(int) (*model.Iot, error)
	Create(*model.Iot) error
	CreateByUser(*model.Iot) error
	Update(int, *model.Iot) error
	Delete(int) error
	CheckIfPositionSuitable(groupID int, longitude float64, latitude float64, iotType string) (bool, error)
	GetAllSignaling() ([]model.Iot, error)
	ChangeState(int, string) error
}

type DumpRepository interface {
	CreateDump() string
	Execute(dumpingQuery string) error
}
