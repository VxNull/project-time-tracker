package models

import (
	"github.com/VxNull/project-time-tracker/database"
	"golang.org/x/crypto/bcrypt"
)

type Employee struct {
	ID         int
	Name       string
	Username   string
	Password   string
	Department string
}

func CreateEmployee(name, username, password, department string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec("INSERT INTO employees (name, username, password, department) VALUES (?, ?, ?, ?)",
		name, username, string(hashedPassword), department)
	return err
}

func GetEmployeeByUsername(username string) (*Employee, error) {
	var e Employee
	err := database.DB.QueryRow("SELECT id, name, username, password, department FROM employees WHERE username = ?", username).
		Scan(&e.ID, &e.Name, &e.Username, &e.Password, &e.Department)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func GetEmployeeByID(id int) (*Employee, error) {
	var e Employee
	err := database.DB.QueryRow("SELECT id, name, username, password, department FROM employees WHERE id = ?", id).
		Scan(&e.ID, &e.Name, &e.Username, &e.Password, &e.Department)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
