package apiserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	tokenHandler "github.com/RudeGalaxy1010/jwt-test-task/internal/app/tokenhandler"
	"github.com/RudeGalaxy1010/jwt-test-task/internal/model"
	"github.com/RudeGalaxy1010/jwt-test-task/internal/store"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	router     *mux.Router
	store      store.Store
	jwtHandler tokenHandler.JwtHandler
}

type TokensPair struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func NewServer(store store.Store, jwtHandler tokenHandler.JwtHandler) *Server {
	server := &Server{
		router:     mux.NewRouter(),
		store:      store,
		jwtHandler: jwtHandler,
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

		user, err := server.store.User().Find(guid)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrUserNotFound)
			return
		}

		ipAddress := strings.Split(r.RemoteAddr, ":")[0]
		refreshToken := server.jwtHandler.NewRefreshToken()
		base64RefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))
		access, err := server.jwtHandler.NewAccessToken(user, ipAddress, refreshToken)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrTokenCreationFailed)
			return
		}

		hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)

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

		server.Respond(w, r, http.StatusOK, TokensPair{
			Access:  access,
			Refresh: base64RefreshToken,
		})
	}
}

func (server *Server) HandleJwtRefresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &TokensPair{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.Error(w, r, http.StatusBadRequest, err)
			return
		}

		userClaims, err := server.jwtHandler.Decode(req.Access)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, ErrTokenInvalid)
		}

		decodedRefreshToken, err := base64.StdEncoding.DecodeString(req.Refresh)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, ErrTokenInvalid)
			return
		}

		err = server.jwtHandler.ValidateRefresh(req.Access, string(decodedRefreshToken))

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, ErrTokenInvalid)
			return
		}

		user, err := server.store.User().Find(userClaims.Id)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, ErrUserNotFound)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Refresh), decodedRefreshToken)

		if err != nil {
			server.Error(w, r, http.StatusBadRequest, ErrTokenInvalid)
			return
		}

		ipAddress := strings.Split(r.RemoteAddr, ":")[0]

		if userClaims.IpAddress != ipAddress {
			// TODO: Send email
			fmt.Printf("ip diff detected, expected: %s , was: %s", userClaims.IpAddress, ipAddress)
			fmt.Println()
		}

		refreshToken := uuid.New().String()
		base64RefreshToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))
		access, err := server.jwtHandler.NewAccessToken(user, ipAddress, refreshToken)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrTokenCreationFailed)
			return
		}

		newHashedRefresh, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)

		if err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrTokenCreationFailed)
			return
		}

		user.Refresh = string(newHashedRefresh)
		if err := server.store.User().UpdateRefreshToken(user); err != nil {
			server.Error(w, r, http.StatusUnprocessableEntity, ErrUserUpdateFailed)
			return
		}

		server.Respond(w, r, http.StatusOK, TokensPair{
			Access:  access,
			Refresh: base64RefreshToken,
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
