package model

import (
	"html/template"
	"net/http"
)

const homePageTemplatePath = "./web/template/index.html"

type HomePage struct {
	Domain string
	Error  string
}

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(homePageTemplatePath))
	homePage := HomePage{
		Domain: r.Host + r.URL.Path,
		Error:  r.URL.Query().Get("error"),
	}
	tmpl.Execute(w, homePage)
}
