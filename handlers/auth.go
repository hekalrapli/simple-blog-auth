package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"simple-blog-auth/auth"
	"simple-blog-auth/database"

	"golang.org/x/crypto/bcrypt"
)



func CreateAccount(db *sql.DB)  http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		err  := r.ParseForm()

	if err != nil {
		log.Fatalf("Failed to read the value:%v",err)
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	confirm_password := r.PostFormValue("confirm_password")
	if password != confirm_password {
		render(w, "register", PageData{
			Title:        "Register",
			ConfirmError: "Password tidak sesuai",
			FormUsername : username,
			FormPassword: password,
			FormConfirmPassword: confirm_password,
		})
		return
	}
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword,bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash the password:%v",err)
	}
	if err := database.Insert(db, username,string(hashedPassword)); err != nil {
		render(w, "register", PageData{
			Title:        "Register",
			Error: "Gagal mendaftar: " + err.Error(),
			FormUsername : username,
			FormPassword: password,
			FormConfirmPassword: confirm_password,
		})
		return 
	}

	render(w, "register", PageData{
			Title:        "Register",
			Message: "Berhasil mendaftar silakan login",
		
	})
	// http.Redirect(w,r,"/login/",http.StatusSeeOther)
	}
	
	
}

func LoginHandler(db *sql.DB, store *auth.SessionStore,jwtSecret string ) http.HandlerFunc {
	return func(w http.ResponseWriter, r*http.Request) {
		err := r.ParseForm()
	if err != nil {
		log.Fatalf("Failed to read the value:%v",err)
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	mode := r.PostFormValue("mode")

	bytePassword := []byte(password)
	userId,HashPassword, err := database.Select(db,username)
	if err != nil {
		render(w, "login", PageData{
			Title:        "Login",
			Error: "Username atau Password salah",
			FormUsername: username,
		})
		return
	}
	byteHashPassword := []byte(HashPassword)
	
	if err := bcrypt.CompareHashAndPassword(byteHashPassword, bytePassword); err != nil {
		render(w, "login", PageData{
			Title:        "Login",
			Error: "Username atau Password salah",
			FormUsername: username,
	})
		return 
	}
	if(mode == "session") {
		token,err := store.Create(userId)
			if err != nil {
				log.Fatalf("Failed to create token:%v",err)

			}
	http.SetCookie(w, &http.Cookie{
		Name : "session_id",
		Value : token,
		Path:  "/",
		HttpOnly: true,
	})
	http.Redirect(w,r,"/dashboard/",http.StatusSeeOther)
	}else {
		token, err := auth.Generate(userId, jwtSecret)
		if err != nil {
			log.Fatalf("Failed to login: %v",err)
		}
		
		http.SetCookie(w, &http.Cookie{
			Name: "token",
			Value : token,
			Path : "/",
			HttpOnly: true,
			MaxAge: 900,
		})

		http.Redirect(w,r,"/dashboard/",http.StatusSeeOther)


	}
	

	

	}
}


func LogoutHandler(store *auth.SessionStore) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			token := cookie.Value
			store.Delete(token)
			http.SetCookie(w, &http.Cookie{
			Name : "session_id",
			Value : "",
			Path:  "/",
			HttpOnly: true,
			MaxAge: -1,
		})
		} 
			http.SetCookie(w, &http.Cookie{
				Name : "token",
				Value : "",
				Path : "/",
				HttpOnly: true,
				MaxAge: -1,
			})
		
	http.Redirect(w,r,"/login/",http.StatusSeeOther)
	}
}
