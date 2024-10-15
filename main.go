package main

import (
	"log"
	"net/http"

	"your-project-path/database"
	"your-project-path/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// 初始化数据库
	err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// 管理员路由
	r.HandleFunc("/admin/login", handlers.AdminLogin).Methods("GET", "POST")
	r.HandleFunc("/admin/dashboard", handlers.AdminDashboard).Methods("GET")
	r.HandleFunc("/admin/project", handlers.ManageProject).Methods("GET", "POST")
	r.HandleFunc("/admin/employee", handlers.ManageEmployee).Methods("GET", "POST")
	r.HandleFunc("/admin/export", handlers.ExportTimesheet).Methods("GET")

	// 员工路由
	r.HandleFunc("/employee/login", handlers.EmployeeLogin).Methods("GET", "POST")
	r.HandleFunc("/employee/dashboard", handlers.EmployeeDashboard).Methods("GET")
	r.HandleFunc("/employee/timesheet", handlers.SubmitTimesheet).Methods("POST")

	// 静态文件服务
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
