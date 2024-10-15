package models

import (
	"log"

	"github.com/VxNull/project-time-tracker/database"
)

type Project struct {
	ID   int
	Name string
	Code string
}

func CreateProject(name, code string) error {
	_, err := database.DB.Exec("INSERT INTO projects (name, code) VALUES (?, ?)", name, code)
	if err != nil {
		log.Printf("数据库插入项目失败: %v", err)
	}
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
		log.Printf("查询所有项目失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Code); err != nil {
			log.Printf("扫描项目数据失败: %v", err)
			return nil, err
		}
		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("遍历项目行时发生错误: %v", err)
		return nil, err
	}

	return projects, nil
}

func GetProjectCount() (int, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	return count, err
}
