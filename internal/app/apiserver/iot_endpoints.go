package apiserver

import (
	"deforestation.detection.com/server/internal/app/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *server) getAllIots() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iots, err := s.store.Iot().GetAll()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, iots)
	}
}

func (s *server) getAllIotsInGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errIncorrectID)
			return
		}

		iots, err := s.store.Iot().FindAllInGroup(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, iots)
	}
}

func (s *server) getIotById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errIncorrectID)
			return
		}

		iot, err := s.store.Iot().FindByID(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, iot)
	}
}

func (s *server) createIot() http.HandlerFunc {
	type request struct {
		UserID             int     `json:"user_id"`
		GroupID            int     `json:"group_id"`
		Longitude          float64 `json:"longitude"`
		Latitude           float64 `json:"latitude"`
		LastUpdateTimeUnix int64   `json:"last_update_time_unix"`
		IotState           string  `json:"iot_state"`
		IotType            string  `json:"iot_type"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		iot := &model.Iot{
			User: &model.User{
				ID: req.UserID,
			},
			Group: &model.IotGroup{
				ID: req.GroupID,
			},
			Longitude:          req.Longitude,
			Latitude:           req.Latitude,
			LastUpdateTimeUnix: req.LastUpdateTimeUnix,
			IotState:           req.IotState,
			IotType:            req.IotType,
		}
		if err := s.store.Iot().Create(iot); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, iot)
	}
}

func (s *server) createIotByUser() http.HandlerFunc {
	type request struct {
		UserID    int     `json:"user_id"`
		GroupID   int     `json:"group_id"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		IotState  string  `json:"iot_state"`
		IotType   string  `json:"iot_type"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		iot := &model.Iot{
			User: &model.User{
				ID: req.UserID,
			},
			Group: &model.IotGroup{
				ID: req.GroupID,
			},
			Longitude: req.Longitude,
			Latitude:  req.Latitude,
			IotState:  req.IotState,
			IotType:   req.IotType,
		}
		if err := s.store.Iot().CreateByUser(iot); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, iot)
	}
}

func (s *server) updateIotById() http.HandlerFunc {
	type request struct {
		Longitude          float64 `json:"longitude"`
		Latitude           float64 `json:"latitude"`
		LastUpdateTimeUnix int64   `json:"last_update_time_unix"`
		IotState           string  `json:"iot_state"`
		IotType            string  `json:"iot_type"`
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

		iot := &model.Iot{
			User:               &model.User{},
			Group:              &model.IotGroup{},
			Longitude:          req.Longitude,
			Latitude:           req.Latitude,
			LastUpdateTimeUnix: req.LastUpdateTimeUnix,
			IotState:           req.IotState,
			IotType:            req.IotType,
		}
		if err := s.store.Iot().Update(id, iot); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, iot)
	}
}

func (s *server) deleteIotById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errIncorrectID)
			return
		}

		if err := s.store.Iot().Delete(id); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusNoContent, nil)
	}
}
