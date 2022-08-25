package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"ticket_app_api/show"
	"ticket_app_api/config"
)

func main() {
	fmt.Println("Application running at Port 12345")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	config.Client, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	err := config.Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer config.Client.Disconnect(ctx)
	router := mux.NewRouter()
	router.HandleFunc("/show", show.CreateShowEndpoint).Methods("POST")
	router.HandleFunc("/getshows", show.GetShowsEndpoint).Methods("GET")
	router.HandleFunc("/getshow{id}", show.GetShowByIdEndpoint).Methods("GET")
	router.HandleFunc("/updateshow{id}", show.UpdateShowEndpoint).Methods("POST")
	router.HandleFunc("/deleteshow{id}", show.DeleteShowEndpoint).Methods("DELETE")
	http.ListenAndServe(":12345", router)
}