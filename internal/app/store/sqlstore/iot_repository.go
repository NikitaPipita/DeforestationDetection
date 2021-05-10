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

func (r *IotRepository) FindPasswordByID(id int) (*model.Iot, error) {
	i := &model.Iot{}
	if err := r.store.db.QueryRow(
		"SELECT iot_id, iot_serial_password FROM iot WHERE iot_id = $1",
		id,
	).Scan(
		&i.ID,
		&i.EncryptedPassword,
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

	if err := i.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO iot (user_id, group_id, longitude, latitude, last_update_time_unix, iot_state, iot_type, iot_serial_password) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING iot_id",
		i.User.ID,
		i.Group.ID,
		i.Longitude,
		i.Latitude,
		i.LastUpdateTimeUnix,
		i.IotState,
		i.IotType,
		i.EncryptedPassword,
	).Scan(&i.ID)
}

func (r *IotRepository) CreateByUser(i *model.Iot) error {
	if err := i.Validate(); err != nil {
		return err
	}

	if err := i.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO iot (user_id, group_id, longitude, latitude, iot_state, iot_type, iot_serial_password) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING iot_id",
		i.User.ID,
		i.Group.ID,
		i.Longitude,
		i.Latitude,
		i.IotState,
		i.IotType,
		i.EncryptedPassword,
	).Scan(
		&i.ID,
	)
}

func (r *IotRepository) Update(id int, i *model.Iot) error {
	if err := i.Validate(); err != nil {
		return err
	}

	if err := i.BeforeCreate(); err != nil {
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

func (r *IotRepository) CheckIfPositionSuitable(groupID int, longitude float64, latitude float64, iotType string) (bool, float64, error) {
	i := &model.Iot{
		User:      &model.User{},
		Group:     &model.IotGroup{},
		Longitude: longitude,
		Latitude:  latitude,
		IotType:   iotType,
	}

	if err := i.ValidateLongitudeAndLatitude(); err != nil {
		return false, 0.0, err
	}

	if err := i.ValidateType(); err != nil {
		return false, 0.0, err
	}

	rows, err := r.store.db.Query(
		"SELECT longitude, latitude, iot_type "+
			"FROM iot "+
			"WHERE group_id = $1", groupID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0.0, store.ErrRecordNotFound
		}

		return false, 0.0, err
	}

	isSuitable := true
	minimumDistanceToMoveAway := 0.0

	for rows.Next() {
		var curLongitude float64
		var curLatitude float64
		var curType string
		if err := rows.Scan(
			&curLongitude,
			&curLatitude,
			&curType,
		); err != nil {
			return false, 0.0, err
		}

		result := haversineDistanceBetweenTwoPointsInMeters(longitude, latitude, curLongitude, curLatitude)
		distanceToMoveAway := 0.0

		if curType == "microphone" && iotType == "microphone" && result < 200 {
			isSuitable = false
			distanceToMoveAway = 200 - result
		} else if curType == "microphone" && iotType == "gyroscope" && result < 120 {
			isSuitable = false
			distanceToMoveAway = 120 - result
		} else if curType == "gyroscope" && iotType == "gyroscope" && result < 40 {
			isSuitable = false
			distanceToMoveAway = 40 - result
		} else if curType == "gyroscope" && iotType == "microphone" && result < 120 {
			isSuitable = false
			distanceToMoveAway = 120 - result
		}

		if distanceToMoveAway > 0 && distanceToMoveAway > minimumDistanceToMoveAway {
			minimumDistanceToMoveAway = distanceToMoveAway
		}
	}

	return isSuitable, minimumDistanceToMoveAway, nil
}

func (r *IotRepository) GetAllSignaling() ([]model.Iot, error) {
	rows, err := r.store.db.Query(
		"SELECT longitude, latitude, iot_state " +
			"FROM iot " +
			"WHERE iot_state = 'active' OR iot_state = 'lost'")

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
			&i.Longitude,
			&i.Latitude,
			&i.IotState,
		); err != nil {
			return nil, err
		}
		i.Longitude = i.Longitude + 0.000126
		i.Latitude = i.Latitude - 0.000077
		iots = append(iots, i)
	}

	return iots, nil
}

func (r *IotRepository) ChangeState(id int, state string) error {
	i := &model.Iot{
		User:     &model.User{},
		Group:    &model.IotGroup{},
		IotState: state,
	}

	if err := i.ValidateState(); err != nil {
		return err
	}

	_, err := r.store.db.Exec(
		"UPDATE iot SET iot_state = $1 WHERE iot_id = $2",
		i.IotState,
		id,
	)

	return err
}
