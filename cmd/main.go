package main

import (
	"log"
	"net/http"

	"internalTransferSystem/db"
	"internalTransferSystem/handlers"

	_ "github.com/lib/pq"
)

func main() {
	DB, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer DB.Close()

	handler := handlers.NewHandler(DB)
	handler.RegisterRoutes()

	log.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
