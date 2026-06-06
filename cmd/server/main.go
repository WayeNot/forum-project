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

	http.HandleFunc("/", handlers.Chain(handlers.Home, handlers.SecurityHeaders))
	http.HandleFunc("/login", handlers.Chain(handlers.Login, handlers.SecurityHeaders, handlers.GuestOnly))
	http.HandleFunc("/register", handlers.Chain(handlers.Register, handlers.SecurityHeaders, handlers.GuestOnly))
	http.HandleFunc("/logout", handlers.Chain(handlers.Logout, handlers.SecurityHeaders, handlers.RequireAuth))
	http.HandleFunc("/createPost", handlers.Chain(handlers.CreatePost, handlers.SecurityHeaders, handlers.RequireAuth))
	http.HandleFunc("/createTag", handlers.Chain(handlers.CreateTag, handlers.SecurityHeaders, handlers.RequireAuth))

	http.HandleFunc("/post/", handlers.Chain(handlers.PostDetail, handlers.SecurityHeaders))
	http.HandleFunc("/post/like", handlers.Chain(handlers.LikePost, handlers.SecurityHeaders, handlers.RequireAuth))
	http.HandleFunc("/post/dislike", handlers.Chain(handlers.DislikePost, handlers.SecurityHeaders, handlers.RequireAuth))
	http.HandleFunc("/post/comment", handlers.Chain(handlers.CommentPost, handlers.SecurityHeaders, handlers.RequireAuth))
	http.HandleFunc("/comment/like", handlers.Chain(handlers.LikeComment, handlers.SecurityHeaders, handlers.RequireAuth))
	http.HandleFunc("/comment/dislike", handlers.Chain(handlers.DislikeComment, handlers.SecurityHeaders, handlers.RequireAuth))
	http.HandleFunc("/post/edit", handlers.Chain(handlers.EditPost, handlers.SecurityHeaders, handlers.RequireAuth))
	http.HandleFunc("/post/delete", handlers.Chain(handlers.DeletePost, handlers.SecurityHeaders, handlers.RequireAuth))

	http.HandleFunc("/user/", handlers.Chain(handlers.UserProfile, handlers.SecurityHeaders))
	http.HandleFunc("/settings", handlers.Chain(handlers.UserSettings, handlers.SecurityHeaders, handlers.RequireAuth))

	fmt.Printf("✅ Serveur lancé sur http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}