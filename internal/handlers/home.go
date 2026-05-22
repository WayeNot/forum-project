package handlers

import (
	"net/http"
	"github.com/WayeNot/forum-project/internal/templates"
)

func Home(w http.ResponseWriter, r *http.Request) {
	templates.Render("home", w, nil)
}