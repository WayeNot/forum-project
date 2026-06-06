package templates

import (
	"html/template"
	"net/http"
)

func Render(path string, w http.ResponseWriter, data any) {
	tmpl, err := template.ParseFiles("web/templates/" + path + ".html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tmpl.Execute(w, data)
}

func ErrorPage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	Render("error", w, map[string]any{
		"StatusCode": statusCode,
		"Message":    message,
	})
}