package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	dbpool, err := pgxpool.New(context.Background(), "postgres://go:password@127.0.0.1:5432/godb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var title string
	var price float64
	err = dbpool.QueryRow(context.Background(), "select title, price from album where id=$1", 1).Scan(&title, &price)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(title, price)

	albums, err := albumsByArtist("Zakk", dbpool)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	alb, err := albumByID(1, dbpool)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	id, err := addAlbum(Album{
		Title:  "Testing",
		Artist: "Zakk",
		Price:  25.44,
	}, dbpool)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album added: %d\n", id)
}

func albumByID(id int64, dbpool *pgxpool.Pool) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := dbpool.QueryRow(context.Background(), "SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

func albumsByArtist(name string, dbpool *pgxpool.Pool) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := dbpool.Query(context.Background(), "SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	var alb Album
	_, err = pgx.ForEachRow(rows, []any{&alb.ID, &alb.Title, &alb.Artist, &alb.Price}, func() error {
		albums = append(albums, alb)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func addAlbum(alb Album, dbpool *pgxpool.Pool) (int64, error) {
	row := dbpool.QueryRow(context.Background(), "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING ID", alb.Title, alb.Artist, alb.Price)
	var id int64
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("addAlbum %d: no such album", id)
		}
		return 0, fmt.Errorf("addAlbum %d: %v", id, err)
	}
	return id, nil
}
