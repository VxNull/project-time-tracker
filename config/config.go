package config

import (
	"log"

	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	Admin struct {
		DefaultUsername string `yaml:"default_username"`
		DefaultPassword string `yaml:"default_password"`
	} `yaml:"admin"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Session struct {
		SecretKey string `yaml:"secret_key"`
	} `yaml:"session"`
}

var AppConfig Config

func LoadConfig(configPath string) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal("无法打开配置文件:", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Fatal("无法解析配置文件:", err)
	}
}
