package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

/* ### This is deceptively called main, when it only infact, sets up the server
for http service to the sensor function/handler*/

// # Main /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\

func main() {

	log.Println("Starting custom handler HTTP server")
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		customHandlerPort = "8080"
	}

	log.Printf("FUNCTIONS_CUSTOMHANDLER_PORT environment variable exists: %v", exists)
	log.Printf("Port value: %s", customHandlerPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/sensor", sensorHandler)

	addr := "0.0.0.0:" + customHandlerPort
	log.Printf("Attempting to listen on %s", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  TIMEOUT,
		WriteTimeout: TIMEOUT,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

//   /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\ /\/\
