package templates

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func Render(path string, w http.ResponseWriter, data any) {
	name := filepath.Base(path + ".html")
	tmpl, err := template.New(name).Funcs(template.FuncMap{
		"contains": func(str, substr string) bool {
			return strings.Contains(str, substr)
		},
		"split": func(str, sep string) []string {
			if str == "" {
				return nil
			}
			return strings.Split(str, sep)
		},
		"hasTag": func(tags []string, tag string) bool {
			for _, t := range tags {
				if t == tag {
					return true
				}
			}
			return false
		},
	}).ParseFiles("web/templates/" + path + ".html")
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