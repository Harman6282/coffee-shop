package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const dsn = "postgresql://admin:mypassword@localhost:5432/coffee_shop_db?sslmode=disable"
const port = ":8081"


type Application struct {
	config Config
	db     *sql.DB
	User   UserRepository
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

	cfg := Config{
		addr: port,
	}

	app := &Application{
		config: cfg,
		db:     db,
		User:   NewUserRepo(db),
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User service is running..."))
	})

	r.Post("/register", app.register)
	r.Post("/login", app.login)
	r.Post("/logout", app.logout)

	log.Println("User service started at port:", port)
	err = http.ListenAndServe(app.config.addr, r)
	if err != nil {
		log.Fatal("failed to start server")
	}
}
