package handlers

import (
	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/templates"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func SetCookie(userID int) (*http.Cookie, error) {
	idSessionId := uuid.New().String()

	request := `UPDATE sessions SET is_active = FALSE WHERE user_id = ?`
	_, err := db.DB.Exec(request, userID)
	if err != nil {
		return nil, err
	}

	request = `INSERT INTO sessions (user_id, session_id, is_active) VALUES (?, ?, TRUE)`
	_, err = db.DB.Exec(request, userID, idSessionId)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    idSessionId,
		MaxAge:   86400,
		Path:     "/",
		HttpOnly: true,
		Secure:   secureCookie(),
		SameSite: http.SameSiteLaxMode,
	}

	return cookie, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		templates.Render("auth/login", w, nil)
		return
	}

	usernameOrMail := r.FormValue("usernameOrMail")
	password := r.FormValue("password")

	if len(usernameOrMail) == 0 || len(password) == 0 {
		templates.Render("auth/login", w, map[string]any{"Error": "Champs obligatoires manquants."})
		return
	}

	var passwordDB string
	var userID int

	query := `SELECT password, id FROM users WHERE username = ? OR mail = ?`
	err := db.DB.QueryRow(query, usernameOrMail, usernameOrMail).Scan(&passwordDB, &userID)
	if err != nil || len(passwordDB) == 0 {
		templates.Render("auth/login", w, map[string]any{"Error": "Identifiant ou mot de passe incorrect."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(password))
	if err != nil {
		templates.Render("auth/login", w, map[string]any{"Error": "Identifiant ou mot de passe incorrect."})
		return
	}

	cookie, err := SetCookie(userID)
	if err != nil {
		templates.Render("auth/login", w, map[string]any{"Error": "Impossible de creer la session."})
		return
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		templates.Render("auth/register", w, nil)
		return
	}

	username := r.FormValue("username")
	mail := r.FormValue("mail")
	password := r.FormValue("password")

	if len(username) == 0 || len(mail) == 0 || len(password) == 0 {
		templates.Render("auth/register", w, map[string]any{"Error": "Champs obligatoires manquants."})
		return
	}

	banner := ""
	pp_url := "/static/images/default-avatar.svg"
	bio := "Aucune bio pour le moment !"
	status := "En ligne"

	userID := 0
	reqIsMailExist := `SELECT id FROM users WHERE mail = ? OR username = ? LIMIT 1`
	err := db.DB.QueryRow(reqIsMailExist, mail, username).Scan(&userID)
	if err == nil {
		templates.Render("auth/register", w, map[string]any{"Error": "Le nom d'utilisateur ou l'adresse email est déjà utilisé."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		templates.Render("auth/register", w, map[string]any{"Error": "Erreur serveur de chiffrement."})
		return
	}

	request := `INSERT INTO users (username, mail, password, banner, pp_url, bio, status) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = db.DB.Exec(request, username, mail, hashedPassword, banner, pp_url, bio, status)
	if err != nil {
		templates.Render("auth/register", w, map[string]any{"Error": "Erreur lors de la création du compte."})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == nil && session.Value != "" {
		_, _ = db.DB.Exec(`UPDATE sessions SET is_active = FALSE WHERE session_id = ?`, session.Value)
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   secureCookie(),
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
