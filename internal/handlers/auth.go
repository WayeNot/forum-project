package handlers

import (
	"database/sql"
	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var DB *sql.DB

func Login(w http.ResponseWriter, r *http.Request) {
	var (
		usernameOrMail string
		password       string
	)

	if r.Method == "POST" {
		usernameOrMail = r.FormValue("usernameOrMail")
		password = r.FormValue("password")

		if len(usernameOrMail) == 0 || len(password) == 0 {
			println("Champ(s) manquant(s) !")
			templates.Render("auth/login", w, r)
			return
		}

		var passwordDB string
		var userID int

		query := `SELECT password, id FROM users WHERE username = ? OR mail = ?`
		err := db.DB.QueryRow(query, usernameOrMail, usernameOrMail).Scan(&passwordDB, &userID)

		if len(passwordDB) == 0 {
			println("Erreur d'authentification !")
			templates.Render("auth/login", w, r)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(password))

		if err != nil {
			println(err.Error())
			println("Erreur d'authentification !")
			templates.Render("auth/login", w, r)
			return
		}

		idSessionId := uuid.New().String()

		request := `UPDATE sessions SET is_active = FALSE WHERE user_id = ?`
		_, err = db.DB.Exec(request, userID)

		request = `INSERT INTO sessions (user_id, session_id) VALUES (?, ?)`
		_, err = db.DB.Exec(request, userID, idSessionId)

		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    idSessionId,
			MaxAge:   86400, // Au bout de 24 heures, le cookie se delete automatiquement !
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	templates.Render("auth/login", w, r)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var (
		username string
		mail     string
		password string
		banner   string
		pp_url   string
		bio      string
		status   string
	)

	if r.Method == "POST" {
		username = r.FormValue("username")
		mail = r.FormValue("mail")
		password = r.FormValue("password")

		if len(username) == 0 || len(mail) == 0 || len(password) == 0 {
			println("Champ(s) manquant(s) !")
			templates.Render("auth/register", w, r)
			return
		}

		banner = ""
		pp_url = "https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExZTNud3o0NzV1eHZkOGl4ZmhmcDJycWNndTNmODcxdDZoMWY3ZTd3aCZlcD12MV9naWZzX3NlYXJjaCZjdD1n/GeG3Ulpo8WrwpNMpUz/giphy.gif"
		bio = "No bio yet."
		status = "Online"

		var mailDb any

		//

		reqIsMailExist := `SELECT id FROM users WHERE mail = ? OR username = ? LIMIT 1`
		err := db.DB.QueryRow(reqIsMailExist, username, mail).Scan(&mailDb)

		if mailDb != nil {
			println("Erreur d'authentification !")
			templates.Render("auth/register", w, r)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			println(err.Error())
			println("Erreur d'authentification !")
			templates.Render("auth/register", w, r)
			return
		}

		request := `INSERT INTO users (username, mail, password, banner, pp_url, bio, status) VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err = db.DB.Exec(request, username, mail, hashedPassword, banner, pp_url, bio, status)

		if err != nil {
			println(err.Error())
			templates.Render("auth/register", w, r)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.Render("auth/register", w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		MaxAge:   0,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}