package frameworks

import "os"

type HtmlName struct {
	Home          string
	PrivacyPolicy string
}

func Web() *HtmlName {
	pwd, _ := os.Getwd()
	VIEWS := &HtmlName{
		Home:          pwd + "/web/index.html",
		PrivacyPolicy: pwd + "/web/privacy-policy.html",
	}

	return VIEWS
}
