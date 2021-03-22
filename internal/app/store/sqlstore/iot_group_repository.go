package sqlstore

import (
	"database/sql"
	"deforestation.detection.com/server/internal/app/model"
	"deforestation.detection.com/server/internal/app/store"
)

type IotGroupRepository struct {
	store *Store
}

func (r *IotGroupRepository) GetAll() ([]model.IotGroup, error) {
	rows, err := r.store.db.Query(
		"SELECT group_id, system_user.user_id, email, user_role, full_name, update_duration_seconds, last_iot_changes_time_unix " +
			"FROM iot_group " +
			"JOIN system_user " +
			"ON iot_group.user_id = system_user.user_id")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	defer rows.Close()

	var groups []model.IotGroup

	for rows.Next() {
		var g model.IotGroup
		g.User = &model.User{}
		if err := rows.Scan(
			&g.ID,
			&g.User.ID,
			&g.User.Email,
			&g.User.Role,
			&g.User.FullName,
			&g.UpdateDurationSeconds,
			&g.LastIotChangesTimeUnix,
		); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

func (r *IotGroupRepository) FindByID(id int) (*model.IotGroup, error) {
	g := &model.IotGroup{}
	g.User = &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT group_id, system_user.user_id, email, user_role, full_name, update_duration_seconds, last_iot_changes_time_unix "+
			"FROM iot_group "+
			"JOIN system_user "+
			"ON iot_group.user_id = system_user.user_id "+
			"WHERE group_id = $1",
		id,
	).Scan(
		&g.ID,
		&g.User.ID,
		&g.User.Email,
		&g.User.Role,
		&g.User.FullName,
		&g.UpdateDurationSeconds,
		&g.LastIotChangesTimeUnix,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return g, nil
}

func (r *IotGroupRepository) Create(g *model.IotGroup) error {
	if err := g.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO iot_group (user_id, update_duration_seconds, last_iot_changes_time_unix) VALUES ($1, $2, $3) RETURNING group_id",
		g.User.ID,
		g.UpdateDurationSeconds,
		g.LastIotChangesTimeUnix,
	).Scan(&g.ID)
}

func (r *IotGroupRepository) CreateByUser(g *model.IotGroup) error {
	return r.store.db.QueryRow(
		"INSERT INTO iot_group (user_id) VALUES ($1) RETURNING group_id",
		g.User.ID,
	).Scan(
		&g.ID,
	)
}

func (r *IotGroupRepository) Update(id int, g *model.IotGroup) error {
	if err := g.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"UPDATE iot_group SET update_duration_seconds = $1, last_iot_changes_time_unix = $2 "+
			"WHERE group_id = $3 RETURNING group_id, user_id, update_duration_seconds, last_iot_changes_time_unix",
		g.UpdateDurationSeconds,
		g.LastIotChangesTimeUnix,
		id,
	).Scan(
		&g.ID,
		&g.User.ID,
		&g.UpdateDurationSeconds,
		&g.LastIotChangesTimeUnix,
	)
}

func (r *IotGroupRepository) Delete(id int) error {
	_, err := r.store.db.Exec("DELETE FROM iot_group WHERE group_id = $1", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}

		return err
	}

	return nil
}
