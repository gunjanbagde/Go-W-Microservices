package main

import (
	"auth/cmd/api/data"
	"errors"
	"fmt"
	"net/http"
)

func Authenticate(u *data.User,w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := ReadJson(w, r, &requestPayload)
	if err != nil {
		errorJson(w, err, http.StatusBadRequest)
		return
	}
	// models := data.Models{}

	user, err := u.GetByEmail(requestPayload.Email)
	if err != nil {
		errorJson(w, errors.New("invalid creds"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		errorJson(w, errors.New("invalid creds"), http.StatusBadRequest)
		return
	} 

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	
	WriteJson(w, http.StatusAccepted, payload)
}
