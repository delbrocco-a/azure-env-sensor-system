package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// # handler.go
// Input  : HTTP request from Azure Function
// Output : JSON response with execution results

/*
This has two principle parts. The first part, is the core of the function, and
creates the database connection, generates the values, and writes them to the
database, doing it's own error checking and failsafes. The next part ie the
"handler" acts a sort of HTTP wrapper for this function, handling the actual 
I/O HTTP response for the function. It recieves a request, and then it tries to
write the data, and returns information on how it got on to the requester.

It is key to note that this is an executable function for main to set up,
and for azure to find in the function.json
*/

// ## JSON Structures =========================================================

type InvokeRequest struct {
	Data     map[string]interface{} `json:"Data"`
	Metadata map[string]interface{} `json:"Metadata"`
}

type InvokeResponse struct {
	Outputs     map[string]interface{} `json:"Outputs"`
	Logs        []string               `json:"Logs"`
	ReturnValue interface{}            `json:"ReturnValue"`
}

// ## Core Function ===========================================================

func executeDataIngestion() error {
	db, err := getDBConnection(SERVER, DATABASE)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer db.Close()

	sensorData := generateSensorData(SENSOR_COUNT)
	err = ingestReadings(sensorData, db)
	if err != nil {
		return fmt.Errorf("storage failed: %w", err)
	}

	return nil
}

// ## HTTP Handler ============================================================

func sensorHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	var invokeReq InvokeRequest

	// ### Recieve request
	if r.Method == "POST" && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&invokeReq); err != nil {
			log.Printf("Error decoding request: %v", err)
		}
	}

	// ### Execute the function (sensors)
	logs := []string{}
	err := executeDataIngestion()

	execTime := time.Since(startTime).Milliseconds()

	var message string
	var statusCode int
	if err != nil {
		message = fmt.Sprintf("Error: %v", err)
		statusCode = 500
		logs = append(logs, message)
	} else {
		message = "Data inserted successfully!"
		statusCode = 200
		logs = append(logs, fmt.Sprintf(
			"Inserted %d sensor readings", SENSOR_COUNT))
	}

	// ## Format and send responce to requester
	response := map[string]interface{}{
		"statusCode":        statusCode,
		"body":              message,
		"execution_time_ms": execTime,
		"sensor_count":      SENSOR_COUNT,
		"timestamp":         time.Now().UTC().Format(time.RFC3339),
	}

	if r.Header.Get("X-Azure-Functions-InvocationId") != "" {
		invokeResp := InvokeResponse{
			Outputs: map[string]interface{}{
				"res": response,
			},
			Logs:        logs,
			ReturnValue: nil,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invokeResp)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}
}