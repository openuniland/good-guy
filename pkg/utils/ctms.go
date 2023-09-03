package utils

import "github.com/openuniland/good-guy/external/types"

const (
	STUDY           = "Học"
	STUDY_ONLINE    = "Học trực tuyến"
	REST            = "Nghỉ"
	EXAM            = "Thi"
	EXTRACURRICULAR = "Ngoại khoá"
)

func ExamScheduleMessage(message string, examSchedule types.ExamSchedule) string {
	return message + ":\n" +
		`-----------------------` + "\n" +
		`STT: ` + examSchedule.SerialNumber + "\n" +
		`Thời gian: ` + examSchedule.Time + "\n" +
		`Phòng thi: ` + examSchedule.ClassRoom + "\n" +
		`Môn thi: ` + examSchedule.SubjectName + "\n" +
		`Mã DS thi: ` + examSchedule.ExamListCode + "\n"
}

func IsExamScheduleExisted(oldExamSchedule, newExamSchedule types.ExamSchedule) bool {
	return oldExamSchedule.SubjectName == newExamSchedule.SubjectName
}

func IsExamScheduleRoomChanged(oldExamSchedule, newExamSchedule types.ExamSchedule) bool {
	return oldExamSchedule.SubjectName == newExamSchedule.SubjectName &&
		oldExamSchedule.ClassRoom != newExamSchedule.ClassRoom
}

func IsExamScheduleTimeChanged(oldExamSchedule, newExamSchedule types.ExamSchedule) bool {
	return oldExamSchedule.SubjectName == newExamSchedule.SubjectName &&
		oldExamSchedule.Time != newExamSchedule.Time
}

func DailyScheduleMessage(message string, dailySchedule *types.DailySchedule) string {
	return message + "\n" +
		`-----------------------` + "\n" +
		`Giờ: ` + dailySchedule.Time + "\n" +
		`Phòng: ` + dailySchedule.ClassRoom + "\n" +
		`Môn học: ` + dailySchedule.SubjectName + "\n" +
		`Giảng viên: ` + dailySchedule.Lecturer + "\n" +
		`Lớp: ` + dailySchedule.ClassCode + "\n"
}

func GetClassStatus(status string, session string) string {
	switch status {
	case STUDY:
		return "📝 Bạn có môn học vào " + session + " nha:"
	case STUDY_ONLINE:
		return "📝 Bạn có môn học trực tuyến vào " + session + " nha:"
	case REST:
		return "🆘🆘🆘 Môn học " + session + " nay của bạn đã bị hủy (hoặc nghỉ học) nha:"
	case EXAM:
		return "💯 Bạn có môn thi vào " + session + " nay nha:"
	case EXTRACURRICULAR:
		return "🫦 Bạn có môn học ngoại khóa vào " + session + " nay nha:"
	default:
		return "😱 Bạn có môn học với trạng thái không xác định vào " + session + " nay nha:"
	}
}
