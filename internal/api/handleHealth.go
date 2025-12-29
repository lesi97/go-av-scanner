package api

import (
	"context"
	"net/http"

	"github.com/lesi97/go-av-scanner/internal/store"
	"github.com/lesi97/go-av-scanner/internal/utils"
)

func (h *ApiHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {

	ctx := context.WithValue(
		r.Context(),
		store.ContextKey,
		store.Context{},
	)

	message, err := h.apiStore.Health(ctx)
	if err != nil {
		h.logger.Printf("ERROR in health handler: %v\n", err)
		utils.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.Success(w, http.StatusOK, message)
}