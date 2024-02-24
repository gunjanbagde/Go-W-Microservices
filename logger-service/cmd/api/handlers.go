package main

import (
	"log"
	"logger/data"
	"net/http"

)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func WriteLog(m data.Models, w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload

	ReadJson(w, r, &requestPayload)

	// insert the data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := m.LogEntry.Insert(event)
	if err != nil {
		log.Println(err)
		errorJson(w, err, http.StatusBadRequest)
		return
	}

	// create the response we'll send back as JSON
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	// write the response back as JSON
	WriteJson(w, http.StatusAccepted, resp)
}
