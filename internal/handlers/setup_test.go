package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/nickhildpac/bookings-goapp/internal/config"
	"github.com/nickhildpac/bookings-goapp/internal/models"
	"github.com/nickhildpac/bookings-goapp/internal/render"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

var app config.AppConfig
var session *scs.SessionManager

func getRoutes() http.Handler {
	gob.Register(models.Reservations{})
	app.InProd = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProd

	app.Session = session
	tc, err := CreateTestTemplateCacheNew()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	log.Println("cache colod")
	app.TemplateCache = tc
	app.UseCache = false
	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplate(&app)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(WriteToConsole)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/general-quarters", Repo.General)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/search-availability", Repo.Availabilty)
	mux.Post("/search-availability", Repo.PostAvailabilty)
	mux.Post("/search-availability-json", Repo.AvailabiltyJSON)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	mux.Get("/contact", Repo.Contact)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux

}
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hit on the page")
		next.ServeHTTP(w, r)
	})
}

// Nosurf adds CSRF protection to all host requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProd,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// Session load loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCacheNew() (map[string]*template.Template, error) {
	// myCache := make(map[string] *template.Template)
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./../../templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./../../templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./../../templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
