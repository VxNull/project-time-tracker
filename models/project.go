package models

import (
	"github.com/VxNull/project-time-tracker/database"
)

type Project struct {
	ID   int
	Name string
	Code string
}

func CreateProject(name, code string) error {
	_, err := database.DB.Exec("INSERT INTO projects (name, code) VALUES (?, ?)", name, code)
	return err
}

func UpdateProject(id, name, code string) error {
	_, err := database.DB.Exec("UPDATE projects SET name = ?, code = ? WHERE id = ?", name, code, id)
	return err
}

func DeleteProject(id string) error {
	_, err := database.DB.Exec("DELETE FROM projects WHERE id = ?", id)
	return err
}

func GetAllProjects() ([]Project, error) {
	rows, err := database.DB.Query("SELECT id, name, code FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Code); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func GetProjectCount() (int, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	return count, err
}
