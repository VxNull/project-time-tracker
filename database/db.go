package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// 创建表的SQL语句
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			superior_id INTEGER,
			FOREIGN KEY (superior_id) REFERENCES employees (id)
		);

		CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			code TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS timesheets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			employee_id INTEGER,
			project_id INTEGER,
			hours REAL NOT NULL,
			month DATE NOT NULL,
			description TEXT,
			submit_date DATETIME NOT NULL,
			FOREIGN KEY (employee_id) REFERENCES employees (id),
			FOREIGN KEY (project_id) REFERENCES projects (id)
		);

		CREATE TABLE IF NOT EXISTS admins (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		);
	`)

	return err
}

func TestConnection() error {
	return DB.Ping()
}
