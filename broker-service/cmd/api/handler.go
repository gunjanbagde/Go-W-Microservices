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
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
	case "log":
		logItem(w, requestPayload.Log)
	default:
		errorJson(w, errors.New("unknown error"))
	}

}

func logItem(w http.ResponseWriter, entry LogPayload) {

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service/writelog", bytes.NewBuffer(jsonData))
	if err != nil {
		errorJson(w, err)
		return
	}

	request.Header.Set("Content-Type", "applicatin/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		errorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		errorJson(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	WriteJson(w,http.StatusAccepted,payload)

}

func authenticate(w http.ResponseWriter, authpayload AuthPayload) {

	jsonData, _ := json.MarshalIndent(authpayload, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
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

	err = json.NewDecoder(response.Body).Decode(&jsonFromSvc)
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
