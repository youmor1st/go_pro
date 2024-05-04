package models

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

var db *pgx.Conn

func ConnectDB() {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:45238@localhost:5432/shop")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	db = conn
	log.Println("Connected to database")
}

func CloseDB() {
	db.Close(context.Background())
	log.Println("Closed database connection")
}
