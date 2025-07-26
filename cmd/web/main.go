package main

import (
	"fmt"
	"net/http"

	"github.com/nickhildpac/learn-go-web/pkg/handlers"
)

const portNumber = ":8081"

func main() {
	fmt.Println("Hello go web")
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)
	http.ListenAndServe(portNumber, nil)
}
