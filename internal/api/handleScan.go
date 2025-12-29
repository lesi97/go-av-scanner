package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/lesi97/go-av-scanner/internal/scanner"
	"github.com/lesi97/go-av-scanner/internal/utils"
)

func (h *ApiHandler) HandleScan(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	file, _, err := r.FormFile("file")
	if err == nil {
		defer func() { _ = file.Close() }()

		message, err := h.apiStore.Scan(ctx, file)
		if err != nil {
			var scanErr *scanner.ScanError
			if errors.As(err, &scanErr) {
				utils.Error(w, http.StatusUnprocessableEntity, scanErr.Result)
				return
			}
			utils.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}
		utils.Success(w, http.StatusOK, message)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		utils.Error(w, http.StatusBadRequest, "missing file field")
		return
	}

	message, err := h.apiStore.Scan(ctx, strings.NewReader(content))
	if err != nil {
		var scanErr *scanner.ScanError
		if errors.As(err, &scanErr) {
			utils.Error(w, http.StatusUnprocessableEntity, scanErr.Result)
			return
		}
		utils.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.Success(w, http.StatusOK, message)

}