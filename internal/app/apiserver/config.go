package apiserver

import (
	"github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	dbhost  = "DBHOST"
	dbport  = "DBPORT"
	dbuser  = "DBUSER"
	dbpass  = "DBPASS"
	dbname  = "DBNAME"
	dburl   = "DATABASE_URL"
	dumpDir = "DUMPS_DIR"
)

type Config struct {
	BindAddr        string `toml:"bind_addr"`
	LogLevel        string `toml:"log_level"`
	DatabaseURL     string `os:"database_url"`
	PGDatabaseURL   string
	DatabaseDumpDir string
}

func NewConfig() *Config {
	psqlInfo, pgDBConn := dbConfig()
	dbDumpDir := os.Getenv(dumpDir)
	if dbDumpDir == "" {
		dbDumpDir = "dumps"
	}
	cwd, _ := os.Getwd()
	dbDumpDir = filepath.Join(cwd, dbDumpDir)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT ERROR")
	}

	return &Config{
		BindAddr:        ":" + port,
		LogLevel:        "debug",
		DatabaseURL:     psqlInfo,
		PGDatabaseURL:   pgDBConn,
		DatabaseDumpDir: dbDumpDir,
	}
}

func dbConfig() (string, string) {
	PGDatabaseURL := os.Getenv(dburl)
	driverURL, err := parsePostgreConn(PGDatabaseURL)
	if err != nil {
		log.Fatalf("INCORRECT DATABASE URL. ENV: %v, PARSED %v", PGDatabaseURL, driverURL)
	}
	return driverURL, PGDatabaseURL

	//conf := make(map[string]string)
	//host, ok := os.LookupEnv(dbhost)
	//if !ok {
	//	panic("DBHOST environment variable required but not set")
	//}
	//port, ok := os.LookupEnv(dbport)
	//if !ok {
	//	panic("DBPORT environment variable required but not set")
	//}
	//user, ok := os.LookupEnv(dbuser)
	//if !ok {
	//	panic("DBUSER environment variable required but not set")
	//}
	//password, ok := os.LookupEnv(dbpass)
	//if !ok {
	//	panic("DBPASS environment variable required but not set")
	//}
	//name, ok := os.LookupEnv(dbname)
	//if !ok {
	//	panic("DBNAME environment variable required but not set")
	//}
	//conf[dbhost] = host
	//conf[dbport] = port
	//conf[dbuser] = user
	//conf[dbpass] = password
	//conf[dbname] = name
	//return conf
}

func parsePostgreConn(url string) (string, error) {
	if !strings.HasPrefix(url, "postgre://") || !strings.HasPrefix(url, "postgresql://") {
		ret, err := pq.ParseURL(url)
		if err != nil {
			return "", err
		}
		return ret, nil
	}
	return url, nil
}
