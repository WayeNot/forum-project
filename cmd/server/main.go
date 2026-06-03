package main

import (
	"fmt"
	"net/http"

	"github.com/WayeNot/forum-project/internal/db"
	"github.com/WayeNot/forum-project/internal/handlers"
)

const port = ":5500"

func main() {
	db.Init("forum.db")

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/createPost", handlers.CreatePost)
	http.HandleFunc("/createTag", handlers.CreateTag)
	
	fmt.Printf("✅ Serveur lancé sur http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}