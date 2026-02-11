package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpHandler "github.com/Peyton232/weather-forecast/pkg/http"
	"github.com/Peyton232/weather-forecast/pkg/nws"
	"github.com/Peyton232/weather-forecast/pkg/service"
)

func main() {
	handler := buildHandler()
	router := buildRouter(handler)

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func buildHandler() *httpHandler.Handler {
	nwsClient := nws.NewClient()
	forecastService := service.NewForecastService(nwsClient)
	return httpHandler.NewHandler(forecastService)
}

func buildRouter(handler *httpHandler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Get("/health", handler.Health)
	r.Get("/forecast", handler.GetForecast)
	r.Get("/", handler.Index)

	return r
}
