package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickhildpac/bookings-goapp/internal/config"
	"github.com/nickhildpac/bookings-goapp/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	// mux := pat.New()
	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	_ = app
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(WriteToConsole)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/general-quarters", handlers.Repo.General)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.Availabilty)
	mux.Post("/search-availability", handlers.Repo.PostAvailabilty)
	mux.Post("/search-availability-json", handlers.Repo.AvailabiltyJSON)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/contact", handlers.Repo.Contact)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
