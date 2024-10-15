package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func Init(s *sessions.CookieStore) {
	store = s
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		if _, ok := session.Values["employee_id"].(int); !ok {
			http.Redirect(w, r, "/employee/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		if _, ok := session.Values["admin"].(bool); !ok {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
