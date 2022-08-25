package show

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	
	"ticket_app_api/uuidGen"
	"ticket_app_api/ticketholder"
	"ticket_app_api/config"
)

type Show struct {
	ID            string          	  			  `json:"id" bson:"id"`
	Name  		  string             			  `json:"name" bson:"name"`
	Date  		  string             			  `json:"date" bson:"date"`
	Price 		  string		      			  `json:"price" bson:"price"`
	TicketHolders []ticketholder.TicketHolder 	  `json:"ticketholders" bson:"ticketholders"`
}

func CreateShowEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type",  "application/json")
	var show Show
	generatedID := uuidGen.GenerateUUID()
	show.ID = generatedID.String()
	err := json.NewDecoder(request.Body).Decode(&show)
	if err != nil {
		log.Fatal(err)
	}
	collection := config.Client.Database("theatre").Collection("showlist")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, show)
	json.NewEncoder(response).Encode(result)
}

func GetShowByIdEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id := params["id"]
	var show Show
	collection := config.Client.Database("theatre").Collection("showlist")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&show)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(show)
}

func GetShowsEndpoint(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("content-type", "application/json")
		var Shows []Show
		collection := config.Client.Database("theatre").Collection("showlist")
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

func UpdateShowEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var show Show
	params := mux.Vars(request)
	id := params["id"]
	collection := config.Client.Database("theatre").Collection("showlist")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	jErr := json.NewDecoder(request.Body).Decode(&show)
	if jErr != nil {
		log.Fatal(jErr)
	}
	show.ID = id
	del, dErr := collection.DeleteOne(ctx, bson.M{"id": id})
	if dErr != nil {
		log.Fatal(dErr)
	}
	fmt.Println(del)
	result, _ := collection.InsertOne(ctx, show)
	json.NewEncoder(response).Encode(result)
}

func DeleteShowEndpoint(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id := params["id"]
	collection := config.Client.Database("theatre").Collection("showlist")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	del, err := collection.DeleteOne(ctx, bson.M{"id" : id})
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(del)
}