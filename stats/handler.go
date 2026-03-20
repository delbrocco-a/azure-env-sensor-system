package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// # handler.go
// Input  : HTTP request from Azure queue trigger
// Output : JSON response with statistics

/*
The same as the sensor implementation, only responds to a HTTP trigger from
a message queue update. It also sends stats to listening client.
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

func executeStatisticsCalculation() ([]SensorStats, error) {
	db, err := getDBConnection(SERVER, DATABASE)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer db.Close()

	stats, err := calcStats(db)
	if err != nil {
		return nil, fmt.Errorf("calculation failed: %w", err)
	}

	return stats, nil
}

// ## HTTP Handler ============================================================

func statisticsHandler(w http.ResponseWriter, r *http.Request) {
	var invokeReq InvokeRequest

	// ### HTTP requests are the same, now coming from a message queue
	if err := json.NewDecoder(r.Body).Decode(&invokeReq); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ### Note queue trigger
	logs := []string{}
	logs = append(logs, fmt.Sprintf("Queue trigger fired at %s",
		time.Now().UTC().Format(time.RFC3339),
	))

	stats, err := executeStatisticsCalculation()

	if err != nil {
		log.Printf("Failed to calculate statistics: %v", err)
		logs = append(logs, fmt.Sprintf("Error: %v", err))

		invokeResp := InvokeResponse{
			Outputs:     map[string]interface{}{},
			Logs:        logs,
			ReturnValue: nil,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(invokeResp)
		return
	}

	logs = append(logs, fmt.Sprintf(
		"Successfully calculated statistics for %d sensors", len(stats),
	))

	// ### Now put the data in your response to the server updating the db
	for _, s := range stats {
		logs = append(
			logs, fmt.Sprintf(
				"Sensor %d - Temp: %.2f-%.2f (avg: %.2f), " +
				"Wind: %.2f-%.2f (avg: %.2f), " +
				"Humidity: %.2f-%.2f (avg: %.2f), " +
				"CO2: %d-%d (avg: %.2f)",
				s.SensorID,
				s.MinTemp, s.MaxTemp, s.AvgTemp,
				s.MinWind, s.MaxWind, s.AvgWind,
				s.MinHumidity, s.MaxHumidity, s.AvgHumidity,
				s.MinCO2, s.MaxCO2, s.AvgCO2,
			)
		)
	}

	invokeResp := InvokeResponse{
		Outputs: map[string]interface{}{
			"stats": stats,
		},
		Logs: logs,
		ReturnValue: fmt.Sprintf(
			"Processed statistics for %d sensors", len(stats)),
	}

	// ### For meaningful testing/debug
	for _, entry := range logs {
		log.Print(entry)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invokeResp)
}