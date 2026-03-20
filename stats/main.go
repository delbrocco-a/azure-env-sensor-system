package main

import (
	"log"
	"net/http"
	"os"
)

// Main.go

/* Sets up the port and service for the stats function. Note that function and
host.json are now considering http requests from the message queue  */

// # Main /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\

func main() {
	http.HandleFunc("/stats", statisticsHandler)
	
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		customHandlerPort = "8080"
	}

	log.Printf("Starting Go custom handler on port %s", customHandlerPort)
	log.Printf("Queue trigger function ready - waiting for messages...")

	if err := http.ListenAndServe(":"+customHandlerPort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// --- /\/\ /\/\ /\/\ /\/\ /\/\ /\/\/\/\ /\/\ /\/\ /\/\ /\/\ /\/\/\/\ /\/\ /\/\
