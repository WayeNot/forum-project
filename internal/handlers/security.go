package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func GetOrCreateCSRFToken(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("csrf_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}
	token := uuid.New().String()
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return token
}

func VerifyCSRFToken(r *http.Request) bool {
	cookie, err := r.Cookie("csrf_token")
	if err != nil || cookie.Value == "" {
		return false
	}
	formToken := r.FormValue("csrf_token")
	return cookie.Value == formToken
}
