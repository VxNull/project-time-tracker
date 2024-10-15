package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/VxNull/project-time-tracker/models"
	"github.com/VxNull/project-time-tracker/store"
	"golang.org/x/crypto/bcrypt"
)

func EmployeeLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		employee, err := models.GetEmployeeByUsername(username)
		if err != nil {
			http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password)); err != nil {
			http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
			return
		}

		session, _ := store.Store.Get(r, "session")
		session.Values["employee_id"] = employee.ID
		session.Save(r, w)

		http.Redirect(w, r, "/employee/dashboard", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/employee_login.html")
	if err != nil {
		http.Error(w, "模板加载失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "模板渲染失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func EmployeeDashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Store.Get(r, "session")
	employeeID, ok := session.Values["employee_id"].(int)
	if !ok {
		http.Redirect(w, r, "/employee/login", http.StatusSeeOther)
		return
	}

	employee, err := models.GetEmployeeByID(employeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projects, err := models.GetAllProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timesheets, err := models.GetTimesheetsByEmployee(employeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Employee   *models.Employee
		Projects   []models.Project
		Timesheets []models.Timesheet
	}{
		Employee:   employee,
		Projects:   projects,
		Timesheets: timesheets,
	}

	tmpl := template.Must(template.ParseFiles("templates/employee_dashboard.html"))
	tmpl.Execute(w, data)
}

func SubmitTimesheet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Store.Get(r, "session")
	employeeID, ok := session.Values["employee_id"].(int)
	if !ok {
		http.Redirect(w, r, "/employee/login", http.StatusSeeOther)
		return
	}

	projectID, _ := strconv.Atoi(r.FormValue("project_id"))
	hours, _ := strconv.ParseFloat(r.FormValue("hours"), 64)
	date, _ := time.Parse("2006-01-02", r.FormValue("date"))
	description := r.FormValue("description")

	err := models.SubmitTimesheet(employeeID, projectID, hours, date, description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/employee/dashboard", http.StatusSeeOther)
}
