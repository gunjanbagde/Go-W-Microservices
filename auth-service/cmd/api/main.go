package main

import (
	"auth/cmd/api/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// type Config struct {
// 	Db   *sql.DB
// 	User data.Models
// }

func main() {

	log.Println("Starting authentication service")

	db := connectToDb()
	defer db.Close()

	models := data.New(db)

	r := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)

	r.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) { Authenticate(&models.User, w, r) }).Methods("POST")

	fmt.Println("Server is running on :80")
	err := http.ListenAndServe(":80", cors(r))
	if err != nil {
		log.Panic(err)
	}

}

func OpenDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDb() *sql.DB {
	dsn := os.Getenv("DSN")
	// dsn := "host=localhost port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"
	counts := 0

	for {
		connection, err := OpenDB(dsn)
		if err != nil {
			log.Println("Postgress not yet ready")
			counts++
		} else {
			log.Println("connected to Postgress")
			return connection
		}
		if counts > 10 {
			log.Println(err)
		}

		log.Println("Backing off for 2 seconds")
		time.Sleep(2 * time.Second)
	}
}
