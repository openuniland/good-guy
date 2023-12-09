package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Provider string

type User struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username            string             `bson:"username" json:"username" validate:"required"`
	Password            string             `bson:"password" json:"password" validate:"required"`
	SubscribedID        string             `bson:"subscribed_id" json:"subscribed_id" validate:"required"`
	SubjectHTML         string             `bson:"subject_html" json:"subject_html"`
	IsSubscribedSubject bool               `bson:"is_subscribed_subject" json:"is_subscribed_subject"`
	IsTrackTimetable    bool               `bson:"is_track_timetable" json:"is_track_timetable"`
	IsExamDay           bool               `bson:"is_exam_day" json:"is_exam_day"`
	IsDeleted           bool               `bson:"is_deleted" json:"is_deleted" default:"false"`
	IsDisabled          bool               `bson:"is_disabled" json:"is_disabled" default:"false"`
	CreatedAt           primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt           primitive.DateTime `bson:"updated_at" json:"updated_at"`
	LoginProvider       Provider           `bson:"login_provider" json:"login_provider"`
	SessionId           string             `bson:"session_id" json:"session_id"`
	AspxAuth            string             `bson:"aspx_auth" json:"aspx_auth"`
	IS_VIP              bool               `bson:"is_vip" json:"is_vip"`
}
