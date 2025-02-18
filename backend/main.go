package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var db *sql.DB

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:password@127.0.0.1:5432/recordings")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var title string
	var price float64
	err = conn.QueryRow(context.Background(), "select title, price from album where id=$1", 1).Scan(&title, &price)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(title, price)
}
