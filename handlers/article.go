package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"simple-blog-auth/database"
	"strconv"
)

func CreateArticle(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := r.Context().Value("userId").(int64)
		username,err := database.GetUsername(db,userId)
		if err != nil {
			log.Fatalf("Failed to get username:%v",err)
		return 
	}
		render(w, "create", PageData{
			Title: "Create",
			IsLoggedIn: true,
			Username: username,
			Article: nil,
		})
	}
}

func CreateNewArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatalf("Failed to read the value:%v",err)
		}
		title := r.PostFormValue("title")
		content := r.PostFormValue("content")
		userId, _ := r.Context().Value("userId").(int64)
		err = database.InsertArticle(db,title,content,userId)
		if err != nil {
			log.Fatalf("Failed to create article: %v",err)

		}
		http.Redirect(w,r,"/dashboard/",http.StatusSeeOther)

	}
}

func EditArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := r.Context().Value("userId").(int64)
		username,err := database.GetUsername(db,userId)
		if err != nil {
			log.Printf("Failed to get username: %v ",err)

		}
		articleId,err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Failed to convert articleId: %v",err)
			return 
		}
		
		articles,err := database.GetArticle(db,int64(articleId))
		if err != nil {
			log.Printf("Failed to get articles: %v ",err)
		}
		if articles.AuthorID == userId {
			render(w,"create",PageData{
				Title: "edit",
				IsLoggedIn : true,
				IsOwner: true,
				Username: username,
				Article: articles,
			})
		}else {
			http.Error(w, "Akses ditolak: bukan pemilik artikel", http.StatusForbidden)
		}
			
		
	}
}

func UpdateArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to read the form: %v", err)
			return
		}
		userId, _ := r.Context().Value("userId").(int64)
	
		title := r.PostFormValue("title")
		content := r.PostFormValue("content")
		articleId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Failed to convert articleId: %v",err)
			return 
		}

		articles,err := database.GetArticle(db,int64(articleId))
		if err != nil {
			log.Printf("Failed to get articles: %v ",err)
			return
		}
		if articles.AuthorID != userId {
			http.Error(w, "Akses ditolak: bukan pemilik artikel", http.StatusForbidden)
			return
		}
		err = database.UpdateArticles(db,int64(articleId),title,content)
		if err != nil {
			log.Printf("Failed to update article:%v", err)
			return 
		}
		http.Redirect(w,r,"/dashboard/",http.StatusSeeOther)


	}
}

func DeleteArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r*http.Request) {
		userId, _ := r.Context().Value("userId").(int64)
		articleId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Failed to convert articleId: %v",err)
			return
		}
		articles,err := database.GetArticle(db,int64(articleId))
		if err != nil {
			log.Printf("Failed to get articles: %v ",err)
			return
		}
		if articles.AuthorID != userId {
			http.Error(w, "Akses ditolak: bukan pemilik artikel", http.StatusForbidden)
			return
		}
		err = database.DeleteArticle(db, int64(articleId))
		if err != nil {
			log.Printf("Failed to delete article: %v", err) 
			return
		}
		http.Redirect(w,r,"/dashboard/",http.StatusSeeOther)
	}
}