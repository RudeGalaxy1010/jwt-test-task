package apiserver

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/RudeGalaxy1010/jwt-test-task/internal/app"
	"github.com/RudeGalaxy1010/jwt-test-task/internal/model"
	"github.com/RudeGalaxy1010/jwt-test-task/internal/store"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	router *mux.Router
	store  store.Store
	key    string
}

func NewServer(store store.Store, key string) *Server {
	server := &Server{
		router: mux.NewRouter(),
		store:  store,
		key:    key,
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
	server.router.HandleFunc("/jwt-refresh", server.HandleJwtRefresh()).Methods("POST")
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
			return
		}

		uuid := uuid.New().String()
		access, err := jwt.NewAccessToken(user, uuid, server.key)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, errors.New("failed to create access token"))
			return
		}

		hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(uuid), bcrypt.DefaultCost)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, errors.New("failed to create refresh token"))
			return
		}

		if user == nil {
			user := &model.User{
				Id:        guid,
				IpAddress: ipAddress,
				Refresh:   string(hashedRefresh),
			}

			if err := server.store.User().Create(user); err != nil {
				server.Error(w, r, http.StatusUnprocessableEntity, errors.New("failed to create user"))
				return
			}
		} else {
			user.Refresh = string(hashedRefresh)
			if err := server.store.User().UpdateRefreshToken(user); err != nil {
				server.Error(w, r, http.StatusUnprocessableEntity, errors.New("failed to update user"))
				return
			}
		}

		server.Respond(w, r, http.StatusOK, jwt.TokensPair{
			Access:  access,
			Refresh: base64.StdEncoding.EncodeToString([]byte(uuid)),
		})
	}
}

func (server *Server) HandleJwtRefresh() http.HandlerFunc {
	type request struct {
		Guid       string         `json:"guid"`
		TokensPair jwt.TokensPair `json:"tokens"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.Error(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := server.store.User().Find(req.Guid)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, errors.New("user not found"))
			return
		}

		decodedRefresh, err := base64.StdEncoding.DecodeString(req.TokensPair.Refresh)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, errors.New("invalid token"))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Refresh), decodedRefresh)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, errors.New("invalid token"))
			return
		}

		ip, err := jwt.ValidateRefresh(req.TokensPair.Access, string(decodedRefresh), server.key)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, errors.New("invalid token"))
			return
		}

		if ip != user.IpAddress {
			// TODO: Send email
			fmt.Printf("ip diff detected, expected: %s , was: %s", user.IpAddress, ip)
			fmt.Println()
		}

		uuid := uuid.New().String()
		access, err := jwt.NewAccessToken(user, uuid, server.key)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, errors.New("failed to create access token"))
			return
		}

		newHashedRefresh, err := bcrypt.GenerateFromPassword([]byte(uuid), bcrypt.DefaultCost)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, errors.New("failed to create refresh token"))
			return
		}

		user.Refresh = string(newHashedRefresh)
		if err := server.store.User().UpdateRefreshToken(user); err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, errors.New("failed to update user"))
			return
		}

		server.Respond(w, r, http.StatusOK, jwt.TokensPair{
			Access:  access,
			Refresh: base64.StdEncoding.EncodeToString([]byte(uuid)),
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
