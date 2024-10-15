package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./timesheet.db")
	if err != nil {
		return err
	}

	// 创建必要的表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			code TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS employees (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			department TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS timesheets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			employee_id INTEGER,
			project_id INTEGER,
			hours FLOAT,
			date DATE,
			FOREIGN KEY(employee_id) REFERENCES employees(id),
			FOREIGN KEY(project_id) REFERENCES projects(id)
		);
	`)

	return err
}
