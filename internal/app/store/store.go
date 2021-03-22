package store

type Store interface {
	User() UserRepository
	IotGroup() IotGroupRepository
}
