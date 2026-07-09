package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const port = ":8080"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Api gateway working")
	})

	userServiceURL, _ := url.Parse("http://localhost:8081")
	orderServiceURL, _ := url.Parse("http://localhost:8082")
	
	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	orderProxy := httputil.NewSingleHostReverseProxy(orderServiceURL)

	r.Route("/users", func(r chi.Router) {
		r.Handle("/*", http.StripPrefix("/users", userProxy))
	})
	
	r.Route("/orders", func(r chi.Router) {
		r.Handle("/*", http.StripPrefix("/orders", orderProxy))
	})

	log.Println("api gateway running on port: ", port)
	log.Fatal(http.ListenAndServe(port, r))
}
