package middleware

import (
	"net/http"

	"github.com/VxNull/project-time-tracker/store"
	"github.com/gorilla/sessions"
)

func InitStore(s *sessions.CookieStore) {
	store.Store = s
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Store.Get(r, "session")
		if _, ok := session.Values["employee_id"].(int); !ok {
			http.Redirect(w, r, "/employee/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Store.Get(r, "session")
		if auth, ok := session.Values["admin"].(bool); !ok || !auth {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
