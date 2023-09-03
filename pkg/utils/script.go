package utils

func HelpScript() string {
	return `*Các lệnh hỗ trợ:*
/login: Đăng nhập CTMS.
/remove: Đăng xuất CTMS.
/help: Hiển thị hướng dẫn.
/examday: Đăng ký nhận thông báo ngày thi.
/un_examday: Hủy đăng ký nhận thông báo ngày thi.
/force_logout: Đăng xuất tất cả phiên hoạt động ở hệ thống`
}

type CommandType struct {
	Login     string
	Remove    string
	Help      string
	Examday   string
	UnExamday string
}

var Command = CommandType{
	Login:     "/login",
	Remove:    "/remove",
	Help:      "/help",
	Examday:   "/examday",
	UnExamday: "/un_examday",
}
