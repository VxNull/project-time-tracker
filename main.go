package main

import (
	"log"
	"net/http"

	"github.com/VxNull/project-time-tracker/database"
	"github.com/VxNull/project-time-tracker/handlers"
	"github.com/VxNull/project-time-tracker/middleware"

	"github.com/gorilla/sessions"
)

func main() {
	// 初始化数据库
	err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	store := sessions.NewCookieStore([]byte("your-secret-key"))
	handlers.InitStore(store)

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/admin/login", handlers.AdminLogin)
	http.HandleFunc("/admin/dashboard", middleware.AdminAuthMiddleware(handlers.AdminDashboard))
	http.HandleFunc("/admin/project", middleware.AdminAuthMiddleware(handlers.ManageProject))
	http.HandleFunc("/admin/employee", middleware.AdminAuthMiddleware(handlers.ManageEmployee))
	http.HandleFunc("/admin/export", middleware.AdminAuthMiddleware(handlers.ExportTimesheet))

	http.HandleFunc("/employee/login", handlers.EmployeeLogin)
	http.HandleFunc("/employee/dashboard", middleware.AuthMiddleware(handlers.EmployeeDashboard))
	http.HandleFunc("/employee/submit", middleware.AuthMiddleware(handlers.SubmitTimesheet))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
