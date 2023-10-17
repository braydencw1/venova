package db

import (
	"fmt"
	"database/sql"
//	"log"
	"os"
	_ "github.com/lib/pq"
)


func Dquery(query string) (*sql.DB, *sql.Rows, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_DB")

	con := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", con)
	if err != nil {
		return nil, nil, err
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Connected")
	return db, rows, nil
}
