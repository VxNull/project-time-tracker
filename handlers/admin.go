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

		start, _ := time.Parse("2006-01", startMonth)
		end, _ := time.Parse("2006-01", endMonth)

		// 获取时间范围内的所有月份
		months := getMonthsBetween(start, end)

		f := excelize.NewFile()

		// 创建综合统计工作表
		summarySheet := "综合统计"
		f.NewSheet(summarySheet)
		// 删除默认的Sheet1
		f.DeleteSheet("Sheet1")

		// 设置列宽
		f.SetColWidth(summarySheet, "A", "A", 20)
		f.SetColWidth(summarySheet, "B", "Z", 15)

		// 设置综合统计表头
		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Size: 12},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DDEBF7"}, Pattern: 1},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 1},
				{Type: "top", Color: "#000000", Style: 1},
				{Type: "bottom", Color: "#000000", Style: 1},
				{Type: "right", Color: "#000000", Style: 1}},
		})
		f.SetCellStyle(summarySheet, "A1", "B1", headerStyle)
		f.SetCellValue(summarySheet, "A1", "统计项目")
		f.SetCellValue(summarySheet, "B1", "数值")

		// 获取项目和员工总数
		projectCount, _ := models.GetProjectCount()
		employeeCount, _ := models.GetEmployeeCount()

		// 设置项目和员工总数
		f.SetCellValue(summarySheet, "A2", "项目总数")
		f.SetCellValue(summarySheet, "B2", projectCount)
		f.SetCellValue(summarySheet, "A3", "员工总数")
		f.SetCellValue(summarySheet, "B3", employeeCount)

		// 计算总工时
		var totalHours float64
		projectTotalHours := make(map[int]float64)
		employeeTotalHours := make(map[int]float64)
		monthlyHours := make(map[string]float64)
		projectMonthlyHours := make(map[int]map[string]float64)
		employeeMonthlyHours := make(map[int]map[string]float64)

		for _, month := range months {
			sheetName := month.Format("2006-01")
			f.NewSheet(sheetName)

			// 获取该月的工时数据
			timesheets, err := models.GetTimesheetsByMonth(month)
			if err != nil {
				http.Error(w, "获取工时数据失败", http.StatusInternalServerError)
				return
			}

			// 获取所有项目和员工
			projects, _ := models.GetAllProjects()
			employees, _ := models.GetAllEmployees()

			// 设置表头
			f.SetCellValue(sheetName, "A1", "员工姓名")
			for col, project := range projects {
				f.SetCellValue(sheetName, getColumnName(col+1+1)+"1", project.Name+" ("+project.Code+")")
			}

			// 填充数据
			monthTotal := 0.0
			for row, employee := range employees {
				f.SetCellValue(sheetName, "A"+strconv.Itoa(row+2), employee.Name)
				employeeMonthTotal := 0.0
				for col, project := range projects {
					hours := getHours(timesheets, employee.ID, project.ID)
					f.SetCellValue(sheetName, getColumnName(col+1+1)+strconv.Itoa(row+2), hours)
					totalHours += hours
					projectTotalHours[project.ID] += hours
					employeeTotalHours[employee.ID] += hours
					monthTotal += hours
					employeeMonthTotal += hours

					// 记录每个项目每月的工时
					if _, ok := projectMonthlyHours[project.ID]; !ok {
						projectMonthlyHours[project.ID] = make(map[string]float64)
					}
					projectMonthlyHours[project.ID][sheetName] += hours
				}
				// 记录每个员工每月的工时
				if _, ok := employeeMonthlyHours[employee.ID]; !ok {
					employeeMonthlyHours[employee.ID] = make(map[string]float64)
				}
				employeeMonthlyHours[employee.ID][sheetName] = employeeMonthTotal
			}

			monthlyHours[sheetName] = monthTotal

			// 添加汇总统计
			totalRow := len(employees) + 3
			f.SetCellValue(sheetName, "A"+strconv.Itoa(totalRow-1), "项目总计")
			for col := range projects {
				colName := getColumnName(col + 1 + 1)
				f.SetCellFormula(sheetName, colName+strconv.Itoa(totalRow-1), "SUM("+colName+"2:"+colName+strconv.Itoa(totalRow-1-1)+")")
			}

			// 设置总计
			f.SetCellValue(sheetName, "A"+strconv.Itoa(totalRow), "总计")
			f.SetCellFormula(sheetName, "B"+strconv.Itoa(totalRow), "SUM(B"+strconv.Itoa(totalRow-1)+":"+getColumnName(len(projects)+1)+strconv.Itoa(totalRow-1)+")")
		}

		// 在综合统计工作表中添加总工时
		f.SetCellValue(summarySheet, "A4", "总工时")
		f.SetCellValue(summarySheet, "B4", totalHours)

		// 添加项目工时统计（二维表格）
		f.SetCellValue(summarySheet, "A6", "项目工时统计")
		row := 7

		// 设置表头
		f.SetCellStyle(summarySheet, "A"+strconv.Itoa(row), getColumnName(len(months)+2)+strconv.Itoa(row), headerStyle)
		f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), "项目名称")
		for col, month := range months {
			f.SetCellValue(summarySheet, getColumnName(col+1+1)+strconv.Itoa(row), month.Format("2006-01"))
		}
		f.SetCellValue(summarySheet, getColumnName(len(months)+2)+strconv.Itoa(row), "项目总计")
		row++

		// 填充项目工时数据
		contentStyle, _ := f.NewStyle(&excelize.Style{
			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 1},
				{Type: "top", Color: "#000000", Style: 1},
				{Type: "bottom", Color: "#000000", Style: 1},
				{Type: "right", Color: "#000000", Style: 1}},
		})
		projects, _ := models.GetAllProjects()
		for _, project := range projects {
			f.SetCellStyle(summarySheet, "A"+strconv.Itoa(row), getColumnName(len(months)+2)+strconv.Itoa(row), contentStyle)
			f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), project.Name)
			totalProjectHours := 0.0
			for col, month := range months {
				monthStr := month.Format("2006-01")
				hours := projectMonthlyHours[project.ID][monthStr]
				f.SetCellValue(summarySheet, getColumnName(col+1+1)+strconv.Itoa(row), hours)
				totalProjectHours += hours
			}
			f.SetCellValue(summarySheet, getColumnName(len(months)+2)+strconv.Itoa(row), totalProjectHours)
			row++
		}

		// 添加月度总计行
		f.SetCellStyle(summarySheet, "A"+strconv.Itoa(row), getColumnName(len(months)+2)+strconv.Itoa(row), headerStyle)
		f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), "月度总计")
		for col, month := range months {
			monthStr := month.Format("2006-01")
			f.SetCellValue(summarySheet, getColumnName(col+1+1)+strconv.Itoa(row), monthlyHours[monthStr])
		}
		f.SetCellValue(summarySheet, getColumnName(len(months)+2)+strconv.Itoa(row), totalHours)

		// 添加员工工时统计（二维表格）
		row += 2
		f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), "员工工时统计")
		row++

		// 设置表头
		f.SetCellStyle(summarySheet, "A"+strconv.Itoa(row), getColumnName(len(months)+2)+strconv.Itoa(row), headerStyle)
		f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), "员工姓名")
		for col, month := range months {
			f.SetCellValue(summarySheet, getColumnName(col+1+1)+strconv.Itoa(row), month.Format("2006-01"))
		}
		f.SetCellValue(summarySheet, getColumnName(len(months)+2)+strconv.Itoa(row), "员工总计")
		row++

		// 填充员工工时数据
		employees, _ := models.GetAllEmployees()
		for _, employee := range employees {
			f.SetCellStyle(summarySheet, "A"+strconv.Itoa(row), getColumnName(len(months)+2)+strconv.Itoa(row), contentStyle)
			f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), employee.Name)
			totalEmployeeHours := 0.0
			for col, month := range months {
				monthStr := month.Format("2006-01")
				hours := employeeMonthlyHours[employee.ID][monthStr]
				f.SetCellValue(summarySheet, getColumnName(col+1+1)+strconv.Itoa(row), hours)
				totalEmployeeHours += hours
			}
			f.SetCellValue(summarySheet, getColumnName(len(months)+2)+strconv.Itoa(row), totalEmployeeHours)
			row++
		}

		// 添加月度总计行
		f.SetCellStyle(summarySheet, "A"+strconv.Itoa(row), getColumnName(len(months)+2)+strconv.Itoa(row), headerStyle)
		f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), "月度总计")
		for col, month := range months {
			monthStr := month.Format("2006-01")
			f.SetCellValue(summarySheet, getColumnName(col+1+1)+strconv.Itoa(row), monthlyHours[monthStr])
		}
		f.SetCellValue(summarySheet, getColumnName(len(months)+2)+strconv.Itoa(row), totalHours)

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

// 辅助函数
func getMonthsBetween(start, end time.Time) []time.Time {
	var months []time.Time
	for current := start; current.Before(end) || current.Equal(end); current = current.AddDate(0, 1, 0) {
		months = append(months, current)
	}
	return months
}

func getColumnName(index int) string {
	name := ""
	for index > 0 {
		index--
		name = string(rune('A'+index%26)) + name
		index /= 26
	}
	return name
}

func getHours(timesheets []models.Timesheet, employeeID, projectID int) float64 {
	var total float64
	for _, ts := range timesheets {
		if ts.EmployeeID == employeeID && ts.ProjectID == projectID {
			total += ts.Hours
		}
	}
	return total
}

func GetTimesheetData(w http.ResponseWriter, r *http.Request) {
	startMonth := r.URL.Query().Get("start_month")
	endMonth := r.URL.Query().Get("end_month")

	start, _ := time.Parse("2006-01", startMonth)
	end, _ := time.Parse("2006-01", endMonth)

	// 修改start的日期值为1号，end的日期值为最后一天
	start = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())

	daysInMonth := time.Date(end.Year(), end.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	end = time.Date(end.Year(), end.Month(), daysInMonth, 23, 59, 59, 0, end.Location())

	// 获取工时数据
	timesheets, err := models.GetTimesheetsByDateRange(start, end)
	if err != nil {
		http.Error(w, "获取工时数据失败", http.StatusInternalServerError)
		return
	}

	// 统计每个项目的工时
	projectHours := make(map[string]float64)
	// 统计每个员工的工时
	employeeHours := make(map[string]float64)

	for _, ts := range timesheets {
		project, err := models.GetProjectByID(ts.ProjectID)
		if err == nil {
			projectHours[project.Name] += ts.Hours
		}
		employee, err := models.GetEmployeeByID(ts.EmployeeID)
		if err == nil {
			employeeHours[employee.Name] += ts.Hours
		}
	}

	// 转换为可返回的格式
	var projectResult []struct {
		ProjectName string  `json:"projectName"`
		Hours       float64 `json:"hours"`
	}

	for projectName, hours := range projectHours {
		projectResult = append(projectResult, struct {
			ProjectName string  `json:"projectName"`
			Hours       float64 `json:"hours"`
		}{
			ProjectName: projectName,
			Hours:       hours,
		})
	}

	var employeeResult []struct {
		EmployeeName string  `json:"employeeName"`
		Hours        float64 `json:"hours"`
	}

	for employeeName, hours := range employeeHours {
		employeeResult = append(employeeResult, struct {
			EmployeeName string  `json:"employeeName"`
			Hours        float64 `json:"hours"`
		}{
			EmployeeName: employeeName,
			Hours:        hours,
		})
	}

	// 按工时从高到低排序
	sort.Slice(projectResult, func(i, j int) bool {
		return projectResult[i].Hours > projectResult[j].Hours
	})
	sort.Slice(employeeResult, func(i, j int) bool {
		return employeeResult[i].Hours > employeeResult[j].Hours
	})

	result := struct {
		ProjectHours  interface{} `json:"projectHours"`
		EmployeeHours interface{} `json:"employeeHours"`
	}{
		ProjectHours:  projectResult,
		EmployeeHours: employeeResult,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
