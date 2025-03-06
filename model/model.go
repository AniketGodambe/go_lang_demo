package model

type Contact struct {
	ID          int    `json:"id,omitempty" bson:"_id,omitempty"`
	ContactName string `json:"contact_name,omitempty" bson:"contact_name,omitempty"`
	Age         int    `json:"age,omitempty" bson:"age,omitempty"`
	Mobile      string `json:"mobile,omitempty" bson:"mobile,omitempty"`
}

type Question struct {
	ID         int      `json:"id"`
	Question   string   `json:"question"`
	Options    []string `json:"options"`
	CorrectAns string   `json:"correct_answer"`
	Reason     string   `json:"reason"`
	Hidden     bool     `json:"hidden"`
}
type Response struct {
	Message    string      `json:"message"`
	StatusCode int         `json:"status"`
	Data       interface{} `json:"data,omitempty"`
}
