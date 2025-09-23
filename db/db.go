package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("postgres", ConnStr)

	if err != nil {
		return fmt.Errorf("database connection error %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("checking connection error %v", err)
	}

	fmt.Println("database connection successful")
	return nil
}
