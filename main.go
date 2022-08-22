package main

import (
	"fmt"
	"time"
	"log"
	"encoding/json"
	"context"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gorilla/mux"
)

type TicketHolder struct {
	FirstName string `json:"firstname" bson:"firstname"`
	LastName  string `json:"lastname" bson:"lastname"`
	Email     string `json:"email" bson:"email"`
	Phone     string `json"phone" bson:"phone"`
}

type Show struct {
	Name  		  string              `json:"name" bson:"name"`
	Date  		  string              `json:"date" bson:"date"`
	Price 		  string		      `json:"price" bson:"price"`
	TicketHolders []TicketHolder 	  `json:"ticketholder" bson:"ticketholder"`
}

var client *mongo.Client

func createShowEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type",  "application/json")
	var show Show
	err := json.NewDecoder(request.Body).Decode(&show)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("theatre").Collection("showlist")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, show)
	json.NewEncoder(response).Encode(result)
}

func getShowsEndpoint(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("content-type", "application/json")
		var Shows []Show
		collection := client.Database("theatre").Collection("showlist")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var show Show
			cursor.Decode(&show)
			Shows = append(Shows, show)
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(Shows)
}

func main() {
	fmt.Println("Application running at Port 12345")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	router := mux.NewRouter()
	router.HandleFunc("/show", createShowEndpoint).Methods("POST")
	router.HandleFunc("/shows", getShowsEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}