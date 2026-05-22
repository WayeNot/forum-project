package main

import (
	"log"
	"net/http"
	"os"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/templates/home.html")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		log.Println("Connexion !")
		http.ServeFile(w, r, "./web/templates/auth/login.html")
		return
	}

	if r.Method == http.MethodPost {

		username := r.FormValue("username")
		password := r.FormValue("password")

		log.Println(username, password)

		return
	}

	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}

func main() {
	fs := http.FileServer(http.Dir("./web/static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Println("✅ Serveur lancé sur le port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
