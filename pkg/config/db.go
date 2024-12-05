package config

import (
	"database/sql"
	"fmt"
	"log"
	"practice-go/pkg/constants"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", constants.DB_USERNAME, constants.DB_PASSWORD, constants.DB_NAME)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
		db.Close()
	}

	return db
}
