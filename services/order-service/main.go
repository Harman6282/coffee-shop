package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const port = ":8082"

const dsn = "postgresql://admin:mypassword@localhost:5432/coffee_shop_db?sslmode=disable"

type Application struct {
	config Config
	db     *sql.DB
}

type Config struct {
	addr string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := NewPostgres(ctx, dsn)
	if err != nil {
		log.Panic("failed to connect Postgres", err)
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Order service is running..."))
	})

	cfg := Config{
		addr: port,
	}

	app := &Application{
		config: cfg,
		db: db,
	}

	log.Println("order service started at port:", port)
	err = http.ListenAndServe(app.config.addr, r)
	if err != nil {
		log.Fatal("failed to start server")
	}
}
