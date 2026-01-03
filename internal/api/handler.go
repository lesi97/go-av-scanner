package api

import (
	"github.com/lesi97/go-av-scanner/internal/store"
	"github.com/lesi97/go-av-scanner/internal/utils"
)

type ApiHandler struct {
	logger          *utils.Logger
    apiStore      	store.ApiStore
}

func NewApiHandler(logger *utils.Logger, apiStore store.ApiStore)  *ApiHandler {
	return &ApiHandler{
		logger: logger,
        apiStore: apiStore,
	}
}
