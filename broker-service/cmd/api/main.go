package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", Broker).Methods("POST")
	r.HandleFunc("/handle", HandleSubmission).Methods("POST")

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Panic(err)
	}
}
