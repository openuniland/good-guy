package frameworks

type HtmlName struct {
	Home          string
	PrivacyPolicy string
}

var TemplateEndpoints = "./pkg/frameworks"

var VIEWS = &HtmlName{
	Home:          TemplateEndpoints + "/web/index.html",
	PrivacyPolicy: TemplateEndpoints + "/web/privacy-policy.html",
}
