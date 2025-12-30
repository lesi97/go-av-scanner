package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/go-av-scanner/internal/app"
	"github.com/lesi97/go-av-scanner/internal/router"
	"github.com/lesi97/go-av-scanner/internal/utils"
)



func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	routes := router.SetupRoutes(app)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		IdleTimeout: time.Minute,
		Handler: routes,
		ReadTimeout: 0,
		WriteTimeout: 0,
	}

	utils.Startup(fmt.Sprintf(":%d", port))

	err = server.ListenAndServe() 
	if err != nil {
		app.Logger.Fatal(err)
	}
}
