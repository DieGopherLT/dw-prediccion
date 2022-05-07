package main

import (
	"context"
	"database/sql"
	"errors"
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

	connectionUri := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;",
					os.Getenv("DB_SERVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("PORT"))

	db, err := ConnectToDatabase(connectionUri)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	stores, err := QueryStores(db)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("\t¿En qué tienda será la sig. venta con base en las ordenes pasadas?")
	fmt.Println(">>Ejecutando algoritmo ZeroR...")

	rule := GenerateZeroRRule(stores)

	fmt.Println("\tLa tienda a predecir será", rule)
	fmt.Println(">>Calculando precisión de la predicción...")

	accuracy := CalculatePredictionAccuracy(stores, rule)

	fmt.Printf("\tLa precisión de las predicciones fue de un %.2f por ciento.\n", accuracy * 100)
}

func ConnectToDatabase(connectionUri string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", connectionUri)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid uri: %s", err.Error()))
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ping failed: %s", err.Error()))
	}

	return db, nil
}

func QueryStores(db *sql.DB) ([]Store, error) {
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
		return nil, errors.New(fmt.Sprintf("error at querying database: %s", err.Error()))
	}

	stores := []Store{}
	for rows.Next() {
		store := Store{}
		err := rows.Scan(&store.Name)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error while scanning rows: %s", err.Error()))
		}
		stores = append(stores, store)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("error at scanning: %s", err.Error()))
	}

	return stores, nil
}

func GenerateZeroRRule(stores []Store) string {
	storeRepetitions := make(map[string]int)
	trainingSetLength := len(stores) / 3

	for _, store := range stores[:trainingSetLength - 1] {
		storeRepetitions[store.Name]++
	}

	max, selectedStore := 0, ""
	for store, orders := range storeRepetitions {
		if orders > max {
			selectedStore = store
		}
	}

	return selectedStore
}

func CalculatePredictionAccuracy(stores []Store, selectedStore string) float32 {
	trainingSetLength := len(stores) / 3
	testSetLength := len(stores) - trainingSetLength
	asserts := 0
	
	for _, store := range stores[trainingSetLength - 1:] {
		if store.Name == selectedStore {
			asserts++
		}
	}

	accuracy := float32(asserts) / float32(testSetLength)
	return accuracy
}