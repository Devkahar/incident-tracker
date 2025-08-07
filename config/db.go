package config

import (
	"fmt"
	"log"

	"incident-tracker/utils"
)

type DBConfig struct {
	Port     string
	Host     string
	Username string
	Password string
	Name     string
}

func LoadDBConfig() DBConfig {
	port, err := utils.GetRequiredEnv("DB_PORT")
	if err != nil {
		log.Fatal(err)
	}
	host, err := utils.GetRequiredEnv("DB_HOST")
	if err != nil {
		log.Fatal(err)
	}
	username, err := utils.GetRequiredEnv("DB_USERNAME")
	if err != nil {
		log.Fatal(err)
	}
	password, err := utils.GetRequiredEnv("DB_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}
	name, err := utils.GetRequiredEnv("DB_NAME")
	if err != nil {
		log.Fatal(err)
	}
	return DBConfig{
		Port:     port,
		Host:     host,
		Username: username,
		Password: password,
		Name:     name,
	}
}

func (c DBConfig) GetDBURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Name)
}
