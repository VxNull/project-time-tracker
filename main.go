package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/VxNull/project-time-tracker/config"
	"github.com/VxNull/project-time-tracker/database"
	"github.com/VxNull/project-time-tracker/handlers"
	"github.com/VxNull/project-time-tracker/middleware"
	"github.com/VxNull/project-time-tracker/models"
	"github.com/VxNull/project-time-tracker/store"
)

func main() {
	// 定义命令行参数
	configPath := flag.String("c", "config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置文件
	config.LoadConfig(*configPath)

	// 初始化数据库
	err := database.InitDB(config.AppConfig.Database.Path)
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 初始化默认管理员账号
	err = models.InitDefaultAdmin(config.AppConfig.Admin.DefaultUsername, config.AppConfig.Admin.DefaultPassword)
	if err != nil {
		log.Fatal("默认管理员账号创建失败:", err)
	}

	// 初始化 session store
	store.InitStore(config.AppConfig.Session.SecretKey)

	// 设置路由
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/admin/login", handlers.AdminLogin)
	http.HandleFunc("/admin/dashboard", middleware.AdminAuthMiddleware(handlers.AdminDashboard))
	http.HandleFunc("/admin/logout", middleware.AdminAuthMiddleware(handlers.AdminLogout))
	http.HandleFunc("/admin/project", middleware.AdminAuthMiddleware(handlers.ManageProject))
	http.HandleFunc("/admin/employee", middleware.AdminAuthMiddleware(handlers.ManageEmployee))
	http.HandleFunc("/admin/export", middleware.AdminAuthMiddleware(handlers.ExportTimesheet))
	http.HandleFunc("/admin/get-timesheet-data", middleware.AdminAuthMiddleware(handlers.GetTimesheetData))
	http.HandleFunc("/admin/change-password", middleware.AdminAuthMiddleware(handlers.ChangeAdminPassword))

	http.HandleFunc("/employee/login", handlers.EmployeeLogin)
	http.HandleFunc("/employee/logout", middleware.AuthMiddleware(handlers.EmployeeLogout))
	http.HandleFunc("/employee/dashboard", middleware.AuthMiddleware(handlers.EmployeeDashboard))
	http.HandleFunc("/employee/submit", middleware.AuthMiddleware(handlers.SubmitTimesheet))
	http.HandleFunc("/employee/monthly-hours", middleware.AuthMiddleware(handlers.GetEmployeeMonthlyHours))
	http.HandleFunc("/employee/update/", middleware.AuthMiddleware(handlers.UpdateTimesheet))
	http.HandleFunc("/employee/change-password", middleware.AuthMiddleware(handlers.ChangeEmployeePassword))

	// 设置静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 数据库连接测试
	err = database.TestConnection()
	if err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}
	log.Println("数据库连接测试成功")

	// 启动服务器
	log.Printf("服务器启动在 http://localhost:%d\n", config.AppConfig.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.AppConfig.Server.Port), nil))
}
