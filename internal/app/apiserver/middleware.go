package apiserver

import (
	"net/http"
)

func (s *server) tokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		err := TokenValid(r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		next(w, r)
	}
}

func (s *server) adminAccessMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenAuth, err := ExtractTokenMetadata(r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		userId, err := FetchAuth(tokenAuth)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		role, err := s.store.User().GetRole(int(userId))
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if role != "admin" {
			s.error(w, r, http.StatusUnauthorized, errAccessDenied)
			return
		}

		next(w, r)
	}
}

func (s *server) managerAccessMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenAuth, err := ExtractTokenMetadata(r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		userId, err := FetchAuth(tokenAuth)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		role, err := s.store.User().GetRole(int(userId))
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if role != "admin" && role != "manager" {
			s.error(w, r, http.StatusUnauthorized, errAccessDenied)
			return
		}

		next(w, r)
	}
}

func (s *server) employeeAccessMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenAuth, err := ExtractTokenMetadata(r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		userId, err := FetchAuth(tokenAuth)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		role, err := s.store.User().GetRole(int(userId))
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if role != "admin" && role != "manager" && role != "employee" {
			s.error(w, r, http.StatusUnauthorized, errAccessDenied)
			return
		}

		next(w, r)
	}
}

func (s *server) observerAccessMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenAuth, err := ExtractTokenMetadata(r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		userId, err := FetchAuth(tokenAuth)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		role, err := s.store.User().GetRole(int(userId))
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if role != "admin" && role != "manager" && role != "employee" && role != "observer" {
			s.error(w, r, http.StatusUnauthorized, errAccessDenied)
			return
		}

		next(w, r)
	}
}
