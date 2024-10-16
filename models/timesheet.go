package models

import (
	"fmt"
	"time"

	"github.com/VxNull/project-time-tracker/database"
)

type Timesheet struct {
	ID          int
	EmployeeID  int
	ProjectID   int
	Hours       float64
	Month       time.Time // 修改为月份
	Description string
	SubmitDate  time.Time // 新增提交日期字段
}

func SubmitTimesheet(employeeID, projectID int, hours float64, month time.Time, description string) error {
	_, err := database.DB.Exec(`
		INSERT INTO timesheets (employee_id, project_id, hours, month, description, submit_date)
		VALUES (?, ?, ?, ?, ?, ?)
	`, employeeID, projectID, hours, month, description, time.Now())
	return err
}

func GetTimesheetsByEmployee(employeeID int) ([]Timesheet, error) {
	rows, err := database.DB.Query(`
		SELECT id, employee_id, project_id, hours, month, description, submit_date
		FROM timesheets
		WHERE employee_id = ?
		ORDER BY month DESC, submit_date DESC
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timesheets []Timesheet
	for rows.Next() {
		var t Timesheet
		if err := rows.Scan(&t.ID, &t.EmployeeID, &t.ProjectID, &t.Hours, &t.Month, &t.Description, &t.SubmitDate); err != nil {
			return nil, err
		}
		timesheets = append(timesheets, t)
	}
	return timesheets, nil
}

func GetCurrentMonthTotalHours() (float64, error) {
	var totalHours float64
	currentYear, currentMonth, _ := time.Now().Date()
	currentMonthStart := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local)

	err := database.DB.QueryRow(`
		SELECT COALESCE(SUM(hours), 0) 
		FROM timesheets 
		WHERE month = ?
	`, currentMonthStart.Format("2006-01")).Scan(&totalHours)

	if err != nil {
		return 0, fmt.Errorf("计算本月总工时失败: %v", err)
	}

	return totalHours, nil
}

type MonthlyProjectHours struct {
	ProjectID   int
	ProjectName string
	Hours       float64
}

func GetEmployeeMonthlyHours(employeeID int, month time.Time) ([]MonthlyProjectHours, float64, error) {
	// 计算月份的开始和结束时间
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	rows, err := database.DB.Query(`
		SELECT t.project_id, p.name, SUM(t.hours)
		FROM timesheets t
		JOIN projects p ON t.project_id = p.id
		WHERE t.employee_id = ? AND t.month >= ? AND t.month <= ?
		GROUP BY t.project_id, p.name
	`, employeeID, startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02"))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projectHours []MonthlyProjectHours
	var totalHours float64

	for rows.Next() {
		var ph MonthlyProjectHours
		if err := rows.Scan(&ph.ProjectID, &ph.ProjectName, &ph.Hours); err != nil {
			return nil, 0, err
		}
		projectHours = append(projectHours, ph)
		totalHours += ph.Hours
	}

	return projectHours, totalHours, nil
}
