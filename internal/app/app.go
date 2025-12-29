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
	sc := clamscan.New("")

	apiStore := store.NewApiStore(logger, sc)
	apiHandler := api.NewApiHandler(logger, apiStore)

	app := &Application{
		Logger: logger,
		ApiHandler: apiHandler,
	}

	return app, nil
}