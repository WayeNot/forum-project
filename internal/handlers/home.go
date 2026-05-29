package handlers

import (
	"net/http"
	"strings"

	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
)

type UserData struct {
	Username string
	Mail     string
	Banner   string
	PpURL    string
	Bio      string
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	isLogged := false
	var userData UserData

	session, err := r.Cookie("session_id")

	if err == nil && session.Value != "" {
		var user_id int
		const requestUserId = `SELECT user_id FROM sessions WHERE session_id = ? LIMIT 1`
		cleanSessionValue := strings.TrimSpace(session.Value)
		err = db.DB.QueryRow(requestUserId, cleanSessionValue).Scan(&user_id)

		if err != nil {
			println(err.Error())
		} else {
			const requestUser = `SELECT username, mail, banner, pp_url, bio FROM users WHERE id = ?`
			err = db.DB.QueryRow(requestUser, user_id).Scan(&userData.Username, &userData.Mail, &userData.Banner, &userData.PpURL, &userData.Bio)

			if err != nil {
				println(err.Error())
			}
			isLogged = true
		}
	}

	data := map[string]any{
		"IsLogged": isLogged,
		"UserData": userData,
	}

	templates.Render("home", w, data)
}
