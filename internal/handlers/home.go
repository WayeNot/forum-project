package handlers

import (
	"github.com/WayeNot/forum-project/internal/templates"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")

	var is_logged bool

	if err != nil {
		if err == http.ErrNoCookie {
			is_logged = false
		} else {
			println(err.Error())
			is_logged = false
		}
	} else {
		is_logged = session.Value != ""
	}

	data := map[string]any{
        "IsLogged": is_logged,
    }

	templates.Render("home", w, data)
}