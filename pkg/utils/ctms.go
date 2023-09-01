package utils

import "github.com/openuniland/good-guy/external/types"

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
