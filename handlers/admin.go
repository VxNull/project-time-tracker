package handlers

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/VxNull/project-time-tracker/database"
	"github.com/VxNull/project-time-tracker/models"
	"github.com/VxNull/project-time-tracker/store"
	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		admin, err := models.GetAdminByUsername(username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)) != nil {
			// 鉴权失败，跳转到首页并传递错误信息
			http.Redirect(w, r, "/?error=用户名或密码错误", http.StatusSeeOther)
			return
		}

		session, _ := store.Store.Get(r, "session")
		session.Values["admin"] = true
		session.Save(r, w)

		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	// 直接跳转到首页
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	projectCount, err := models.GetProjectCount()
	if err != nil {
		http.Error(w, "获取项目数量失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	employeeCount, err := models.GetEmployeeCount()
	if err != nil {
		http.Error(w, "获取员工数量失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	currentMonthHours, err := models.GetCurrentMonthTotalHours()
	if err != nil {
		http.Error(w, "获取本月总工时失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		ProjectCount      int
		EmployeeCount     int
		CurrentMonthHours float64
	}{
		ProjectCount:      projectCount,
		EmployeeCount:     employeeCount,
		CurrentMonthHours: currentMonthHours,
	}

	tmpl, err := template.ParseFiles("templates/admin_dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ManageProject(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		action := r.FormValue("action")
		switch action {
		case "add":
			name := r.FormValue("name")
			code := r.FormValue("code")

			nameExist, err := models.IsProjectNameExist(name)
			if err != nil {
				http.Error(w, "检查项目名称失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if nameExist {
				http.Error(w, "项目名称已存在", http.StatusBadRequest)
				return
			}

			codeExist, err := models.IsProjectCodeExist(code)
			if err != nil {
				http.Error(w, "检查项目代码失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if codeExist {
				http.Error(w, "项目代码已存在", http.StatusBadRequest)
				return
			}

			err = models.CreateProject(name, code)
			if err != nil {
				http.Error(w, "创建项目失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
		case "edit":
			id := r.FormValue("id")
			name := r.FormValue("name")
			code := r.FormValue("code")
			err := models.UpdateProject(id, name, code)
			if err != nil {
				http.Error(w, "更新项目失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
		case "delete":
			id := r.FormValue("id")
			err := models.DeleteProject(id)
			if err != nil {
				http.Error(w, "删除项目失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/admin/project", http.StatusSeeOther)
		return
	}

	projects, err := models.GetAllProjects()
	if err != nil {
		log.Printf("获取项目列表失败: %v", err)
		http.Error(w, "获取项目列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("获取到 %d 个项目", len(projects))

	tmpl, err := template.ParseFiles("templates/manage_project.html")
	if err != nil {
		log.Printf("解析模板失败: %v", err)
		http.Error(w, "解析模板失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, projects)
	if err != nil {
		log.Printf("渲染模板失败: %v", err)
		http.Error(w, "渲染模板失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func ManageEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		action := r.FormValue("action")
		switch action {
		case "add":
			name := r.FormValue("name")
			username := r.FormValue("username")
			password := r.FormValue("password")
			superiorID := r.FormValue("superior_id")

			var superiorIDPtr *int
			if superiorID != "" {
				id, _ := strconv.Atoi(superiorID)
				superiorIDPtr = &id
			}

			err := models.CreateEmployee(name, username, password, superiorIDPtr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case "edit":
			id := r.FormValue("id")
			name := r.FormValue("name")
			username := r.FormValue("username")
			superiorID := r.FormValue("superior_id")

			var superiorIDPtr *int
			if superiorID != "" {
				id, _ := strconv.Atoi(superiorID)
				superiorIDPtr = &id
			}

			err := models.UpdateEmployee(id, name, username, superiorIDPtr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case "delete":
			id := r.FormValue("id")
			err := models.DeleteEmployee(id)
			if err != nil {
				http.Error(w, "删除员工失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
		case "reset_password":
			id := r.FormValue("id")
			newPassword := r.FormValue("new_password")
			err := models.ResetEmployeePassword(id, newPassword)
			if err != nil {
				http.Error(w, "重置密码失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("密码重置成功"))
			return
		}
		http.Redirect(w, r, "/admin/employee", http.StatusSeeOther)
		return
	}

	employees, err := models.GetAllEmployees()
	if err != nil {
		http.Error(w, "获取员工列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/manage_employee.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, employees)
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
