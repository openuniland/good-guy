package types

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Cookie    string `json:"cookie"`
	IsSuccess bool   `json:"is_success"`
}

type LogoutRequest struct {
	Cookie string `json:"cookie" validate:"required"`
}

type GetDailyScheduleRequest struct {
	Cookie string `json:"cookie" validate:"required"`
}

type DailySchedule struct {
	SerialNumber string `json:"serial_number"`
	Time         string `json:"time"`
	ClassRoom    string `json:"class_room"`
	CourseName   string `json:"course_name"`
	Lecturer     string `json:"lecturer"`
	ClassCode    string `json:"class_code"`
	Status       string `json:"status"`
}
