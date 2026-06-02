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

	if r.Method == "POST" {
		postData.Name = r.FormValue("name")
		postData.Description = r.FormValue("description")

		if postData.Name == "" || postData.Description == "" {
			println("Le nom et la description sont requis")
			templates.Render("creator/createTag", w, r)
			return
		}

		const insertTag = `INSERT INTO tags (name, description) VALUES (?, ?)`
		_, err := db.DB.Exec(insertTag, postData.Name, postData.Description)

		if err != nil {
			println(err.Error())
			println("Erreur lors de la création du tag")
			templates.Render("creator/createTag", w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	templates.Render("creator/createTag", w, r)
}
