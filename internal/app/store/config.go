package store

import (
	"fmt"
	"os"
)

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

type Config struct {
	DatabaseURL string `os:"database_url"`
}

func NewConfig() *Config {
	config := dbConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport], config[dbuser], config[dbpass], config[dbname])
	return &Config{
		DatabaseURL: psqlInfo,
	}
}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}
