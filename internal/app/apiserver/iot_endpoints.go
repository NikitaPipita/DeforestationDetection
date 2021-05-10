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
		enableCors(&w)
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
		enableCors(&w)
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
		enableCors(&w)
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
		enableCors(&w)
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
			Password:           "qwerty",
		}
		if err := s.store.Iot().Create(iot); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		iot.Sanitize()

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
		enableCors(&w)
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
			Password:  "qwerty",
		}
		if err := s.store.Iot().CreateByUser(iot); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		iot.Sanitize()

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
		enableCors(&w)
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
		enableCors(&w)
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

func (s *server) checkIfPositionSuitable() http.HandlerFunc {
	type request struct {
		GroupID   int     `json:"group_id"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		IotType   string  `json:"iot_type"`
	}

	type response struct {
		Suitable                  bool    `json:"suitable"`
		MinimumDistanceToMoveAway float64 `json:"minimum_distance_to_move_away"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		suitable, minimumDistanceToMoveAway, err := s.store.Iot().CheckIfPositionSuitable(req.GroupID, req.Longitude, req.Latitude, req.IotType)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		res := response{
			Suitable:                  suitable,
			MinimumDistanceToMoveAway: minimumDistanceToMoveAway,
		}

		s.respond(w, r, http.StatusOK, res)
	}
}

func (s *server) getAllSignaling() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		iots, err := s.store.Iot().GetAllSignaling()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, iots)
	}
}

func (s *server) changeIotState() http.HandlerFunc {
	type request struct {
		ID       int    `json:"iot_id"`
		IotState string `json:"iot_state"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Iot().ChangeState(req.ID, req.IotState); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleIotSessionsCreate() http.HandlerFunc {
	type request struct {
		ID       int    `json:"serial_id"`
		Password string `json:"serial_password"`
	}

	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		iot, err := s.store.Iot().FindPasswordByID(req.ID)
		if err != nil || !iot.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		ts, err := CreateIotToken(uint64(iot.ID))
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := CreateAuth(uint64(iot.ID), ts); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		res := response{
			AccessToken:  ts.AccessToken,
			RefreshToken: ts.RefreshToken}

		s.respond(w, r, http.StatusOK, res)
	}
}

func (s *server) changeConnectedIotState() http.HandlerFunc {
	type request struct {
		ID       int    `json:"iot_id"`
		IotState string `json:"iot_state"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Iot().ChangeState(req.ID, req.IotState); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}
