package model

import "time"

type Contact struct {
	ID          int    `json:"id,omitempty" bson:"_id,omitempty"`
	ContactName string `json:"contact_name,omitempty" bson:"contact_name,omitempty"`
	Age         int    `json:"age,omitempty" bson:"age,omitempty"`
	Mobile      string `json:"mobile,omitempty" bson:"mobile,omitempty"`
}

type Question struct {
	ID           int       `json:"id" bson:"id,omitempty"`
	Question     string    `json:"question" bson:"question"`
	Options      []string  `json:"options" bson:"options"`
	CorrectAns   string    `json:"correct_answer" bson:"correct_answer"`
	Reason       string    `json:"reason" bson:"reason"`
	Hidden       bool      `json:"hidden" bson:"hidden"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	LastModified time.Time `json:"last_modified" bson:"last_modified"`
}

type Response struct {
	Message    string      `json:"message"`
	StatusCode int         `json:"status"`
	Data       interface{} `json:"data,omitempty"`
}
