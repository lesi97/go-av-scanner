package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/lesi97/go-av-scanner/internal/utils"
)

func (h *ApiHandler) HandleScan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, h.apiStore.MaxUploadBytes())

	mr, err := r.MultipartReader()
	if err != nil {
		h.logger.Error(err)
		h.writeMultipartErr(w, err)
		return
	}

	var content string

	for {
		part, partErr := mr.NextPart()
		if partErr == io.EOF {
			break
		}
		if partErr != nil {
			h.writeMultipartErr(w, partErr)
			return
		}
		fmt.Println()
		switch part.FormName() {
		case "file":
			if h.handleFilePart(ctx, w, part) {
				return
			}

		case "content":
			if content != "" {
				_ = part.Close()
				continue
			}

			value, ok := h.handleContentPart(w, part)
			if !ok {
				return
			}
			content = value

		default:
			_ = part.Close()
		}
	}

	if content == "" {
		utils.Error(w, http.StatusBadRequest, "missing file field")
		return
	}

	h.scanAndWrite(ctx, w, strings.NewReader(content))
}

func (h *ApiHandler) handleFilePart(ctx context.Context, w http.ResponseWriter, part *multipart.Part) bool {
	filename := part.FileName()
	logged := utils.NewLoggingReader(io.NopCloser(part), h.logger, filename, time.Second)

	h.scanAndWrite(ctx, w, logged)

	_ = logged.Close()
	return true
}

func (h *ApiHandler) handleContentPart(w http.ResponseWriter, part *multipart.Part) (string, bool) {
	b, readErr := io.ReadAll(io.LimitReader(part, 1<<20))
	_ = part.Close()
	if readErr != nil {
		utils.Error(w, http.StatusBadRequest, "invalid multipart form")
		return "", false
	}
	return string(b), true
}

func (h *ApiHandler) scanAndWrite(ctx context.Context, w http.ResponseWriter, r io.Reader) {
	message, scanErr := h.apiStore.Scan(ctx, r)
	if scanErr != nil {
		switch message {
		case nil:
			utils.Error(w, http.StatusInternalServerError, scanErr)
			return
		default:
			utils.Error(w, http.StatusInternalServerError, message.Error)
			return
		}
	}
	utils.Success(w, http.StatusOK, message)
}

func (h *ApiHandler) writeMultipartErr(w http.ResponseWriter, err error) {
	var maxErr *http.MaxBytesError
	if errors.As(err, &maxErr) {
		maxFileSize := utils.FormatBytes(h.apiStore.MaxUploadBytes())
		utils.Error(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("file size too large, max file size is: %s", maxFileSize))
		return
	}
	utils.Error(w, http.StatusBadRequest, "invalid multipart form")
}
