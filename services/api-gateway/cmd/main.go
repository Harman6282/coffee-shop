package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	internalMiddleware "github.com/Harman6282/coffee-shop/services/api-gateway/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

const port = ":8080"

func main() {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.RequestID)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
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
		r.Use(internalMiddleware.RequireAuth)
		r.Handle("/*", http.StripPrefix("/orders", orderProxy))
	})

	log.Println("api gateway running on port: ", port)
	log.Fatal(http.ListenAndServe(port, r))
}
