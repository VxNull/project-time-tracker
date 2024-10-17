package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/VxNull/project-time-tracker/models"
	"github.com/VxNull/project-time-tracker/store"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

func AdminLogout(w http.ResponseWriter, r *http.Request) {
	// 清除会话
	session, _ := store.Store.Get(r, "session")
	delete(session.Values, "admin")
	session.Save(r, w)

	// 重定向到首页
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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

	// 生成月份选项
	months := []string{}
	for year := 2020; year <= time.Now().Year(); year++ {
		for month := 1; month <= 12; month++ {
			months = append(months, strconv.Itoa(year)+"-"+fmt.Sprintf("%02d", month))
		}
	}

	// 渲染模板
	tmpl := template.Must(template.ParseFiles("templates/admin_dashboard.html"))
	tmpl.Execute(w, map[string]interface{}{
		"ProjectCount":      projectCount,
		"EmployeeCount":     employeeCount,
		"CurrentMonthHours": currentMonthHours,
		"Months":            months,
	})
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
	if r.Method == "POST" {
		startMonth := r.FormValue("start_month")
		endMonth := r.FormValue("end_month")

		// 解析开始和结束月份
		start, _ := time.Parse("2006-01", startMonth)
		end, _ := time.Parse("2006-01", endMonth)

		// 获取工时数据
		timesheets, err := models.GetTimesheetsByDateRange(start, end)
		if err != nil {
			http.Error(w, "获取工时数据失败", http.StatusInternalServerError)
			return
		}

		// 创建 Excel 文件
		f := excelize.NewFile()
		sheetName := "工时数据"
		f.NewSheet(sheetName)

		// 写入表头
		f.SetCellValue(sheetName, "A1", "员工ID")
		f.SetCellValue(sheetName, "B1", "项目ID")
		f.SetCellValue(sheetName, "C1", "工时")
		f.SetCellValue(sheetName, "D1", "月份")
		f.SetCellValue(sheetName, "E1", "描述")

		// 写入数据
		for i, ts := range timesheets {
			f.SetCellValue(sheetName, "A"+strconv.Itoa(i+2), ts.EmployeeID)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(i+2), ts.ProjectID)
			f.SetCellValue(sheetName, "C"+strconv.Itoa(i+2), ts.Hours)
			f.SetCellValue(sheetName, "D"+strconv.Itoa(i+2), ts.Month.Format("2006-01"))
			f.SetCellValue(sheetName, "E"+strconv.Itoa(i+2), ts.Description)
		}

		// 设置响应头
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename=timesheet.xlsx")

		// 写入文件到响应
		if err := f.Write(w); err != nil {
			http.Error(w, "导出 Excel 失败", http.StatusInternalServerError)
			return
		}
		return
	}

	// 处理 GET 请求，渲染导出页面
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func GetTimesheetData(w http.ResponseWriter, r *http.Request) {
	startMonth := r.URL.Query().Get("start_month")
	endMonth := r.URL.Query().Get("end_month")

	start, _ := time.Parse("2006-01", startMonth)
	end, _ := time.Parse("2006-01", endMonth)

	// 获取工时数据
	timesheets, err := models.GetTimesheetsByDateRange(start, end)
	if err != nil {
		http.Error(w, "获取工时数据失败", http.StatusInternalServerError)
		return
	}

	// 统计每个项目的工时
	projectHours := make(map[string]float64)
	for _, ts := range timesheets {
		project, err := models.GetProjectByID(ts.ProjectID)
		if err == nil {
			projectHours[project.Name] += ts.Hours
		}
	}

	// 转换为可返回的格式
	var result []struct {
		ProjectName string  `json:"projectName"`
		Hours       float64 `json:"hours"`
	}

	for projectName, hours := range projectHours {
		result = append(result, struct {
			ProjectName string  `json:"projectName"`
			Hours       float64 `json:"hours"`
		}{
			ProjectName: projectName,
			Hours:       hours,
		})
	}

	// 按工时从高到低排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Hours > result[j].Hours
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
