package handlers

import (
	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
	"net/http"
)

type TagData struct {
	Id          int
	Name        string
	Description string
}

func CreateTag(w http.ResponseWriter, r *http.Request) {
	var postData TagData

	userData, _ := getLoggedUser(r)
	csrfToken := GetOrCreateCSRFToken(w, r)

	if r.Method == "POST" {
		if !VerifyCSRFToken(r) {
			http.Error(w, "Invalid CSRF Token", http.StatusForbidden)
			return
		}

		postData.Name = r.FormValue("name")
		postData.Description = r.FormValue("description")

		if postData.Name == "" || postData.Description == "" {
			templates.Render("creator/createTag", w, map[string]any{"Error": "Le nom et la description sont requis", "CSRFToken": csrfToken, "IsLogged": true, "UserData": userData})
			return
		}

		const insertTag = `INSERT INTO tags (name, description) VALUES (?, ?)`
		_, err := db.DB.Exec(insertTag, postData.Name, postData.Description)
		if err != nil {
			templates.Render("creator/createTag", w, map[string]any{"Error": "Erreur lors de la création du tag", "CSRFToken": csrfToken, "IsLogged": true, "UserData": userData})
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates.Render("creator/createTag", w, map[string]any{"CSRFToken": csrfToken, "IsLogged": true, "UserData": userData})
}
