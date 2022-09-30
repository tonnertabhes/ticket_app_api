package show

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strconv"
	
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	
	"ticket_app_api/uuidGen"
	"ticket_app_api/config"
)

type Show struct {
	ID            string          	  `json:"id" bson:"id"`
	Name  		  string              `json:"name" bson:"name"`
	Description   string			  `json:"description" bson:"description"`
	Date  		  string              `json:"date" bson:"date"`
	Price 		  string		      `json:"price" bson:"price"`
	TicketHolders []TicketHolder 	  `json:"ticketholders" bson:"ticketholders"`
	Capacity	  string			  `json:"capacity" bson:"capacity"`
}

type TicketHolder struct {
	ID 		  string `json:"id" bson:"id"`
	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email     string `json:"email,omitempty" bson:"email,omitempty"`
	Phone     string `json:"phone,omitempty" bson:"phone,omitempty"`
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
	err := json.NewDecoder(request.Body).Decode(&show)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
	}
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", show}}
	result, _ := collection.UpdateOne(ctx, filter, update)
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


func AddTicketHolderEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	collection := config.Client.Database("theatre").Collection("showlist")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	params := mux.Vars(request)
	id := params["id"]
	amt, _ := strconv.Atoi(params["amt"])
	var ticketholder TicketHolder
	var show Show
	find := collection.FindOne(ctx, bson.M{"id": id}).Decode(&show)
	cap, _ := strconv.Atoi(show.Capacity)
	if find != nil {
		fmt.Println(find)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + find.Error() + `" }`))
		return
	}
	if (len(show.TicketHolders) + amt) > cap {
		response.Write([]byte(`{ "message": "Show is at capacity" }`))
		return
	}
	ticketholder.ID = uuidGen.GenerateUUID().String()
	err := json.NewDecoder(request.Body).Decode(&ticketholder)
	if err != nil {
		log.Fatal(err)
		return
	}
	filter := bson.D{{"id", id}}
 
	update := bson.D{
    	{"$push", bson.D{
        {"ticketholders", ticketholder},
    	}},
	}
 
	i := 0
	for i < amt {		
		_, e := collection.UpdateOne(ctx, filter, update)
		if e != nil {
			log.Fatal(e)
			return
		}
		i++
	}
}

func GetTicketHolderByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	collection := config.Client.Database("theatre").Collection("showlist")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	params := mux.Vars(request)
	id := params["id"]
	var show Show
	var ticketholder TicketHolder
	err := collection.FindOne(ctx, bson.M{"ticketholders.id": id}).Decode(&show)
	if err != nil {
		fmt.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	for i := range show.TicketHolders {
		if show.TicketHolders[i].ID == id {
			ticketholder = show.TicketHolders[i]
		}
	}
	json.NewEncoder(response).Encode(ticketholder)
}

func GetTicketHolders(response http.ResponseWriter, request *http.Request) {
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
	json.NewEncoder(response).Encode(show.TicketHolders)
}

func UpdateTicketHolder(response http.ResponseWriter, request *http.Request) {
	
}

func DeleteTicketHolder(response http.ResponseWriter, request *http.Request) {
	
}