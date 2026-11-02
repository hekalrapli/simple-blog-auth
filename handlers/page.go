package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"simple-blog-auth/database"
	"simple-blog-auth/models"
	"strconv"
)

type PageData struct {
	Title        string
	IsLoggedIn   bool
	FormUsername   string
	FormPassword string
	FormConfirmPassword string
	Error        string
	Message      string
	ConfirmError string
	Username string
	MyArticles []models.Articles 
	Article *models.Articles
	Articles []models.Articles
	IsOwner     bool

}

func render(w http.ResponseWriter, page string, data any) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/"+page+".html",
	)
	if err != nil {
		log.Println("template.ParseFiles error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("tmpl.ExecuteTemplate error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	render(w, "login", nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	render(w, "register", nil)
}

func Home(db *sql.DB) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		userId, _ := r.Context().Value("userId").(int64)
		articles,err := database.GetAllArticles(db)
		if err != nil {
			log.Printf("Failed to get articles: %v",err)
			return 

		}
		if userId != 0 {
			username,err := database.GetUsername(db,userId)
			if err != nil {
			log.Printf("Failed to get username: %v",err)
				return
			}
			render(w, "home", PageData{
			Title: "home",
			IsLoggedIn: true,
			Username: username,
			Articles: articles,
		})
		}else {

			render(w, "home", PageData{
				Title: "home",
				Articles: articles,
			})
		}
		


	
		
	
		
	}
}

func Dashboard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("userId").(int64)
	username,err := database.GetUsername(db,userId)
	if err != nil {
		log.Fatalf("Failed to get username: %v",err)
		return 
	}
	articles,err := database.GetArticles(db,userId)
	render(w, "dashboard", PageData{
		Title: "Dashboard",
		IsLoggedIn: true,
		Username: username,
		MyArticles: articles,
	})
	}
	
}

func ViewArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId,_ := r.Context().Value("userId").(int64)
		articleId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Failed to convert articleId: %v",err)
			return
		}
		if userId != 0 {
			username,err := database.GetUsername(db,userId)
			if err != nil {
				log.Printf("Failed to get username: %v",err)
				return
			}
			article, err := database.GetArticle(db,int64(articleId))
			if err != nil {
				render(w, "view", PageData{
				Title: "View",
				IsLoggedIn: true,
				Username: username,
				Article: nil,
			})	
				log.Printf("Failed to get article: %v", err)
				return
			}
			render(w,"view",PageData{
			Title: "edit",
			IsLoggedIn: true,
			IsOwner: article.AuthorID == userId,
			Username: username,
			Article: article,
		})
		}else {
			article, err := database.GetArticle(db,int64(articleId))
			if err != nil {
				render(w, "view", PageData{
				Title: "View",
				Article: nil,
			})	
				log.Printf("Failed to get article: %v", err)
				return
			}
			render(w, "view", PageData{
					Title: "View",
					Article: article,
				})	
		}


		
		
		
		
		
		
			
		
	}
}

