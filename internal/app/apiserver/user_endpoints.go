package apiserver

import (
	"deforestation.detection.com/server/internal/app/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ID           int    `json:"user_id"`
		Email        string `json:"email"`
		Role         string `json:"user_role"`
		FullName     string `json:"full_name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !user.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		ts, err := CreateToken(uint64(user.ID))
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := CreateAuth(uint64(user.ID), ts); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		res := response{
			AccessToken:  ts.AccessToken,
			RefreshToken: ts.RefreshToken,
			ID:           user.ID,
			Email:        user.Email,
			Role:         user.Role,
			FullName:     user.FullName}

		s.respond(w, r, http.StatusOK, res)
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"user_role"`
		FullName string `json:"full_name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user := &model.User{
			Email:    req.Email,
			Password: req.Password,
			Role:     req.Role,
			FullName: req.FullName,
		}
		if err := s.store.User().Create(user); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		user.Sanitize()
		s.respond(w, r, http.StatusCreated, user)
	}
}

func (s *server) getAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		users, err := s.store.User().GetAll()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, users)
	}
}

func (s *server) getUserById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errIncorrectID)
			return
		}

		user, err := s.store.User().FindByID(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, user)
	}
}

func (s *server) getUserByIdWithPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errIncorrectID)
			return
		}

		user, err := s.store.User().FindByIDWithPassword(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, user)
	}
}

func (s *server) updateUserById() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Role     string `json:"user_role"`
		FullName string `json:"full_name"`
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

		user := &model.User{
			Email:    req.Email,
			Role:     req.Role,
			FullName: req.FullName,
		}
		if err := s.store.User().Update(id, user); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, user)
	}
}

func (s *server) deleteUserById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errIncorrectID)
			return
		}

		if err := s.store.User().Delete(id); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusNoContent, nil)
	}
}
