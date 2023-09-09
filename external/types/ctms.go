package types

type LoginCtmsRequest struct {
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	SubscribedID string `json:"subscribed_id" validate:"required"`
}

type LoginCtmsResponse struct {
	Cookie   string `json:"cookie"`
	Username string `json:"username"`
}

type LogoutCtmsRequest struct {
	Cookie string `json:"cookie" validate:"required"`
}

type GetDailyScheduleRequest struct {
	Cookie string `json:"cookie" validate:"required"`
}

type DailySchedule struct {
	SerialNumber string `json:"serial_number"`
	Time         string `json:"time"`
	ClassRoom    string `json:"class_room"`
	SubjectName  string `json:"subject_name"`
	Lecturer     string `json:"lecturer"`
	ClassCode    string `json:"class_code"`
	Status       string `json:"status"`
}

type ExamSchedule struct {
	SerialNumber string `json:"serial_number"`
	Time         string `json:"time"`
	ClassRoom    string `json:"class_room"`
	SubjectName  string `json:"subject_name"`
	ExamListCode string `json:"exam_list_code"`
}

type GetExamScheduleRequest struct {
	Cookie string `json:"cookie" validate:"required"`
}

type GetUpcomingExamScheduleResponse struct {
	CurrentExamsSchedules []ExamSchedule `json:"current_exams_schedule"`
	OldExamsSchedules     []ExamSchedule `json:"old_exams_schedule"`
}
