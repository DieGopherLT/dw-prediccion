package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

type Store struct {
	Name string
}

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	query := `
		SELECT s.store_name
			FROM
				BikeStores.sales.order_items oi
		INNER JOIN
			BikeStores.sales.orders o ON o.order_id =  oi.product_id
		INNER JOIN
			BikeStores.sales.stores s ON o.store_id = s.store_id;
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatalln("error at querying database:", err.Error())
	}

	stores := []Store{}
	for rows.Next() {
		store := Store{}
		err := rows.Scan(&store.Name)
		if err != nil {
			log.Fatalln("error while scanning rows", err.Error())
		}
		stores = append(stores, store)
	}

	fmt.Println("number of rows", len(stores))
}
