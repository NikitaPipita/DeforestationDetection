package apiserver

import (
	"deforestation.detection.com/server/internal/app/store"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errIncorrectID              = errors.New("incorrect id")
)

type server struct {
	router *mux.Router
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/users", s.getAllUsers()).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{id:[0-9]+}", s.getUserById()).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{id:[0-9]+}/info", s.getUserByIdWithPassword()).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{id:[0-9]+}", s.updateUserById()).Methods(http.MethodPut)
	s.router.HandleFunc("/user/{id:[0-9]+}", s.deleteUserById()).Methods(http.MethodDelete)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
