package handlers

import (
	"encoding/json"
	"html/template"
	"log"
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
		if err != nil || bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password)) != nil {
			// 鉴权失败，跳转到首页
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		session, _ := store.Store.Get(r, "session")
		session.Values["employee_id"] = employee.ID
		session.Save(r, w)

		http.Redirect(w, r, "/employee/dashboard", http.StatusSeeOther)
		return
	}

	// 直接跳转到首页
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Error(w, "获取员工信息失败", http.StatusInternalServerError)
		return
	}

	projects, err := models.GetAllProjects()
	if err != nil {
		http.Error(w, "获取项目列表失败", http.StatusInternalServerError)
		return
	}

	timesheets, err := models.GetTimesheetsByEmployee(employeeID)
	if err != nil {
		http.Error(w, "获取工时记录失败", http.StatusInternalServerError)
		return
	}

	type TimesheetWithProjectName struct {
		models.Timesheet
		ProjectName string
	}

	var timesheetsWithProjectNames []TimesheetWithProjectName
	for _, ts := range timesheets {
		project, err := models.GetProjectByID(ts.ProjectID)
		if err != nil {
			http.Error(w, "获取项目信息失败", http.StatusInternalServerError)
			return
		}
		timesheetsWithProjectNames = append(timesheetsWithProjectNames, TimesheetWithProjectName{
			Timesheet:   ts,
			ProjectName: project.Name,
		})
	}

	data := struct {
		Employee          *models.Employee
		Projects          []models.Project
		Timesheets        []TimesheetWithProjectName
		CurrentMonth      string
		MonthlyHours      []models.MonthlyProjectHours
		TotalMonthlyHours float64
	}{
		Employee:     employee,
		Projects:     projects,
		Timesheets:   timesheetsWithProjectNames,
		CurrentMonth: time.Now().Format("2006-01"),
	}

	// 获取当前月份的工时统计
	currentMonth := time.Now().UTC().Truncate(time.Hour * 24 * 30)
	monthlyHours, totalHours, err := models.GetEmployeeMonthlyHours(employeeID, currentMonth)
	if err != nil {
		log.Printf("获取月度工时统计失败: %v", err)
	} else {
		data.MonthlyHours = monthlyHours
		data.TotalMonthlyHours = totalHours
	}

	tmpl, err := template.ParseFiles("templates/employee_dashboard.html")
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
	month, _ := time.Parse("2006-01", r.FormValue("month"))
	description := r.FormValue("description")

	err := models.SubmitTimesheet(employeeID, projectID, hours, month, description)
	if err != nil {
		http.Error(w, "Failed to submit timesheet: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/employee/dashboard", http.StatusSeeOther)
}

func GetEmployeeMonthlyHours(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Store.Get(r, "session")
	employeeID, ok := session.Values["employee_id"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	monthStr := r.URL.Query().Get("month")
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		http.Error(w, "Invalid month format", http.StatusBadRequest)
		return
	}

	log.Printf("获取员工 %d 在 %s 月份的工时统计", employeeID, monthStr)

	projectHours, totalHours, err := models.GetEmployeeMonthlyHours(employeeID, month)
	if err != nil {
		log.Printf("获取月度工时统计失败: %v", err)
		// 即使出错，也返回一个空的有效响应
		projectHours = []models.MonthlyProjectHours{}
		totalHours = 0
	}

	log.Printf("找到 %d 个项目的工时记录，总工时: %.2f", len(projectHours), totalHours)

	response := struct {
		ProjectHours []models.MonthlyProjectHours `json:"projectHours"`
		TotalHours   float64                      `json:"totalHours"`
	}{
		ProjectHours: projectHours,
		TotalHours:   totalHours,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func EmployeeLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Store.Get(r, "session")
	session.Values["employee_id"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/employee/login", http.StatusSeeOther)
}

func UpdateTimesheet(w http.ResponseWriter, r *http.Request) {
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

	timesheetID := r.URL.Path[len("/employee/update/"):]
	projectID, _ := strconv.Atoi(r.FormValue("project_id"))
	hours, _ := strconv.ParseFloat(r.FormValue("hours"), 64)
	month, _ := time.Parse("2006-01", r.FormValue("month"))
	description := r.FormValue("description")

	err := models.UpdateTimesheet(timesheetID, employeeID, projectID, hours, month, description)
	if err != nil {
		http.Error(w, "Failed to update timesheet: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/employee/dashboard", http.StatusSeeOther)
}
