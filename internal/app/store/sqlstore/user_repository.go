package sqlstore

import (
	"database/sql"
	"deforestation.detection.com/server/internal/app/model"
	"deforestation.detection.com/server/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO system_user (email, password, user_role, full_name) VALUES ($1, $2, $3, $4) RETURNING user_id",
		u.Email,
		u.EncryptedPassword,
		u.Role,
		u.FullName,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT user_id, email, password, user_role, full_name FROM system_user WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.Role,
		&u.FullName,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	rows, err := r.store.db.Query("SELECT user_id, email, user_role, full_name FROM system_user")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Role,
			&u.FullName,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) FindByID(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT user_id, email, user_role, full_name FROM system_user WHERE user_id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Role,
		&u.FullName,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByIDWithPassword(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT user_id, email, password, user_role, full_name FROM system_user WHERE user_id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.FullName,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) Update(id int, u *model.User) error {
	if err := u.UpdateValidate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"UPDATE system_user SET email = $1, user_role = $2, full_name = $3 WHERE user_id = $4 RETURNING user_id, email, user_role, full_name",
		u.Email,
		u.Role,
		u.FullName,
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Role,
		&u.FullName,
	)
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.store.db.Exec("DELETE FROM system_user WHERE user_id = $1", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}

		return err
	}

	return nil
}
