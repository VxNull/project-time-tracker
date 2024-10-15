package handlers

import (
	"html/template"
	"net/http"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// 这里应该有更安全的认证方式,比如使用bcrypt加密密码
		if username == "admin" && password == "password" {
			http.SetCookie(w, &http.Cookie{
				Name:  "admin",
				Value: "true",
				Path:  "/",
			})
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/admin_login.html"))
	tmpl.Execute(w, nil)
}
