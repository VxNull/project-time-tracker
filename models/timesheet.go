package models

import (
	"time"

	"github.com/VxNull/project-time-tracker/database"
)

type Timesheet struct {
	ID          int
	EmployeeID  int
	ProjectID   int
	Hours       float64
	Date        time.Time
	Description string
}

func SubmitTimesheet(employeeID, projectID int, hours float64, date time.Time, description string) error {
	_, err := database.DB.Exec(`
		INSERT INTO timesheets (employee_id, project_id, hours, date, description)
		VALUES (?, ?, ?, ?, ?)
	`, employeeID, projectID, hours, date, description)
	return err
}

func GetTimesheetsByEmployee(employeeID int) ([]Timesheet, error) {
	rows, err := database.DB.Query(`
		SELECT id, employee_id, project_id, hours, date, description
		FROM timesheets
		WHERE employee_id = ?
		ORDER BY date DESC
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timesheets []Timesheet
	for rows.Next() {
		var t Timesheet
		if err := rows.Scan(&t.ID, &t.EmployeeID, &t.ProjectID, &t.Hours, &t.Date, &t.Description); err != nil {
			return nil, err
		}
		timesheets = append(timesheets, t)
	}
	return timesheets, nil
}

func GetCurrentMonthTotalHours() (float64, error) {
	var totalHours float64
	currentYear, currentMonth, _ := time.Now().Date()
	startOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	err := database.DB.QueryRow(`
		SELECT COALESCE(SUM(hours), 0) 
		FROM timesheets 
		WHERE date >= ? AND date <= ?
	`, startOfMonth, endOfMonth).Scan(&totalHours)

	return totalHours, err
}
