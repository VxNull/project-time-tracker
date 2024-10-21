package models

import (
	"github.com/VxNull/project-time-tracker/database"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID       int
	Username string
	Password string
}

func CreateAdmin(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec("INSERT INTO admins (username, password) VALUES (?, ?)", username, string(hashedPassword))
	return err
}

func GetAdminByUsername(username string) (*Admin, error) {
	var admin Admin
	err := database.DB.QueryRow("SELECT id, username, password FROM admins WHERE username = ?", username).
		Scan(&admin.ID, &admin.Username, &admin.Password)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func InitDefaultAdmin(username, password string) error {
	// 检查是否已存在管理员账号
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM admins").Scan(&count)
	if err != nil {
		return err
	}

	// 如果不存在管理员账号,则创建默认账号
	if count == 0 {
		return CreateAdmin(username, password)
	}

	return nil
}

func UpdateAdminPassword(adminID int, newPassword string) error {
	_, err := database.DB.Exec("UPDATE admins SET password = ? WHERE id = ?", newPassword, adminID)
	return err
}
