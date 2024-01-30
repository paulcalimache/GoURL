package main

import (
	"net/http"

	"log"

	"github.com/paulcalimache/gourl/internal/db"
	"github.com/paulcalimache/gourl/internal/healthcheck"
	"github.com/paulcalimache/gourl/internal/middleware"
	"github.com/paulcalimache/gourl/internal/model"
	s "github.com/paulcalimache/gourl/internal/shortener"

	"github.com/gorilla/mux"
)

func main() {
	mongo := db.NewMongoDB()
	shortener := s.NewShortener(mongo)
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Healthcheck endpoints
	r.Methods(http.MethodGet).Path("/healthz").HandlerFunc(healthcheck.Healthz)
	r.Methods(http.MethodGet).Path("/readyz").HandlerFunc(healthcheck.Readyz)

	r.Methods(http.MethodPost).Path("/shorten").HandlerFunc(shortener.Shorten)
	r.Methods(http.MethodGet).Path("/").HandlerFunc(model.RenderHomePage)

	r.Methods(http.MethodGet).Path("/{key}").HandlerFunc(shortener.Redirect)

	r.Use(middleware.LoggingMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
