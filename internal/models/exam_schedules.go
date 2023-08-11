package models

import (
	"github.com/openuniland/good-guy/external/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExamSchedules struct {
	Id        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Username  string               `bson:"username" json:"username" validate:"required"`
	Subjects  []types.ExamSchedule `bson:"subjects" json:"subjects"`
	CreatedAt primitive.DateTime   `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime   `bson:"updated_at" json:"updated_at"`
}
