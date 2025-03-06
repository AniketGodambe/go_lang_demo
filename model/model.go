package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Contact struct {
	ID          int    `json:"id,omitempty" bson:"_id,omitempty"`
	ContactName string `json:"contact_name,omitempty" bson:"contact_name,omitempty"`
	Age         int    `json:"age,omitempty" bson:"age,omitempty"`
	Mobile      string `json:"mobile,omitempty" bson:"mobile,omitempty"`
}

type UpdateContactRequest struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	ContactName string             `json:"contact_name,omitempty"`
	Age         int                `json:"age,omitempty"`
}
