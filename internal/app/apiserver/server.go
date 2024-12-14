package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/RudeGalaxy1010/jwt-test-task/internal/model"
	"github.com/RudeGalaxy1010/jwt-test-task/internal/store"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	store  store.Store
}

func NewServer(store store.Store) *Server {
	server := &Server{
		router: mux.NewRouter(),
		store:  store,
	}

	server.ConfigureRouter()

	return server
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

func (server *Server) ConfigureRouter() {
	server.router.HandleFunc("/ping", server.HandlePing()).Methods("GET")
	server.router.HandleFunc("/jwt", server.HandleJwtGet()).Methods("POST")
}

func (server *Server) HandlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.Respond(w, r, http.StatusOK, "hello!")
	}
}

func (server *Server) HandleJwtGet() http.HandlerFunc {
	type request struct {
		Guid string `json:"guid"`
	}

	type response struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.Error(w, r, http.StatusBadRequest, err)
			return
		}

		guid := strings.TrimSpace(req.Guid)

		if len(guid) == 0 {
			server.Error(w, r, http.StatusBadRequest, errors.New("guid can't be empty"))
			return
		}

		ipAddress := strings.Split(r.RemoteAddr, ":")[0]

		user, err := server.store.User().Find(guid)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, errors.New("failed to find user"))
		}

		if user == nil {
			user := &model.User{
				Id:        guid,
				IpAddress: ipAddress,
			}

			if err := server.store.User().Create(user); err != nil {
				server.Error(w, r, http.StatusUnprocessableEntity, errors.New("failed to create user"))
				return
			}
		}

		server.Respond(w, r, http.StatusOK, response{
			Access:  "",
			Refresh: "",
		})
	}
}

func (server *Server) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	server.Respond(w, r, code, map[string]string{"error": err.Error()})
}

func (server *Server) Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
