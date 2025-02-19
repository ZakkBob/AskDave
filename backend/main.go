package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func main() {
	var err error
	dbpool, err = pgxpool.New(context.Background(), "postgres://postgres:password@127.0.0.1:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	t, err := nextTasks(2)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(t, "", "  ")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to serialise tasks: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}
