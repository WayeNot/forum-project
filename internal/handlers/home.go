package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
)

type UserData struct {
	Username string
	Mail     string
	Banner   string
	PpURL    string
	Bio      string
}

func getAllTags() []string {
	const queryTags = `SELECT name FROM tags`
	var allTagsStr string
	err := db.DB.QueryRow(queryTags).Scan(&allTagsStr)
	if err != nil {
		println(err.Error())
	}
	if allTagsStr == "" {
		return []string{}
	}
	allTags := strings.Split(allTagsStr, ",")
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

		if err != nil {
			println(err.Error())
		} else {
			const requestUser = `SELECT username, mail, banner, pp_url, bio FROM users WHERE id = ?`
			err = db.DB.QueryRow(requestUser, user_id).Scan(&userData.Username, &userData.Mail, &userData.Banner, &userData.PpURL, &userData.Bio)

			if err != nil {
				println(err.Error())
			}
			isLogged = true
		}
	}

	const postsQuery = `SELECT posts.id, posts.title, posts.description, posts.author_id, posts.image_url, posts.tags, posts.created_at, users.username FROM posts INNER JOIN users ON posts.author_id = users.id ORDER BY posts.created_at DESC`
	rows, err := db.DB.Query(postsQuery)

	if err != nil {
		println(err.Error())
	}

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
		}

		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
		} else {
			tags = []string{}
		}

		getUserPpURL := func(userID int) string {
			var ppURL string
			const queryPpURL = `SELECT pp_url FROM users WHERE id = ?`
			err := db.DB.QueryRow(queryPpURL, userID).Scan(&ppURL)
			if err != nil {
				println(err.Error())
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
				return strconv.Itoa(days) + "  jours"
			}
		}

		postData := map[string]any{
			"id":          id,
			"title":       title,
			"description": description,
			"author_id":   authorID,
			"image_url":   imageURL,
			"tags":        tags,
			"created_at":  createdAt,
			"time_ago":    updateTimeAgo(createdAt),
			"author_name": authorName,
			"author_pp":   getUserPpURL(authorID),
		}

		rowsData = append(rowsData, postData)
	}

	data := map[string]any{
		"IsLogged": isLogged,
		"UserData": userData,
		"Posts":    rowsData,
		"Tags":     getAllTags(),
	}

	templates.Render("home", w, data)
}
