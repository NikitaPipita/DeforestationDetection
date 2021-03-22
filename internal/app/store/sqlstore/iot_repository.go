package sqlstore

import (
	"database/sql"
	"deforestation.detection.com/server/internal/app/model"
	"deforestation.detection.com/server/internal/app/store"
)

type IotRepository struct {
	store *Store
}

func (r *IotRepository) GetAll() ([]model.Iot, error) {
	rows, err := r.store.db.Query(
		"SELECT iot_id, system_user.user_id, email, user_role, full_name, " +
			"iot_group.group_id, update_duration_seconds, last_iot_changes_time_unix, " +
			"longitude, latitude, last_update_time_unix, iot_state, iot_type " +
			"FROM iot " +
			"JOIN system_user " +
			"ON iot.user_id = system_user.user_id " +
			"JOIN iot_group " +
			"ON iot.group_id = iot_group.group_id")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	defer rows.Close()

	var iots []model.Iot

	for rows.Next() {
		var i model.Iot
		i.User = &model.User{}
		i.Group = &model.IotGroup{}
		if err := rows.Scan(
			&i.ID,
			&i.User.ID,
			&i.User.Email,
			&i.User.Role,
			&i.User.FullName,
			&i.Group.ID,
			&i.Group.UpdateDurationSeconds,
			&i.Group.LastIotChangesTimeUnix,
			&i.Longitude,
			&i.Latitude,
			&i.LastUpdateTimeUnix,
			&i.IotState,
			&i.IotType,
		); err != nil {
			return nil, err
		}
		iots = append(iots, i)
	}

	return iots, nil
}

func (r *IotRepository) FindAllInGroup(groupID int) ([]model.Iot, error) {
	rows, err := r.store.db.Query(
		"SELECT iot_id, system_user.user_id, email, user_role, full_name, "+
			"iot_group.group_id, update_duration_seconds, last_iot_changes_time_unix, "+
			"longitude, latitude, last_update_time_unix, iot_state, iot_type "+
			"FROM iot "+
			"JOIN system_user "+
			"ON iot.user_id = system_user.user_id "+
			"JOIN iot_group "+
			"ON iot.group_id = iot_group.group_id "+
			"WHERE iot_group.group_id = $1", groupID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	defer rows.Close()

	var iots []model.Iot

	for rows.Next() {
		var i model.Iot
		i.User = &model.User{}
		i.Group = &model.IotGroup{}
		if err := rows.Scan(
			&i.ID,
			&i.User.ID,
			&i.User.Email,
			&i.User.Role,
			&i.User.FullName,
			&i.Group.ID,
			&i.Group.UpdateDurationSeconds,
			&i.Group.LastIotChangesTimeUnix,
			&i.Longitude,
			&i.Latitude,
			&i.LastUpdateTimeUnix,
			&i.IotState,
			&i.IotType,
		); err != nil {
			return nil, err
		}
		iots = append(iots, i)
	}

	return iots, nil
}

func (r *IotRepository) FindByID(id int) (*model.Iot, error) {
	i := &model.Iot{}
	i.User = &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT iot_id, system_user.user_id, email, user_role, full_name, "+
			"iot_group.group_id, update_duration_seconds, last_iot_changes_time_unix, "+
			"longitude, latitude, last_update_time_unix, iot_state, iot_type "+
			"FROM iot "+
			"JOIN system_user "+
			"ON iot.user_id = system_user.user_id "+
			"JOIN iot_group "+
			"ON iot.group_id = iot_group.group_id "+
			"WHERE iot_id = $1",
		id,
	).Scan(
		&i.ID,
		&i.User.ID,
		&i.User.Email,
		&i.User.Role,
		&i.User.FullName,
		&i.Group.ID,
		&i.Group.UpdateDurationSeconds,
		&i.Group.LastIotChangesTimeUnix,
		&i.Longitude,
		&i.Latitude,
		&i.LastUpdateTimeUnix,
		&i.IotState,
		&i.IotType,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return i, nil
}

func (r *IotRepository) Create(i *model.Iot) error {
	if err := i.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO iot (user_id, group_id, longitude, latitude, last_update_time_unix, iot_state, iot_type) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING iot_id",
		i.User.ID,
		i.Group.ID,
		i.Longitude,
		i.Latitude,
		i.LastUpdateTimeUnix,
		i.IotState,
		i.IotType,
	).Scan(&i.ID)
}

func (r *IotRepository) CreateByUser(i *model.Iot) error {
	if err := i.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO iot (user_id, group_id, longitude, latitude, iot_state, iot_type) "+
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING iot_id",
		i.User.ID,
		i.Group.ID,
		i.Longitude,
		i.Latitude,
		i.IotState,
		i.IotType,
	).Scan(
		&i.ID,
	)
}

func (r *IotRepository) Update(id int, i *model.Iot) error {
	if err := i.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"UPDATE iot SET longitude = $1, latitude = $2 , last_update_time_unix = $3, iot_state = $4, iot_type = $5"+
			"WHERE iot_id = $6 RETURNING iot_id, user_id, group_id, longitude, latitude, last_update_time_unix, iot_state, iot_type",
		i.Longitude,
		i.Latitude,
		i.LastUpdateTimeUnix,
		i.IotState,
		i.IotType,
		id,
	).Scan(
		&i.ID,
		&i.User.ID,
		&i.Group.ID,
		&i.Longitude,
		&i.Latitude,
		&i.LastUpdateTimeUnix,
		&i.IotState,
		&i.IotType,
	)
}

func (r *IotRepository) Delete(id int) error {
	_, err := r.store.db.Exec("DELETE FROM iot WHERE iot_id = $1", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}

		return err
	}

	return nil
}
