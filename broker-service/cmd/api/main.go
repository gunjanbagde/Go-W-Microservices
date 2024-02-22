package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)

	r.HandleFunc("/", Broker).Methods("POST")
	r.HandleFunc("/handle", HandleSubmission).Methods("POST")

	fmt.Println("Server is running on :80")
	err := http.ListenAndServe(":80", cors(r))
	if err != nil {
		log.Panic(err)
	}
}
