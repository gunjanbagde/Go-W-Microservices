package main

import (
	"auth/cmd/api/data"
	"bytes"
	"encoding/json"
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

	//log Authentication
	err = logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		errorJson(w,err)
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	WriteJson(w, http.StatusAccepted, payload)
}

func logRequest(name string, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service/writelog", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}
