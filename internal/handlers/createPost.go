package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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

	userData, _ := getLoggedUser(r)
	postData.Author_id = userData.ID

	csrfToken := GetOrCreateCSRFToken(w, r)

	if r.Method == "POST" {
		if !VerifyCSRFToken(r) {
			http.Error(w, "Invalid CSRF Token", http.StatusForbidden)
			return
		}

		err := r.ParseForm()
		if err != nil {
			templates.Render("creator/createPost", w, map[string]any{"Error": "Formulaire invalide", "CSRFToken": csrfToken, "Tags": getAllTags(), "IsLogged": true, "UserData": userData})
			return
		}

		action := strings.TrimSpace(r.FormValue("action"))
		if action == "createTag" {
			newTagName := strings.TrimSpace(r.FormValue("newTagName"))
			if newTagName == "" {
				writeJSONError(w, http.StatusBadRequest, "Le nom du tag est requis")
				return
			}

			if tagExists(newTagName) {
				writeJSONError(w, http.StatusBadRequest, "Ce tag existe déjà")
				return
			}

			err = insertTag(newTagName, "")
			if err != nil {
				writeJSONError(w, http.StatusInternalServerError, "Erreur lors de la création du tag")
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"tag": newTagName})
			return
		}

		postData.Title = r.FormValue("title")
		postData.Description = r.FormValue("description")
		postData.Image_url = strings.TrimSpace(r.FormValue("media"))
		if postData.Image_url != "" && !strings.HasPrefix(postData.Image_url, "http://") && !strings.HasPrefix(postData.Image_url, "https://") {
			postData.Image_url = ""
		}

		selectedTags := r.Form["tags"]
		newTagName := strings.TrimSpace(r.FormValue("newTagName"))

		if newTagName != "" {
			if tagExists(newTagName) {
				templates.Render("creator/createPost", w, map[string]any{"Error": "Ce tag existe déjà. Choisissez-le dans la liste ou utilisez un autre nom.", "CSRFToken": csrfToken, "Tags": getAllTags(), "IsLogged": true, "UserData": userData})
				return
			}

			err = insertTag(newTagName, "")
			if err != nil {
				templates.Render("creator/createPost", w, map[string]any{"Error": "Erreur lors de la création du tag", "CSRFToken": csrfToken, "Tags": getAllTags(), "IsLogged": true, "UserData": userData})
				return
			}
			if !containsTag(selectedTags, newTagName) {
				selectedTags = append(selectedTags, newTagName)
			}
		}

		postData.Tags = strings.Join(selectedTags, ",")

		if postData.Title == "" || postData.Description == "" {
			templates.Render("creator/createPost", w, map[string]any{"Error": "Le titre et la description sont requis", "CSRFToken": csrfToken, "Tags": getAllTags(), "IsLogged": true, "UserData": userData})
			return
		}

		const insertPost = `INSERT INTO posts (title, description, author_id, image_url, tags) VALUES (?, ?, ?, ?, ?)`
		_, err = db.DB.Exec(insertPost, postData.Title, postData.Description, postData.Author_id, postData.Image_url, postData.Tags)
		if err != nil {
			templates.Render("creator/createPost", w, map[string]any{"Error": "Erreur lors de la création du post", "CSRFToken": csrfToken, "Tags": getAllTags(), "IsLogged": true, "UserData": userData})
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates.Render("creator/createPost", w, map[string]any{"CSRFToken": csrfToken, "Tags": getAllTags(), "IsLogged": true, "UserData": userData})
}

func tagExists(name string) bool {
	if strings.TrimSpace(name) == "" {
		return false
	}

	const query = `SELECT COUNT(*) FROM tags WHERE LOWER(name) = LOWER(?)`
	var count int
	err := db.DB.QueryRow(query, name).Scan(&count)
	return err == nil && count > 0
}

func insertTag(name, description string) error {
	const query = `INSERT INTO tags (name, description) VALUES (?, ?)`
	_, err := db.DB.Exec(query, name, description)
	return err
}

func containsTag(tags []string, name string) bool {
	name = strings.TrimSpace(strings.ToLower(name))
	for _, tag := range tags {
		if strings.TrimSpace(strings.ToLower(tag)) == name {
			return true
		}
	}
	return false
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{"error": message})
}
