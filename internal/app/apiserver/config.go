package apiserver

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
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `os:"database_url"` // postgres://user:password@host:port/databaseName?sslmode=[enable/disable]?sslmode=disable
	// postgresql://{user/psswd}ec2-3-217-219-146.compute-1.amazonaws.com:5432/d1d17sa4ihs51r
}

// postgres://postgres:1234567@localhost:8080/deforestation_detection_db?sslmode=disable
// DATABASE_URL

func NewConfig() *Config {
	config := dbConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport], config[dbuser], config[dbpass], config[dbname])
	return &Config{
		BindAddr:    ":8080",
		LogLevel:    "debug",
		DatabaseURL: psqlInfo,
	}
}

func dbConfig() map[string]string {
	// DATABASE CONN URL
	// dbURL := os.Getenv("DATABASE_URL")
	// host=host password=password

	// ParsePostgreConn determines if url contains `postgre://` or `postgresql://`
	// And converts it to driver form `user=xxx password=xxx host=xxx...`
	//func ParsePostgreConn(url string) (string, error) {
	//	if !strings.HasPrefix(url, "postgre://") || !strings.HasPrefix(url, "postgresql://") {
	//		ret, err := pq.ParseURL(url)
	//		if err != nil {
	//			return "", err
	//		}
	//		return ret, nil
	//	}
	//	return url, nil
	//}

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
