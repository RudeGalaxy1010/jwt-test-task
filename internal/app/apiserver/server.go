package apiserver

import (
	"encoding/base64"
	"encoding/json"
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
			server.Error(w, r, http.StatusBadRequest, ErrInsufficientRequest)
			return
		}

		ipAddress := strings.Split(r.RemoteAddr, ":")[0]

		user, err := server.store.User().Find(guid)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrUserNotFound)
			return
		}

		uuid := uuid.New().String()
		access, err := jwt.NewAccessToken(user, uuid, ipAddress, server.key)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrTokenCreationFailed)
			return
		}

		hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(uuid), bcrypt.DefaultCost)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrTokenCreationFailed)
			return
		}

		if user == nil {
			user := &model.User{
				Id:      guid,
				Refresh: string(hashedRefresh),
			}

			if err := server.store.User().Create(user); err != nil {
				server.Error(w, r, http.StatusUnprocessableEntity, ErrUserCreationFailed)
				return
			}
		} else {
			user.Refresh = string(hashedRefresh)
			if err := server.store.User().UpdateRefreshToken(user); err != nil {
				server.Error(w, r, http.StatusUnprocessableEntity, ErrUserUpdateFailed)
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
	return func(w http.ResponseWriter, r *http.Request) {
		req := &jwt.TokensPair{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.Error(w, r, http.StatusBadRequest, err)
			return
		}

		userId, ip, err := jwt.DecodeAccess(req.Access, server.key)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, ErrTokenInvalid)
		}

		user, err := server.store.User().Find(userId)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, ErrUserNotFound)
			return
		}

		decodedRefresh, err := base64.StdEncoding.DecodeString(req.Refresh)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, ErrTokenInvalid)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Refresh), decodedRefresh)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, ErrTokenInvalid)
			return
		}

		err = jwt.ValidateRefresh(req.Access, string(decodedRefresh), server.key)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, ErrTokenInvalid)
			return
		}

		ipAddress := strings.Split(r.RemoteAddr, ":")[0]

		if ip != ipAddress {
			// TODO: Send email
			fmt.Printf("ip diff detected, expected: %s , was: %s", ipAddress, ip)
			fmt.Println()
		}

		uuid := uuid.New().String()
		access, err := jwt.NewAccessToken(user, uuid, ipAddress, server.key)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrTokenCreationFailed)
			return
		}

		newHashedRefresh, err := bcrypt.GenerateFromPassword([]byte(uuid), bcrypt.DefaultCost)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrTokenCreationFailed)
			return
		}

		user.Refresh = string(newHashedRefresh)
		if err := server.store.User().UpdateRefreshToken(user); err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrUserUpdateFailed)
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
