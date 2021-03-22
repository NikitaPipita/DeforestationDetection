package sqlstore

import (
	"database/sql"
	"deforestation.detection.com/server/internal/app/store"

	_ "github.com/lib/pq"
)

type Store struct {
	db                 *sql.DB
	userRepository     *UserRepository
	iotGroupRepository *IotGroupRepository
	iotRepository      *IotRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) IotGroup() store.IotGroupRepository {
	if s.iotGroupRepository != nil {
		return s.iotGroupRepository
	}

	s.iotGroupRepository = &IotGroupRepository{
		store: s,
	}

	return s.iotGroupRepository
}

func (s *Store) Iot() store.IotRepository {
	if s.iotRepository != nil {
		return s.iotRepository
	}

	s.iotRepository = &IotRepository{
		store: s,
	}

	return s.iotRepository
}
