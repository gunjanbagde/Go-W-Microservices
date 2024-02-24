package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"time"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

// type Config struct {
// 	Models data.Models
// }

func main() {

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	//context to disconnect mongo
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//close
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// app := Config{
	// 	Models: data.New(client),
	// }

	models := data.New(client)

	r := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)

	r.HandleFunc("/writelog", func(w http.ResponseWriter, r *http.Request) { WriteLog(models, w, r) }).Methods("POST")

	fmt.Println("Server is running on: ",webPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s",webPort), cors(r))
	if err != nil {
		log.Panic(err)
	}

}


func connectToMongo() (*mongo.Client, error) {
	// create connect options
	clientOptions := options.Client().ApplyURI(mongoURL)

	//should not be hardcoded
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Connect to the MongoDB and return Client instance
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		return nil, err
	}

	
	return c, nil
}
