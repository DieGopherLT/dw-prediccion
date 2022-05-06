package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("could not load environment variables:", err.Error())
	}

	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;",
		os.Getenv("DB_SERVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("PORT"))

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatalln("could not connect to database", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalln("ping failed:", err.Error())
	}
}
