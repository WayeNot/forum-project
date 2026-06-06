package handlers

import (
	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
	"net/http"
	"strings"
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
	var user_id int

	session, err := r.Cookie("session_id")
	if err != nil || session.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	const requestUserId = `SELECT user_id FROM sessions WHERE session_id = ? AND is_active = TRUE LIMIT 1`
	err = db.DB.QueryRow(requestUserId, session.Value).Scan(&user_id)
	if err != nil {
		println("Erreur lors de la récupération de l'ID de l'utilisateur :", err.Error())
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	postData.Author_id = user_id

	csrfToken := GetOrCreateCSRFToken(w, r)

	if r.Method == "POST" {
		if !VerifyCSRFToken(r) {
			http.Error(w, "Invalid CSRF Token", http.StatusForbidden)
			return
		}

		err = r.ParseForm()
		if err != nil {
			templates.Render("creator/createPost", w, map[string]any{"Error": "Formulaire invalide", "CSRFToken": csrfToken, "Tags": getAllTags()})
			return
		}

		postData.Title = r.FormValue("title")
		postData.Description = r.FormValue("description")
		postData.Image_url = strings.TrimSpace(r.FormValue("media"))
		if postData.Image_url != "" && !strings.HasPrefix(postData.Image_url, "http://") && !strings.HasPrefix(postData.Image_url, "https://") {
			postData.Image_url = ""
		}
		
		selectedTags := r.Form["tags"]
		postData.Tags = strings.Join(selectedTags, ",")

		if postData.Title == "" || postData.Description == "" {
			templates.Render("creator/createPost", w, map[string]any{"Error": "Le titre et la description sont requis", "CSRFToken": csrfToken, "Tags": getAllTags()})
			return
		}

		const insertPost = `INSERT INTO posts (title, description, author_id, image_url, tags) VALUES (?, ?, ?, ?, ?)`
		_, err := db.DB.Exec(insertPost, postData.Title, postData.Description, postData.Author_id, postData.Image_url, postData.Tags)
		if err != nil {
			templates.Render("creator/createPost", w, map[string]any{"Error": "Erreur lors de la création du post", "CSRFToken": csrfToken, "Tags": getAllTags()})
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates.Render("creator/createPost", w, map[string]any{"CSRFToken": csrfToken, "Tags": getAllTags()})
}
