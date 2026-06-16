package handlers

import "os"

func secureCookie() bool {
	return os.Getenv("COOKIE_SECURE") == "true" || os.Getenv("APP_ENV") == "production"
}
