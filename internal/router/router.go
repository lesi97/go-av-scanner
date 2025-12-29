package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/go-av-scanner/internal/app"
	"github.com/lesi97/go-av-scanner/internal/middleware"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	const apiKey = "1324"

	routes := chi.NewRouter()
	// routes.Use(middleware.ApiKey(apiKey))
	
	routes.Get("/health", http.HandlerFunc(middleware.Run(app.Logger, app.ApiHandler.HandleHealth)))
	routes.Post("/scan", http.HandlerFunc(middleware.Run(app.Logger, app.ApiHandler.HandleScan)))


	return routes
}