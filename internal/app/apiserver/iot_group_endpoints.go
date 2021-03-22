package apiserver

import (
	"deforestation.detection.com/server/internal/app/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *server) getAllIotGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iotGroups, err := s.store.IotGroup().GetAll()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, iotGroups)
	}
}

func (s *server) getIotGroupByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errIncorrectID)
			return
		}

		iotGroup, err := s.store.IotGroup().FindByID(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, iotGroup)
	}
}

func (s *server) createIotGroup() http.HandlerFunc {
	type request struct {
		UserID                 int `json:"user_id"`
		UpdateDurationSeconds  int `json:"update_duration_seconds"`
		LastIotChangesTimeUnix int `json:"last_iot_changes_time_unix"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		group := &model.IotGroup{
			User: &model.User{
				ID: req.UserID,
			},
			UpdateDurationSeconds:  req.UpdateDurationSeconds,
			LastIotChangesTimeUnix: req.LastIotChangesTimeUnix,
		}
		if err := s.store.IotGroup().Create(group); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, group)
	}
}

func (s *server) createIotGroupByUser() http.HandlerFunc {
	type request struct {
		UserID int `json:"user_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		group := &model.IotGroup{
			User: &model.User{
				ID: req.UserID,
			},
		}
		if err := s.store.IotGroup().CreateByUser(group); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, group)
	}
}

func (s *server) updateIotGroupById() http.HandlerFunc {
	type request struct {
		UpdateDurationSeconds  int `json:"update_duration_seconds"`
		LastIotChangesTimeUnix int `json:"last_iot_changes_time_unix"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errIncorrectID)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		group := &model.IotGroup{
			User:                   &model.User{},
			UpdateDurationSeconds:  req.UpdateDurationSeconds,
			LastIotChangesTimeUnix: req.LastIotChangesTimeUnix,
		}
		if err := s.store.IotGroup().Update(id, group); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, group)
	}
}

func (s *server) deleteIotGroupById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errIncorrectID)
			return
		}

		if err := s.store.IotGroup().Delete(id); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusNoContent, nil)
	}
}
