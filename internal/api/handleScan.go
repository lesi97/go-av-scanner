package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lesi97/go-av-scanner/internal/scanner"
	"github.com/lesi97/go-av-scanner/internal/utils"
)

func (h *ApiHandler) HandleScan(w http.ResponseWriter, r *http.Request) {
	// ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	// defer cancel()

	r.Body = http.MaxBytesReader(w, r.Body, h.apiStore.MaxUploadBytes())
	
	err := r.ParseMultipartForm(64 << 24) // 1073741824 bytes
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			maxFileSize := utils.FormatBytes(h.apiStore.MaxUploadBytes())
			newErr := fmt.Errorf("file size too large, max filesize is: %v", maxFileSize)
			utils.Error(w, http.StatusRequestEntityTooLarge, newErr)
			return
		}

		utils.Error(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	file, _, err := r.FormFile("file")
	if err == nil {
		defer func() { _ = file.Close() }()

		// message, err := h.apiStore.Scan(ctx, file)
		message, err := h.apiStore.Scan(r.Context(), file)
		if err != nil {
			var scanErr *scanner.ScanError
			if errors.As(err, &scanErr) {
				utils.Error(w, http.StatusUnprocessableEntity, scanErr.Result)
				return
			}
			utils.Error(w, http.StatusInternalServerError, err)
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

	// message, err := h.apiStore.Scan(ctx, strings.NewReader(content))
	message, err := h.apiStore.Scan(r.Context(), strings.NewReader(content))
	if err != nil {
		var scanErr *scanner.ScanError
		if errors.As(err, &scanErr) {
			utils.Error(w, http.StatusUnprocessableEntity, scanErr.Result)
			return
		}
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.Success(w, http.StatusOK, message)

}