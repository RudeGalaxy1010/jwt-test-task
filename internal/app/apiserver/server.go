package apiserver

import (
	"encoding/json"
	"net/http"

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
}

func (server *Server) HandlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.Respond(w, r, http.StatusOK, "hello!")
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
