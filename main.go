package main

import (
	"log"

	"github.com/VolticFroogo/Launcher-Server/db"
	"github.com/VolticFroogo/Launcher-Server/handle"
)

func main() {
	// Initialise the database.
	err := db.Init()
	if err != nil {
		log.Printf("Error initialising database: %v", err)
		return
	}

	// Start listening for incoming connections.
	handle.Listen()
}
