package main

import (
	"log"
	"net/http"
)

func (app *Application) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user RegisterRequest
	err := readJson(r, &user)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		writeJson(w, "invalid credentials", nil, http.StatusBadRequest)
		return
	}

	hashedPass, err := HashPassword(user.Password)
	if err != nil {
		log.Println("failed to hash password", err)
	}
	user.Password = hashedPass

	err = app.User.Insert(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := writeJson(w, "user registerd", user.Email, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (app *Application) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user RegisterRequest
	err := readJson(r, &user)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		writeJson(w, "invalid credentials", nil, http.StatusBadRequest)
		return
	}

	hashedPassword, err := app.User.GetByEmail(ctx, user.Email)
	if err != nil {
		if err := writeJson(w, "invalid credentials", nil, http.StatusUnauthorized); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if !CheckPasswordHash(user.Password, hashedPassword) {
		if err := writeJson(w, "invalid password", nil, http.StatusUnauthorized); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := app.User.UpdateLastLogin(ctx, user.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString, expiresAt, err := GenerateJWTToken(user.Email)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	SetAuthCookie(w, tokenString, expiresAt)

	if err := writeJson(w, "login successful", user.Email, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) logout(w http.ResponseWriter, r *http.Request) {
	ClearAuthCookie(w)

	if err := writeJson(w, "logout successful", nil, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
