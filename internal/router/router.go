package router

import (
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/go-av-scanner/internal/app"
	"github.com/lesi97/go-av-scanner/internal/middleware"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	const apiKey = "1324"

	routes := chi.NewRouter()
	// routes.Use(middleware.ApiKey(apiKey))
	
	routes.Get("/api/health", http.HandlerFunc(middleware.Run(app.Logger, app.ApiHandler.HandleHealth)))
	routes.Post("/api/scan", http.HandlerFunc(middleware.Run(app.Logger, app.ApiHandler.HandleScan)))

	uiEnabled, err := strconv.ParseBool(os.Getenv("ENABLE_UI"))
	if err == nil && uiEnabled {
		routes.Handle("/*", http.FileServer(http.Dir("/app/ui/dist")))
	}

	return routes
}