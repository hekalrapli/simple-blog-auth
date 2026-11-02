package auth

import (
	"context"
	"net/http"
)


func RequireAuth(store *SessionStore, next http.HandlerFunc, jwtSecret string) http.HandlerFunc {
	return func( w http.ResponseWriter,r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			token := cookie.Value
			userId, ok := store.Get(token)
			if !ok {
				http.Redirect(w,r,"/login/",http.StatusSeeOther)
				return
			}
			ctx := context.WithValue(r.Context(), "userId",userId)
			next(w,r.WithContext(ctx))
			return 
			
		}else {
			cookie_jwt, err := r.Cookie("token")
			if err != nil {
				http.Redirect(w,r,"/login/",http.StatusSeeOther)
				return 
			}
			token := cookie_jwt.Value
			userId,err := Verify(token,jwtSecret)
			if err != nil {
				http.Redirect(w,r,"/login/",http.StatusSeeOther)
				return
			}
			ctx := context.WithValue(r.Context(), "userId",userId)
			next(w,r.WithContext(ctx))
			return 
		
		}

	}

}

func RequireGuest(store *SessionStore, next http.HandlerFunc, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			token := cookie.Value
			_, ok := store.Get(token)
			if !ok {
				next(w,r)
				return
			}
			http.Redirect(w,r,"/dashboard/",http.StatusSeeOther)
			return
		}else {
			cookie,err := r.Cookie("token")
			if err == nil {		
				token := cookie.Value
				_,err := Verify(token,jwtSecret)
				if err != nil {
					next(w,r)
					return
				}		
				http.Redirect(w,r,"/dashboard/",http.StatusSeeOther)
				return
			}else {
				next(w,r)
				return
			}
		}
	}
}

func OptionalAuth(store *SessionStore, next http.HandlerFunc, jwtSecret string) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			token := cookie.Value
			userId, ok := store.Get(token)
			if !ok {
				next(w,r)
				return
			}
			ctx := context.WithValue(r.Context(), "userId",userId)
			next(w,r.WithContext(ctx))
			return
		}else {
			cookie,err := r.Cookie("token")
			if err == nil {		
				token := cookie.Value
				userId,err := Verify(token,jwtSecret)
				if err != nil {
					next(w,r)
					return
				}		
				ctx := context.WithValue(r.Context(), "userId",userId)
				next(w,r.WithContext(ctx))
				return
			}else {
				next(w,r)
				return
			}
		}
	}
}	