package main

import (
	"log"
	"os"

	"Sprint-13-14/pkg/db"
	"Sprint-13-14/pkg/server"
)

func main() {
	logger := log.New(os.Stdout, "Serv: ", log.Ldate|log.Ltime|log.Lshortfile)

	if err := db.Init("scheduler.db"); err != nil {
		logger.Fatal("Database initialization error:", err)
	}

	defer db.Close()

	serv := server.NewServer(logger)
	serv.Start()
}
