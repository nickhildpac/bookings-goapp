package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nickhildpac/learn-go-web/pkg/config"
	"github.com/nickhildpac/learn-go-web/pkg/handlers"
	"github.com/nickhildpac/learn-go-web/pkg/render"
)

const portNumber = ":8081"

func main() {
	var app config.AppConfig
	tc, err := render.CreateTemplateCacheNew()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)
	fmt.Println("Hello go web")
	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)
	http.ListenAndServe(portNumber, nil)
}
