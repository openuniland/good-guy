package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cookie struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	Cookies   []string           `bson:"cookies" json:"cookies"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
