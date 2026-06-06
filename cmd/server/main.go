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

	http.HandleFunc("/post/", handlers.PostDetail)
	http.HandleFunc("/post/like", handlers.LikePost)
	http.HandleFunc("/post/dislike", handlers.DislikePost)
	http.HandleFunc("/post/comment", handlers.CommentPost)
	http.HandleFunc("/comment/like", handlers.LikeComment)
	http.HandleFunc("/comment/dislike", handlers.DislikeComment)
	http.HandleFunc("/post/edit", handlers.EditPost)
	http.HandleFunc("/post/delete", handlers.DeletePost)

	http.HandleFunc("/user/", handlers.UserProfile)
	http.HandleFunc("/settings", handlers.UserSettings)

	fmt.Printf("✅ Serveur lancé sur http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}