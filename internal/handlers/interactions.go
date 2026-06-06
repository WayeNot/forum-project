package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
)

func getLoggedUser(r *http.Request) (UserData, bool) {
	var userData UserData
	session, err := r.Cookie("session_id")
	if err != nil || session.Value == "" {
		return userData, false
	}

	var user_id int
	const requestUserId = `SELECT user_id FROM sessions WHERE session_id = ? AND is_active = TRUE LIMIT 1`
	err = db.DB.QueryRow(requestUserId, session.Value).Scan(&user_id)
	if err != nil {
		return userData, false
	}

	const requestUser = `SELECT id, username, mail, banner, pp_url, bio, favorite_instrument, preferred_genres, profile_theme, custom_status FROM users WHERE id = ?`
	err = db.DB.QueryRow(requestUser, user_id).Scan(&userData.ID, &userData.Username, &userData.Mail, &userData.Banner, &userData.PpURL, &userData.Bio, &userData.FavoriteInstrument, &userData.PreferredGenres, &userData.ProfileTheme, &userData.CustomStatus)
	if err != nil {
		return userData, false
	}

	return userData, true
}

func PostDetail(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		templates.ErrorPage(w, http.StatusNotFound, "Post introuvable.")
		return
	}
	idStr := parts[2]
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		templates.ErrorPage(w, http.StatusNotFound, "Identifiant du post invalide.")
		return
	}

	isLogged := false
	userData, logged := getLoggedUser(r)
	if logged {
		isLogged = true
	}

	var title, description, imageURL, authorName, createdAt string
	var authorID int
	var tagsStr string

	const queryPost = `SELECT posts.title, posts.description, posts.author_id, posts.image_url, posts.tags, posts.created_at, users.username FROM posts INNER JOIN users ON posts.author_id = users.id WHERE posts.id = ?`
	err = db.DB.QueryRow(queryPost, postID).Scan(&title, &description, &authorID, &imageURL, &tagsStr, &createdAt, &authorName)
	if err != nil {
		if err == sql.ErrNoRows {
			templates.ErrorPage(w, http.StatusNotFound, "Post introuvable.")
		} else {
			templates.ErrorPage(w, http.StatusInternalServerError, "Erreur lors de la récupération du post.")
		}
		return
	}

	var tags []string
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
	const queryLikes = `SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND vote = 1`
	_ = db.DB.QueryRow(queryLikes, postID).Scan(&likesCount)

	var dislikesCount int
	const queryDislikes = `SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND vote = -1`
	_ = db.DB.QueryRow(queryDislikes, postID).Scan(&dislikesCount)

	userVote := 0
	if isLogged {
		const queryUserVote = `SELECT vote FROM post_likes WHERE post_id = ? AND user_id = ?`
		_ = db.DB.QueryRow(queryUserVote, postID, userData.ID).Scan(&userVote)
	}

	var authorPp string
	const queryPp = `SELECT pp_url FROM users WHERE id = ?`
	_ = db.DB.QueryRow(queryPp, authorID).Scan(&authorPp)
	if authorPp == "" {
		authorPp = "https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExZTNud3o0NzV1eHZkOGl4ZmhmcDJycWNndTNmODcxdDZoMWY3ZTd3aCZlcD12MV9naWZzX3NlYXJjaCZjdD1n/GeG3Ulpo8WrwpNMpUz/giphy.gif"
	}

	const queryComments = `SELECT comments.id, comments.content, comments.created_at, comments.author_id, users.username, users.pp_url, comments.parent_id FROM comments INNER JOIN users ON comments.author_id = users.id WHERE comments.post_id = ? ORDER BY comments.created_at ASC`
	rows, err := db.DB.Query(queryComments, postID)
	var comments []map[string]any
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cID, cAuthorID int
			var cContent, cCreatedAt, cUsername, cPp string
			var cParentID sql.NullInt64
			if err := rows.Scan(&cID, &cContent, &cCreatedAt, &cAuthorID, &cUsername, &cPp, &cParentID); err == nil {
				var cLikesCount int
				_ = db.DB.QueryRow(`SELECT COUNT(*) FROM comment_likes WHERE comment_id = ? AND vote = 1`, cID).Scan(&cLikesCount)

				var cDislikesCount int
				_ = db.DB.QueryRow(`SELECT COUNT(*) FROM comment_likes WHERE comment_id = ? AND vote = -1`, cID).Scan(&cDislikesCount)

				cUserVote := 0
				if isLogged {
					_ = db.DB.QueryRow(`SELECT vote FROM comment_likes WHERE comment_id = ? AND user_id = ?`, cID, userData.ID).Scan(&cUserVote)
				}

				if cPp == "" {
					cPp = "https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExZTNud3o0NzV1eHZkOGl4ZmhmcDJycWNndTNmODcxdDZoMWY3ZTd3aCZlcD12MV9naWZzX3NlYXJjaCZjdD1n/GeG3Ulpo8WrwpNMpUz/giphy.gif"
				}

				parentIDVal := 0
				if cParentID.Valid {
					parentIDVal = int(cParentID.Int64)
				}

				comments = append(comments, map[string]any{
					"id":             cID,
					"content":        cContent,
					"created_at":     cCreatedAt,
					"author_id":      cAuthorID,
					"author_name":    cUsername,
					"author_pp":      cPp,
					"likes_count":    cLikesCount,
					"dislikes_count": cDislikesCount,
					"user_vote":      cUserVote,
					"parent_id":      parentIDVal,
				})
			}
		}
	}

	csrfToken := GetOrCreateCSRFToken(w, r)

	data := map[string]any{
		"IsLogged":       isLogged,
		"UserData":       userData,
		"PostID":         postID,
		"Title":          title,
		"Description":    description,
		"ImageURL":       imageURL,
		"AuthorID":       authorID,
		"AuthorName":     authorName,
		"AuthorPp":       authorPp,
		"CreatedAt":      createdAt,
		"Tags":           tags,
		"LikesCount":     likesCount,
		"DislikesCount":  dislikesCount,
		"UserVote":       userVote,
		"Comments":       comments,
		"IsAuthor":       isLogged && authorID == userData.ID,
		"CSRFToken":      csrfToken,
	}

	templates.Render("postDetail", w, data)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	votePost(w, r, 1)
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	votePost(w, r, -1)
}

func votePost(w http.ResponseWriter, r *http.Request, voteType int) {
	userData, logged := getLoggedUser(r)
	if !logged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var currentVote int
	const selectVote = `SELECT vote FROM post_likes WHERE post_id = ? AND user_id = ?`
	err = db.DB.QueryRow(selectVote, postID, userData.ID).Scan(&currentVote)
	if err != nil {
		if err == sql.ErrNoRows {
			const insertVote = `INSERT INTO post_likes (post_id, user_id, vote) VALUES (?, ?, ?)`
			_, _ = db.DB.Exec(insertVote, postID, userData.ID, voteType)
		}
	} else {
		if currentVote == voteType {
			const deleteVote = `DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`
			_, _ = db.DB.Exec(deleteVote, postID, userData.ID)
		} else {
			const updateVote = `UPDATE post_likes SET vote = ? WHERE post_id = ? AND user_id = ?`
			_, _ = db.DB.Exec(updateVote, voteType, postID, userData.ID)
		}
	}

	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}

func CommentPost(w http.ResponseWriter, r *http.Request) {
	userData, logged := getLoggedUser(r)
	if !logged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		if !VerifyCSRFToken(r) {
			templates.ErrorPage(w, http.StatusForbidden, "CSRF invalide.")
			return
		}

		postIDStr := r.FormValue("post_id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
			return
		}

		parentIDStr := r.FormValue("parent_id")
		var parentID sql.NullInt64
		if parentIDStr != "" && parentIDStr != "0" {
			pID, err := strconv.Atoi(parentIDStr)
			if err == nil {
				parentID.Int64 = int64(pID)
				parentID.Valid = true
			}
		}

		const insertComment = `INSERT INTO comments (post_id, author_id, content, parent_id) VALUES (?, ?, ?, ?)`
		_, err = db.DB.Exec(insertComment, postID, userData.ID, content, parentID)
		if err != nil {
			templates.ErrorPage(w, http.StatusInternalServerError, "Impossible de sauvegarder votre commentaire.")
			return
		}

		http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LikeComment(w http.ResponseWriter, r *http.Request) {
	voteComment(w, r, 1)
}

func DislikeComment(w http.ResponseWriter, r *http.Request) {
	voteComment(w, r, -1)
}

func voteComment(w http.ResponseWriter, r *http.Request, voteType int) {
	userData, logged := getLoggedUser(r)
	if !logged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentIDStr := r.URL.Query().Get("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var currentVote int
	const selectVote = `SELECT vote FROM comment_likes WHERE comment_id = ? AND user_id = ?`
	err = db.DB.QueryRow(selectVote, commentID, userData.ID).Scan(&currentVote)
	if err != nil {
		if err == sql.ErrNoRows {
			const insertVote = `INSERT INTO comment_likes (comment_id, user_id, vote) VALUES (?, ?, ?)`
			_, _ = db.DB.Exec(insertVote, commentID, userData.ID, voteType)
		}
	} else {
		if currentVote == voteType {
			const deleteVote = `DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?`
			_, _ = db.DB.Exec(deleteVote, commentID, userData.ID)
		} else {
			const updateVote = `UPDATE comment_likes SET vote = ? WHERE comment_id = ? AND user_id = ?`
			_, _ = db.DB.Exec(updateVote, voteType, commentID, userData.ID)
		}
	}

	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	userData, logged := getLoggedUser(r)
	if !logged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		templates.ErrorPage(w, http.StatusNotFound, "Post introuvable.")
		return
	}

	var authorID int
	var title, description, imageURL, tagsStr string
	const queryPost = `SELECT title, description, author_id, image_url, tags FROM posts WHERE id = ?`
	err = db.DB.QueryRow(queryPost, postID).Scan(&title, &description, &authorID, &imageURL, &tagsStr)
	if err != nil {
		templates.ErrorPage(w, http.StatusNotFound, "Post introuvable.")
		return
	}

	if authorID != userData.ID {
		templates.ErrorPage(w, http.StatusForbidden, "Vous n'avez pas l'autorisation de modifier ce post.")
		return
	}

	csrfToken := GetOrCreateCSRFToken(w, r)

	if r.Method == "POST" {
		if !VerifyCSRFToken(r) {
			templates.ErrorPage(w, http.StatusForbidden, "CSRF invalide.")
			return
		}

		err = r.ParseForm()
		if err != nil {
			templates.ErrorPage(w, http.StatusBadRequest, "Formulaire invalide.")
			return
		}

		title = r.FormValue("title")
		description = r.FormValue("description")
		imageURL = strings.TrimSpace(r.FormValue("media"))
		if imageURL != "" && !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
			imageURL = ""
		}
		
		selectedTags := r.Form["tags"]
		tagsStr = strings.Join(selectedTags, ",")

		if title == "" || description == "" {
			data := map[string]any{
				"Error":       "Titre et description requis",
				"PostID":      postID,
				"Title":       title,
				"Description": description,
				"ImageURL":    imageURL,
				"Tags":        tagsStr,
				"AllTags":     getAllTags(),
				"CSRFToken":   csrfToken,
			}
			templates.Render("creator/editPost", w, data)
			return
		}

		const updatePost = `UPDATE posts SET title = ?, description = ?, image_url = ?, tags = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = db.DB.Exec(updatePost, title, description, imageURL, tagsStr, postID)
		if err != nil {
			templates.ErrorPage(w, http.StatusInternalServerError, "Erreur serveur.")
			return
		}

		http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
		return
	}

	data := map[string]any{
		"PostID":      postID,
		"Title":       title,
		"Description": description,
		"ImageURL":    imageURL,
		"Tags":        tagsStr,
		"AllTags":     getAllTags(),
		"CSRFToken":   csrfToken,
	}

	templates.Render("creator/editPost", w, data)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	userData, logged := getLoggedUser(r)
	if !logged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		templates.ErrorPage(w, http.StatusNotFound, "Post introuvable.")
		return
	}

	var authorID int
	const queryAuthor = `SELECT author_id FROM posts WHERE id = ?`
	err = db.DB.QueryRow(queryAuthor, postID).Scan(&authorID)
	if err != nil {
		templates.ErrorPage(w, http.StatusNotFound, "Post introuvable.")
		return
	}

	if authorID != userData.ID {
		templates.ErrorPage(w, http.StatusForbidden, "Vous n'avez pas l'autorisation de supprimer ce post.")
		return
	}

	const deletePost = `DELETE FROM posts WHERE id = ?`
	_, _ = db.DB.Exec(deletePost, postID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
