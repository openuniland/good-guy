package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username            string             `json:"username" validate:"required"`
	Password            string             `json:"password" validate:"required"`
	SubscribedID        string             `json:"subscribed_id"`
	SubjectHTML         string             `json:"subject_html"`
	IsSubscribedSubject bool               `json:"is_subscribed_subject"`
	IsTrackTimetable    bool               `json:"is_track_timetable"`
	IsExamDay           bool               `json:"is_exam_day"`
	IsDeleted           bool               `json:"is_deleted" default:"false"`
	CreatedAt           string             `json:"created_at"`
	UpdatedAt           string             `json:"updated_at"`
}
