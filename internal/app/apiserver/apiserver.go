package apiserver

import (
	"net/http"
)

func Start(config *Config) error {
	server := NewServer()

	return http.ListenAndServe(config.Address, server)
}
