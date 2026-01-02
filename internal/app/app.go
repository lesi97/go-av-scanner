package app

import (
	"log"

	"github.com/lesi97/go-av-scanner/internal/api"
	"github.com/lesi97/go-av-scanner/internal/scanner/clamscan"
	"github.com/lesi97/go-av-scanner/internal/store"
	"github.com/lesi97/go-av-scanner/internal/utils"
)

type Application struct {
	Logger					*log.Logger
	ApiHandler 				*api.ApiHandler
}


func NewApplication() (*Application, error) {
	logger := utils.NewColourLogger("brightMagenta")
	const maxUploadBytes int64 = 10 << 30  // 10.73741824gb or 10,737,418,240 bytes
	sc, err := clamscan.New("", maxUploadBytes)
	if err != nil {
		logger.Fatalf("failed to initialise clamdscan: %v", err)
	}

	apiStore := store.NewApiStore(logger, sc, maxUploadBytes)
	apiHandler := api.NewApiHandler(logger, apiStore)

	app := &Application{
		Logger: logger,
		ApiHandler: apiHandler,
	}

	return app, nil
}