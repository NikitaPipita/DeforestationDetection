package apiserver

import (
	"deforestation.detection.com/server/internal/app/store"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errIncorrectID              = errors.New("incorrect id")
	errUnauthorized             = errors.New("unauthorized")
	errRefreshExpired           = errors.New("refresh expired")
	errAccessDenied             = errors.New("access denied")
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

	initRedis()

	s.configureRouter()

	return s
}

func initRedis() {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/refresh", s.updateToken()).Methods(http.MethodPost)
	s.router.HandleFunc("/users", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.handleUsersCreate()))).Methods(http.MethodPost)
	s.router.HandleFunc("/users", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.getAllUsers()))).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{id:[0-9]+}", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.getUserById()))).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{id:[0-9]+}/info", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.getUserByIdWithPassword()))).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{id:[0-9]+}", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.updateUserById()))).Methods(http.MethodPut)
	s.router.HandleFunc("/user/{id:[0-9]+}", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.deleteUserById()))).Methods(http.MethodDelete)

	s.router.HandleFunc("/groups", s.tokenAuthMiddleware(s.managerAccessMiddleware(s.getAllIotGroups()))).Methods(http.MethodGet)
	s.router.HandleFunc("/group/{id:[0-9]+}", s.tokenAuthMiddleware(s.managerAccessMiddleware(s.getIotGroupByID()))).Methods(http.MethodGet)
	s.router.HandleFunc("/groups", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.createIotGroup()))).Methods(http.MethodPost)
	s.router.HandleFunc("/groups/create", s.tokenAuthMiddleware(s.employeeAccessMiddleware(s.createIotGroupByUser()))).Methods(http.MethodPost)
	s.router.HandleFunc("/group/{id:[0-9]+}", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.updateIotGroupById()))).Methods(http.MethodPut)
	s.router.HandleFunc("/group/{id:[0-9]+}", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.deleteIotGroupById()))).Methods(http.MethodDelete)

	s.router.HandleFunc("/iots", s.tokenAuthMiddleware(s.managerAccessMiddleware(s.getAllIots()))).Methods(http.MethodGet)
	s.router.HandleFunc("/iots/{id:[0-9]+}", s.tokenAuthMiddleware(s.managerAccessMiddleware(s.getAllIotsInGroup()))).Methods(http.MethodGet)
	s.router.HandleFunc("/iot/{id:[0-9]+}", s.tokenAuthMiddleware(s.managerAccessMiddleware(s.getIotById()))).Methods(http.MethodGet)
	s.router.HandleFunc("/iots", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.createIot()))).Methods(http.MethodPost)
	s.router.HandleFunc("/iots/create", s.tokenAuthMiddleware(s.employeeAccessMiddleware(s.createIotByUser()))).Methods(http.MethodPost)
	s.router.HandleFunc("/iot/{id:[0-9]+}", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.updateIotById()))).Methods(http.MethodPut)
	s.router.HandleFunc("/iot/{id:[0-9]+}", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.deleteIotById()))).Methods(http.MethodDelete)
	s.router.HandleFunc("/iot/check", s.tokenAuthMiddleware(s.employeeAccessMiddleware(s.checkIfPositionSuitable()))).Methods(http.MethodPost)
	s.router.HandleFunc("/iot/signal", s.tokenAuthMiddleware(s.observerAccessMiddleware(s.getAllSignaling()))).Methods(http.MethodGet)
	s.router.HandleFunc("/iot/state", s.tokenAuthMiddleware(s.employeeAccessMiddleware(s.changeIotState()))).Methods(http.MethodPut)

	s.router.HandleFunc("/dump/make", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.MakeAndDownloadDump))).Methods(http.MethodGet)
	s.router.HandleFunc("/dump/exec", s.tokenAuthMiddleware(s.adminAccessMiddleware(s.UploadAndExecuteDump))).Methods(http.MethodPost)
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json; charset=utf-8")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, HEAD, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
