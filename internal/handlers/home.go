package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
)

type UserData struct {
	ID                 int
	Username           string
	Mail               string
	Banner             string
	PpURL              string
	Bio                string
	FavoriteInstrument string
	PreferredGenres    string
	ProfileTheme       string
	CustomStatus       string
}

func getAllTags() []string {
	const queryTags = `SELECT name FROM tags`
	rows, err := db.DB.Query(queryTags)
	if err != nil {
		println(err.Error())
		return []string{}
	}
	defer rows.Close()

	var allTags []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			println(err.Error())
			continue
		}
		allTags = append(allTags, name)
	}
	return allTags
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

		if err == nil {
			const requestUser = `SELECT id, username, mail, banner, pp_url, bio, favorite_instrument, preferred_genres, profile_theme, custom_status FROM users WHERE id = ?`
			err = db.DB.QueryRow(requestUser, user_id).Scan(&userData.ID, &userData.Username, &userData.Mail, &userData.Banner, &userData.PpURL, &userData.Bio, &userData.FavoriteInstrument, &userData.PreferredGenres, &userData.ProfileTheme, &userData.CustomStatus)
			if err == nil {
				isLogged = true
			}
		}
	}

	filter := r.URL.Query().Get("filter")
	if (filter == "created" || filter == "liked") && !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tagFilter := strings.TrimSpace(r.URL.Query().Get("tag"))

	var postsQuery string
	if filter == "popular" {
		postsQuery = `SELECT posts.id, posts.title, posts.description, posts.author_id, posts.image_url, posts.tags, posts.created_at, users.username FROM posts INNER JOIN users ON posts.author_id = users.id ORDER BY (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id AND post_likes.vote = 1) DESC, posts.created_at DESC`
	} else {
		postsQuery = `SELECT posts.id, posts.title, posts.description, posts.author_id, posts.image_url, posts.tags, posts.created_at, users.username FROM posts INNER JOIN users ON posts.author_id = users.id ORDER BY posts.created_at DESC`
	}
	rows, err := db.DB.Query(postsQuery)
	if err != nil {
		println(err.Error())
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	rowsData := []map[string]any{}

	for rows.Next() {
		var id int
		var title, description, imageURL, authorName string
		var authorID int
		var tagsStr string
		var tags []string
		var createdAt string

		err = rows.Scan(&id, &title, &description, &authorID, &imageURL, &tagsStr, &createdAt, &authorName)
		if err != nil {
			println(err.Error())
			continue
		}

		if tagsStr != "" {
			rawTags := strings.Split(tagsStr, ",")
			tags = make([]string, 0, len(rawTags))
			for _, t := range rawTags {
				trimmed := strings.TrimSpace(t)
				if trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		} else {
			tags = []string{}
		}

		if tagFilter != "" {
			found := false
			for _, t := range tags {
				if strings.EqualFold(t, tagFilter) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if filter == "created" && authorID != userData.ID {
			continue
		}

		var likesCount int
		const queryLikes = `SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND vote = 1`
		_ = db.DB.QueryRow(queryLikes, id).Scan(&likesCount)

		var dislikesCount int
		const queryDislikes = `SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND vote = -1`
		_ = db.DB.QueryRow(queryDislikes, id).Scan(&dislikesCount)

		userVote := 0
		if isLogged {
			const queryUserVote = `SELECT vote FROM post_likes WHERE post_id = ? AND user_id = ?`
			_ = db.DB.QueryRow(queryUserVote, id, userData.ID).Scan(&userVote)
		}

		if filter == "liked" && userVote != 1 {
			continue
		}

		getUserPpURL := func(userID int) string {
			var ppURL string
			const queryPpURL = `SELECT pp_url FROM users WHERE id = ?`
			err := db.DB.QueryRow(queryPpURL, userID).Scan(&ppURL)
			if err != nil {
				return "https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExZTNud3o0NzV1eHZkOGl4ZmhmcDJycWNndTNmODcxdDZoMWY3ZTd3aCZlcD12MV9naWZzX3NlYXJjaCZjdD1n/GeG3Ulpo8WrwpNMpUz/giphy.gif"
			}
			return ppURL
		}

		updateTimeAgo := func(createdAt string) string {
			const queryTimeAgo = `SELECT strftime('%s', 'now') - strftime('%s', ?) AS time_ago`
			var timeAgo int
			err := db.DB.QueryRow(queryTimeAgo, createdAt).Scan(&timeAgo)
			if err != nil {
				println(err.Error())
			}
			if timeAgo < 60 {
				return "quelques secondes"
			} else if timeAgo < 3600 {
				minutes := timeAgo / 60
				return strconv.Itoa(minutes) + " minute(s)"
			} else if timeAgo < 86400 {
				hours := timeAgo / 3600
				return strconv.Itoa(hours) + " heure(s)"
			} else {
				days := timeAgo / 86400
				return strconv.Itoa(days) + " jours"
			}
		}

		var commentCount int
		const queryCommentCount = `SELECT COUNT(*) FROM comments WHERE post_id = ?`
		_ = db.DB.QueryRow(queryCommentCount, id).Scan(&commentCount)

		postData := map[string]any{
			"id":             id,
			"title":          title,
			"description":    description,
			"author_id":      authorID,
			"image_url":      imageURL,
			"tags":           tags,
			"created_at":     createdAt,
			"time_ago":       updateTimeAgo(createdAt),
			"author_name":    authorName,
			"author_pp":      getUserPpURL(authorID),
			"likes_count":    likesCount,
			"dislikes_count": dislikesCount,
			"user_vote":      userVote,
			"comment_count":  commentCount,
		}

		rowsData = append(rowsData, postData)
	}

	data := map[string]any{
		"IsLogged":   isLogged,
		"UserData":   userData,
		"Posts":      rowsData,
		"Tags":       getAllTags(),
		"CurrentTag": tagFilter,
		"Filter":     filter,
	}

	templates.Render("home", w, data)
}
