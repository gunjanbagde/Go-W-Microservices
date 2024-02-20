package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	_ = WriteJson(w, http.StatusAccepted, payload)

}

func HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err := ReadJson(w, r, &requestPayload)
	if err != nil {
		errorJson(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		authenticate(w, requestPayload.Auth)
	default:
		errorJson(w, errors.New("unknown error"))
	}

}

func authenticate(w http.ResponseWriter, authpayload AuthPayload) {

	jsonData, _ := json.MarshalIndent(authpayload, "", "\t")

	request, err := http.NewRequest("POST", "http://localhost:8081/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		errorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		errorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		errorJson(w, errors.New("invalid creds"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		errorJson(w, errors.New("error calling authentication service"))
		return
	}

	var jsonFromSvc jsonResponse

	err = json.NewDecoder(request.Body).Decode(&jsonFromSvc)
	if err != nil {
		errorJson(w, err)
		return
	}

	if jsonFromSvc.Error {
		errorJson(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromSvc.Data,
	}

	WriteJson(w, http.StatusAccepted, payload)
}
