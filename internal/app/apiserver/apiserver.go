package apiserver

import (
	"net/http"
)

func Start() error {
	server := NewServer()

	return http.ListenAndServe("localhost:5000", server)
}
