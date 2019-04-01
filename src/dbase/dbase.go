package dbase

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "docker"
	password = "docker"
	dbname   = "postgres"

	schema = "./src/dbase/dump.sql"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return DB, nil
}

func DumpDB(db *sql.DB) error {
	buffer, err := ioutil.ReadFile(schema)
	if err != nil {
		return err
	}

	schema := string(buffer)
	_, err = db.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}
