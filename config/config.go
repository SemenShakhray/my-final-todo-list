package config

import (
	"os"
)

type TodoEnv struct {
	Port     string
	DBFile   string
	Password string
}

func checkEnv(env, baseValue string) string {
	if value, ok := os.LookupEnv(env); ok {
		return value
	}
	return baseValue
}

func GetEnv() *TodoEnv {
	port := checkEnv("TODO_PORT", "7540")
	dbfile := checkEnv("TODO_DBFILE", "")
	password := checkEnv("TODO_PASSWORD", "")

	return &TodoEnv{
		Port:     port,
		DBFile:   dbfile,
		Password: password,
	}
}
