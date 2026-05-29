package handlers

import (
	"net/http"
	"strings"

	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
)

type PostData struct {
	Title       string
	Description string
	Author      int
	image       string
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var user_id int
	var PostData PostData

	session, err := r.Cookie("session_id")

	if err == nil && session.Value != "" {
		const requestUserId = `SELECT user_id FROM sessions WHERE session_id = ? LIMIT 1`
		cleanSessionValue := strings.TrimSpace(session.Value)
		err = db.DB.QueryRow(requestUserId, cleanSessionValue).Scan(&user_id)

		if err != nil {
			println(err.Error())
		}
	}

	if r.Method == "POST" {
		PostData.Title = r.FormValue("title")
		PostData.Description = r.FormValue("title")
		PostData.Author = user_id
		PostData.image = r.FormValue("media")
	}

	templates.Render("home", w, r)

	// data := map[string]any{
	// 	"IsLogged": isLogged,
	// 	"UserData": userData,
	// }

	// templates.Render("home", w, data)
}
