package store

type Store interface {
	User() UserRepository
	IotGroup() IotGroupRepository
	Iot() IotRepository
	Dump() DumpRepository
}
