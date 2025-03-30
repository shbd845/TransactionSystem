package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	POSTGRES_CONN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	DB, err := sql.Open("postgres", POSTGRES_CONN)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		DB.Close()
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	fmt.Println("Connected to the database successfully!")
	return DB, nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed.")
	}
}
