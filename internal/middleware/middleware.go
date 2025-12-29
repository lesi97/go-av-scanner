package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/lesi97/go-av-scanner/internal/utils"
)

func Run(logger *log.Logger, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path

		ApplyCORS(w, r)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		defer func() {
			pathColour := utils.Colours["brightBlack"]
			timeColour := utils.Colours["green"]
			duration := time.Since(start)
			if duration > 100*time.Millisecond {
				timeColour = utils.Colours["brightRed"] + utils.Colours["bold"]
			}
			logger.Printf("%s%s %stook %s%s", pathColour, path, timeColour, duration, utils.Colours["reset"])
		}()
		handlerFunc(w, r)
	}
}