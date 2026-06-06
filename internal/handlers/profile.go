package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
)

func UserProfile(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		templates.ErrorPage(w, http.StatusNotFound, "Page non trouvée.")
		return
	}
	usernamePrefixed := parts[2]
	if !strings.HasPrefix(usernamePrefixed, "@") {
		templates.ErrorPage(w, http.StatusNotFound, "Page non trouvée.")
		return
	}
	username := strings.TrimPrefix(usernamePrefixed, "@")

	var targetUser UserData
	const queryUser = `SELECT id, username, mail, banner, pp_url, bio FROM users WHERE username = ?`
	err := db.DB.QueryRow(queryUser, username).Scan(&targetUser.ID, &targetUser.Username, &targetUser.Mail, &targetUser.Banner, &targetUser.PpURL, &targetUser.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			templates.ErrorPage(w, http.StatusNotFound, "Profil introuvable.")
		} else {
			templates.ErrorPage(w, http.StatusInternalServerError, "Erreur lors de la récupération du profil.")
		}
		return
	}

	loggedUser, isLogged := getLoggedUser(r)

	const postsQuery = `SELECT posts.id, posts.title, posts.description, posts.author_id, posts.image_url, posts.tags, posts.created_at, users.username FROM posts INNER JOIN users ON posts.author_id = users.id WHERE posts.author_id = ? ORDER BY posts.created_at DESC`
	rows, err := db.DB.Query(postsQuery, targetUser.ID)
	var posts []map[string]any
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var title, description, imageURL, authorName, createdAt string
			var authorID int
			var tagsStr string
			var tags []string

			err = rows.Scan(&id, &title, &description, &authorID, &imageURL, &tagsStr, &createdAt, &authorName)
			if err != nil {
				continue
			}

			if tagsStr != "" {
				rawTags := strings.Split(tagsStr, ",")
				for _, t := range rawTags {
					trimmed := strings.TrimSpace(t)
					if trimmed != "" {
						tags = append(tags, trimmed)
					}
				}
			}

			var likesCount int
			_ = db.DB.QueryRow(`SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND vote = 1`, id).Scan(&likesCount)

			var dislikesCount int
			_ = db.DB.QueryRow(`SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND vote = -1`, id).Scan(&dislikesCount)

			userVote := 0
			if isLogged {
				_ = db.DB.QueryRow(`SELECT vote FROM post_likes WHERE post_id = ? AND user_id = ?`, id, loggedUser.ID).Scan(&userVote)
			}

			var commentCount int
			_ = db.DB.QueryRow(`SELECT COUNT(*) FROM comments WHERE post_id = ?`, id).Scan(&commentCount)

			posts = append(posts, map[string]any{
				"id":             id,
				"title":          title,
				"description":    description,
				"author_id":      authorID,
				"image_url":      imageURL,
				"tags":           tags,
				"created_at":     createdAt,
				"author_name":    authorName,
				"author_pp":      targetUser.PpURL,
				"likes_count":    likesCount,
				"dislikes_count": dislikesCount,
				"user_vote":      userVote,
				"comment_count":  commentCount,
			})
		}
	}

	data := map[string]any{
		"IsLogged":   isLogged,
		"UserData":   loggedUser,
		"TargetUser": targetUser,
		"Posts":      posts,
		"IsSelf":     isLogged && loggedUser.ID == targetUser.ID,
	}

	templates.Render("profile", w, data)
}

func UserSettings(w http.ResponseWriter, r *http.Request) {
	loggedUser, isLogged := getLoggedUser(r)
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	csrfToken := GetOrCreateCSRFToken(w, r)

	if r.Method == "POST" {
		if !VerifyCSRFToken(r) {
			templates.ErrorPage(w, http.StatusForbidden, "Session expirée ou requête CSRF invalide.")
			return
		}

		bio := strings.TrimSpace(r.FormValue("bio"))
		banner := strings.TrimSpace(r.FormValue("banner"))
		ppURL := strings.TrimSpace(r.FormValue("pp_url"))

		if ppURL == "" {
			ppURL = "https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExZTNud3o0NzV1eHZkOGl4ZmhmcDJycWNndTNmODcxdDZoMWY3ZTd3aCZlcD12MV9naWZzX3NlYXJjaCZjdD1n/GeG3Ulpo8WrwpNMpUz/giphy.gif"
		}

		const updateUser = `UPDATE users SET bio = ?, banner = ?, pp_url = ? WHERE id = ?`
		_, err := db.DB.Exec(updateUser, bio, banner, ppURL, loggedUser.ID)
		if err != nil {
			templates.ErrorPage(w, http.StatusInternalServerError, "Impossible de sauvegarder vos paramètres.")
			return
		}

		http.Redirect(w, r, "/user/@"+loggedUser.Username, http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"IsLogged":  isLogged,
		"UserData":  loggedUser,
		"CSRFToken": csrfToken,
	}

	templates.Render("auth/settings", w, data)
}
