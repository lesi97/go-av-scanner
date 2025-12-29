package api

import (
	"log"

	"github.com/lesi97/go-av-scanner/internal/store"
)

type ApiHandler struct {
	logger          *log.Logger
    apiStore      	store.ApiStore
}

func NewApiHandler(logger *log.Logger, apiStore store.ApiStore)  *ApiHandler {
	return &ApiHandler{
		logger: logger,
        apiStore: apiStore,
	}
}
