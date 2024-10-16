package models

import (
	"database/sql"
	"fmt"

	"github.com/VxNull/project-time-tracker/database"
	"golang.org/x/crypto/bcrypt"
)

type Employee struct {
	ID           int
	Name         string
	Username     string
	Password     string
	SuperiorID   sql.NullInt64
	SuperiorName sql.NullString
}

func CreateEmployee(name, username, password string, superiorID *int) error {
	// 首先检查用户名是否已存在
	exist, err := IsUsernameExist(username)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("用户名 '%s' 已存在", username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var superiorIDValue sql.NullInt64
	if superiorID != nil {
		superiorIDValue.Int64 = int64(*superiorID)
		superiorIDValue.Valid = true
	}

	_, err = database.DB.Exec("INSERT INTO employees (name, username, password, superior_id) VALUES (?, ?, ?, ?)",
		name, username, string(hashedPassword), superiorIDValue)
	return err
}

func GetEmployeeByUsername(username string) (*Employee, error) {
	var e Employee
	err := database.DB.QueryRow("SELECT id, name, username, password, superior_id FROM employees WHERE username = ?", username).
		Scan(&e.ID, &e.Name, &e.Username, &e.Password, &e.SuperiorID)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func GetEmployeeByID(id int) (*Employee, error) {
	var e Employee
	err := database.DB.QueryRow("SELECT id, name, username, password, superior_id FROM employees WHERE id = ?", id).
		Scan(&e.ID, &e.Name, &e.Username, &e.Password, &e.SuperiorID)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func GetEmployeeCount() (int, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM employees").Scan(&count)
	return count, err
}

func GetAllEmployees() ([]Employee, error) {
	rows, err := database.DB.Query(`
		SELECT e.id, e.name, e.username, COALESCE(s.id, 0) as superior_id, COALESCE(s.name, '') as superior_name
		FROM employees e
		LEFT JOIN employees s ON e.superior_id = s.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Username, &e.SuperiorID, &e.SuperiorName); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}

	return employees, nil
}

func UpdateEmployee(id string, name, username string, superiorID *int) error {
	// 首先检查新用户名是否与其他员工冲突
	var currentUsername string
	err := database.DB.QueryRow("SELECT username FROM employees WHERE id = ?", id).Scan(&currentUsername)
	if err != nil {
		return err
	}

	if username != currentUsername {
		exist, err := IsUsernameExist(username)
		if err != nil {
			return err
		}
		if exist {
			return fmt.Errorf("用户名 '%s' 已存在", username)
		}
	}

	var superiorIDValue sql.NullInt64
	if superiorID != nil {
		superiorIDValue.Int64 = int64(*superiorID)
		superiorIDValue.Valid = true
	}

	_, err = database.DB.Exec("UPDATE employees SET name = ?, username = ?, superior_id = ? WHERE id = ?",
		name, username, superiorIDValue, id)
	return err
}

func DeleteEmployee(id string) error {
	_, err := database.DB.Exec("DELETE FROM employees WHERE id = ?", id)
	return err
}

// 添加这个新函数
func IsUsernameExist(username string) (bool, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM employees WHERE username = ?", username).Scan(&count)
	return count > 0, err
}
