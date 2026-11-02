package main

import (
	"log"
	"net/http"
	"os"
	"simple-blog-auth/auth"
	"simple-blog-auth/database"
	"simple-blog-auth/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to read .env: %v",err)
		
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT Secret is not set")
	}


	mux := http.NewServeMux() 

	
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	db := database.InitDB()
	mux.HandleFunc("POST /register", handlers.CreateAccount(db))
	store := auth.NewSessionStore()
	mux.HandleFunc("GET /{$}", auth.OptionalAuth(store,handlers.Home(db),jwtSecret))

	mux.HandleFunc("GET /articles/{id}", auth.OptionalAuth(store,handlers.ViewArticle(db),jwtSecret))
	mux.HandleFunc("GET /articles/{id}/edit",auth.RequireAuth(store,handlers.EditArticle(db),jwtSecret))
	mux.HandleFunc("POST /articles/{id}/edit",auth.RequireAuth(store,handlers.UpdateArticle(db),jwtSecret))
	mux.HandleFunc("POST /articles/{id}/delete",auth.RequireAuth(store,handlers.DeleteArticle(db),jwtSecret))
	mux.HandleFunc("GET /login/", auth.RequireGuest(store,handlers.Login,jwtSecret))
	mux.HandleFunc("GET /register/", auth.RequireGuest(store,handlers.Register,jwtSecret))
	mux.HandleFunc("POST /login", handlers.LoginHandler(db,store,jwtSecret))
	mux.HandleFunc("GET /dashboard/", auth.RequireAuth(store, handlers.Dashboard(db), jwtSecret))
	mux.HandleFunc("GET /articles/new", auth.RequireAuth(store,handlers.CreateArticle(db),jwtSecret))
	mux.HandleFunc("POST /articles", auth.RequireAuth(store,handlers.CreateNewArticle(db),jwtSecret))
	mux.HandleFunc("GET /logout",handlers.LogoutHandler(store))
	defer db.Close()
	log.Println("Database is running")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080",mux))

}