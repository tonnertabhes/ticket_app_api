package ticketholder

type TicketHolder struct {
	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email     string `json:"email,omitempty" bson:"email,omitempty"`
	Phone     string `json:"phone,omitempty" bson:"phone,omitempty"`
}
