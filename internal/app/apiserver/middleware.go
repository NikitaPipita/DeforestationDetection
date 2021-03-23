package apiserver

import (
	"log"
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

		log.Println(userId)

		next(w, r)
	}
}
