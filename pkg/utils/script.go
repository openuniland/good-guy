package utils

func HelpScript() string {
	return `*Các lệnh hỗ trợ:*
/login: Đăng nhập CTMS.
/logout: Đăng xuất CTMS.
/help: Hiển thị hướng dẫn.
/examday: Đăng ký nhận thông báo ngày thi.
/un_examday: Hủy đăng ký nhận thông báo ngày thi.`
}

type CommandType struct {
	Login     string
	Logout    string
	Help      string
	Examday   string
	UnExamday string
}

var command = CommandType{
	Login:     "/login",
	Logout:    "/logout",
	Help:      "/help",
	Examday:   "/examday",
	UnExamday: "/un_examday",
}

func ChatScript(id string, msg string) string {
	switch msg {
	case command.Login:
		//
		return ""
	case command.Logout:
		//
		return ""
	case command.Help:
		//
		return ""
	case command.Examday:
		//	return ""
		return ""
	case command.UnExamday:
		//	return ""
		return ""
	default:
		return "BOT"
	}

}
