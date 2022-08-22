package models

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