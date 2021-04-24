package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strconv"
)

func (s *server) updateToken() http.HandlerFunc {
	type request struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("REFRESH_SECRET")), nil
		})

		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			s.error(w, r, http.StatusUnauthorized, errUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {

			refreshUuid, ok := claims["refresh_uuid"].(string)
			if !ok {
				s.error(w, r, http.StatusUnprocessableEntity, errUnauthorized)
				return
			}

			userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			deleted, delErr := DeleteAuth(refreshUuid)
			if delErr != nil || deleted == 0 {
				s.error(w, r, http.StatusUnauthorized, errUnauthorized)
				return

			}

			ts, createErr := CreateToken(userId)
			if createErr != nil {
				s.error(w, r, http.StatusForbidden, createErr)
				return
			}

			saveErr := CreateAuth(userId, ts)
			if saveErr != nil {
				s.error(w, r, http.StatusForbidden, saveErr)
				return
			}

			tokens := map[string]string{
				"access_token":  ts.AccessToken,
				"refresh_token": ts.RefreshToken,
			}

			s.respond(w, r, http.StatusCreated, tokens)
		} else {
			s.error(w, r, http.StatusUnauthorized, errRefreshExpired)
		}
	}
}
