package handlers

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/VxNull/project-time-tracker/database"
	"github.com/VxNull/project-time-tracker/models"
	"github.com/VxNull/project-time-tracker/store"
	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("访问管理员登录页面")

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		admin, err := models.GetAdminByUsername(username)
		if err != nil {
			http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
		if err != nil {
			http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
			return
		}

		session, _ := store.Store.Get(r, "session")
		session.Values["admin"] = true
		session.Save(r, w)

		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	log.Println("尝试加载模板")
	tmpl, err := template.ParseFiles("templates/admin_login.html")
	if err != nil {
		log.Printf("模板加载失败: %v", err)
		http.Error(w, "模板加载失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("尝试渲染模板")
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("模板渲染失败: %v", err)
		http.Error(w, "模板渲染失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("管理员登录页面渲染完成")
}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/admin_dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ManageProject(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		code := r.FormValue("code")
		err := models.CreateProject(name, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	projects, err := models.GetAllProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/manage_project.html"))
	tmpl.Execute(w, projects)
}

func ManageEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		username := r.FormValue("username")
		password := r.FormValue("password")
		department := r.FormValue("department")
		err := models.CreateEmployee(name, username, password, department)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/manage_employee.html"))
	tmpl.Execute(w, nil)
}

func ExportTimesheet(w http.ResponseWriter, r *http.Request) {
	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")

	// 从数据库获取工时数据
	rows, err := database.DB.Query(`
		SELECT e.name, p.name, t.hours, t.date
		FROM timesheets t
		JOIN employees e ON t.employee_id = e.id
		JOIN projects p ON t.project_id = p.id
		WHERE t.date BETWEEN ? AND ?
	`, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 设置响应头,使浏览器下载文件
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=timesheet_%s_%s.csv", startDate, endDate))

	// 创建CSV writer
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// 写入CSV头
	csvWriter.Write([]string{"员工姓名", "项目名称", "工时", "日期"})

	// 写入数据
	for rows.Next() {
		var employeeName, projectName string
		var hours float64
		var date time.Time
		if err := rows.Scan(&employeeName, &projectName, &hours, &date); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		csvWriter.Write([]string{
			employeeName,
			projectName,
			fmt.Sprintf("%.2f", hours),
			date.Format("2006-01-02"),
		})
	}
}
