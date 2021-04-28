package sqlstore

import (
	"bytes"
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type DumpRepository struct {
	store *Store
}

func (d *DumpRepository) CreateDump(dumpFileDir string) string {
	dumpFileName := "dump.sql" // TODO: RENAME??
	dumpFilePath := filepath.Join(dumpFileDir, dumpFileName)
	if _, err := os.Stat(dumpFileDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dumpFileDir, os.ModePerm); err != nil {
				log.Printf("DUMP DIR CREATE FAILLED: %v", err)
			}
		}
	}

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd := exec.Command("pg_dump", d.store.config.PGDatabaseURL, "--column-inserts", "-f", dumpFilePath)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Printf("DUMPING FAILLED: %v", err)
		}
	}()

	return dumpFilePath
}

func (d *DumpRepository) Execute(dumpFilePath string) error {
	dumpFileContent, err := ioutil.ReadFile(dumpFilePath)
	if err != nil {
		log.Println(err)
		return err
	}
	dumpQueries := string(dumpFileContent)

	transaction, err := d.store.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	defer func(transaction *sql.Tx) {
		err := transaction.Rollback()
		log.Println(err)
	}(transaction)
	if _, err := transaction.Exec(
		`DROP SCHEMA IF EXISTS public CASCADE;
			   CREATE SCHEMA IF NOT EXISTS public;`,
	); err != nil {
		log.Println(err)
		return err
	}
	if _, err := transaction.Exec(dumpQueries); err != nil {
		log.Println(err)
		return err
	}
	if err := transaction.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
