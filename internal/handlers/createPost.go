package handlers

import (
	"net/http"

	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
)

type PostData struct {
	Id          int
	Title       string
	Description string
	Tags        string
	Author_id   int
	Image_url   string
	CreatedAt   string
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var postData PostData
	var userID int

	session, err := r.Cookie("session_id")
	if err != nil || session.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	const requestUserID = `SELECT user_id FROM sessions WHERE session_id = ? AND is_active = TRUE LIMIT 1`
	err = db.DB.QueryRow(requestUserID, session.Value).Scan(&userID)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	postData.Author_id = userID

	if r.Method == "POST" {
		postData.Title = r.FormValue("title")
		postData.Description = r.FormValue("description")
		postData.Image_url = r.FormValue("media")
		postData.Tags = r.FormValue("tags")

		if postData.Title == "" || postData.Description == "" {
			println("Le titre et la description sont requis")
			templates.Render("creator/createPost", w, r)
			return
		}

		const insertPost = `INSERT INTO posts (title, description, author_id, image_url, tags) VALUES (?, ?, ?, ?, ?)`
		_, err := db.DB.Exec(insertPost, postData.Title, postData.Description, postData.Author_id, postData.Image_url, postData.Tags)
		if err != nil {
			println(err.Error())
			println("Erreur lors de la creation du post")
			templates.Render("creator/createPost", w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates.Render("creator/createPost", w, r)
}
