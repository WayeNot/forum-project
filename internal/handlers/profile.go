package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
	"golang.org/x/crypto/bcrypt"
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
	const queryUser = `SELECT id, username, mail, banner, pp_url, bio, favorite_instrument, preferred_genres, profile_theme, custom_status FROM users WHERE username = ?`
	err := db.DB.QueryRow(queryUser, username).Scan(&targetUser.ID, &targetUser.Username, &targetUser.Mail, &targetUser.Banner, &targetUser.PpURL, &targetUser.Bio, &targetUser.FavoriteInstrument, &targetUser.PreferredGenres, &targetUser.ProfileTheme, &targetUser.CustomStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			templates.ErrorPage(w, http.StatusNotFound, "Profil introuvable.")
		} else {
			templates.ErrorPage(w, http.StatusInternalServerError, "Erreur lors de la récupération du profil.")
		}
		return
	}

	if targetUser.PpURL == "" || strings.Contains(targetUser.PpURL, "giphy.gif") {
		targetUser.PpURL = "/static/images/default-avatar.svg"
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

	csrfToken := GetOrCreateCSRFToken(w, r)

	if r.Method == "POST" {
		if !VerifyCSRFToken(r) {
			templates.ErrorPage(w, http.StatusForbidden, "Session expirée ou requête CSRF invalide.")
			return
		}

		err := r.ParseForm()
		if err != nil {
			templates.Render("auth/settings", w, map[string]any{"Error": "Données de formulaire invalides.", "IsLogged": isLogged, "UserData": loggedUser, "CSRFToken": csrfToken})
			return
		}

		action := r.FormValue("action")
		if action == "delete" {
			_, err1 := db.DB.Exec("DELETE FROM sessions WHERE user_id = ?", loggedUser.ID)
			_, err2 := db.DB.Exec("DELETE FROM users WHERE id = ?", loggedUser.ID)
			if err1 != nil || err2 != nil {
				templates.ErrorPage(w, http.StatusInternalServerError, "Erreur lors de la suppression du compte.")
				return
			}
			cookie := &http.Cookie{
				Name:     "session_id",
				Value:    "",
				MaxAge:   -1,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		username := strings.TrimSpace(r.FormValue("username"))
		mail := strings.TrimSpace(r.FormValue("mail"))
		bio := strings.TrimSpace(r.FormValue("bio"))
		banner := strings.TrimSpace(r.FormValue("banner"))
		ppURL := strings.TrimSpace(r.FormValue("pp_url"))
		customStatus := strings.TrimSpace(r.FormValue("custom_status"))
		favoriteInstrument := strings.TrimSpace(r.FormValue("favorite_instrument"))
		profileTheme := strings.TrimSpace(r.FormValue("profile_theme"))

		preferredGenresList := r.Form["preferred_genres"]
		preferredGenres := strings.Join(preferredGenresList, ",")

		if username == "" || mail == "" {
			templates.Render("auth/settings", w, map[string]any{"Error": "Le nom d'utilisateur et l'email sont obligatoires.", "IsLogged": isLogged, "UserData": loggedUser, "CSRFToken": csrfToken})
			return
		}

		if username != loggedUser.Username {
			var exists int
			db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND id != ?", username, loggedUser.ID).Scan(&exists)
			if exists > 0 {
				templates.Render("auth/settings", w, map[string]any{"Error": "Ce nom d'utilisateur est déjà utilisé.", "IsLogged": isLogged, "UserData": loggedUser, "CSRFToken": csrfToken})
				return
			}
		}

		if mail != loggedUser.Mail {
			var exists int
			db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE mail = ? AND id != ?", mail, loggedUser.ID).Scan(&exists)
			if exists > 0 {
				templates.Render("auth/settings", w, map[string]any{"Error": "Cette adresse email est déjà utilisée.", "IsLogged": isLogged, "UserData": loggedUser, "CSRFToken": csrfToken})
				return
			}
		}

		currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")
		if currentPassword != "" || newPassword != "" {
			if currentPassword == "" || newPassword == "" {
				templates.Render("auth/settings", w, map[string]any{"Error": "Veuillez saisir le mot de passe actuel et le nouveau mot de passe pour le modifier.", "IsLogged": isLogged, "UserData": loggedUser, "CSRFToken": csrfToken})
				return
			}

			var passwordDB string
			err = db.DB.QueryRow("SELECT password FROM users WHERE id = ?", loggedUser.ID).Scan(&passwordDB)
			if err != nil {
				templates.Render("auth/settings", w, map[string]any{"Error": "Erreur de base de données.", "IsLogged": isLogged, "UserData": loggedUser, "CSRFToken": csrfToken})
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(currentPassword))
			if err != nil {
				templates.Render("auth/settings", w, map[string]any{"Error": "Le mot de passe actuel est incorrect.", "IsLogged": isLogged, "UserData": loggedUser, "CSRFToken": csrfToken})
				return
			}

			hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			if err != nil {
				templates.Render("auth/settings", w, map[string]any{"Error": "Erreur lors du chiffrement du mot de passe.", "IsLogged": isLogged, "UserData": loggedUser, "CSRFToken": csrfToken})
				return
			}

			_, err = db.DB.Exec("UPDATE users SET password = ? WHERE id = ?", string(hashed), loggedUser.ID)
			if err != nil {
				templates.Render("auth/settings", w, map[string]any{"Error": "Impossible de mettre à jour le mot de passe.", "IsLogged": isLogged, "UserData": loggedUser, "CSRFToken": csrfToken})
				return
			}
		}

		if ppURL == "" || strings.Contains(ppURL, "giphy.gif") {
			ppURL = "/static/images/default-avatar.svg"
		}

		const updateUser = `UPDATE users SET username = ?, mail = ?, bio = ?, banner = ?, pp_url = ?, favorite_instrument = ?, preferred_genres = ?, profile_theme = ?, custom_status = ? WHERE id = ?`
		_, err = db.DB.Exec(updateUser, username, mail, bio, banner, ppURL, favoriteInstrument, preferredGenres, profileTheme, customStatus, loggedUser.ID)
		if err != nil {
			templates.ErrorPage(w, http.StatusInternalServerError, "Impossible de sauvegarder vos paramètres.")
			return
		}

		http.Redirect(w, r, "/user/@"+username, http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"IsLogged":  isLogged,
		"UserData":  loggedUser,
		"CSRFToken": csrfToken,
	}

	templates.Render("auth/settings", w, data)
}
