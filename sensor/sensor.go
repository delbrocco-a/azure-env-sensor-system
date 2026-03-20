
package main

import (
	"math/rand"
	"time"
)

// # sensor.go
// Input  : n := number of sensors that you want to simulate
// Output : n rows of sensor data in an array

/*
Creates as shown sensor data to use as a test for implementation
*/

// ## Sensor Structure ========================================================
// ### Generates data per coursework specification

type SensorData struct {
	SensorID         int    `json:"sensor_id"`
	Temperature      int    `json:"temperature"`
	WindSpeed        int    `json:"wind_speed"`
	RelativeHumidity int    `json:"relative_humidity"`
	CO2Level         int    `json:"co2_level"`
	Timestamp        string `json:"timestamp"`
}

// ## Data Generation =========================================================
// ### Random number generators and SensorData struct random initialisers

func randomIntInRange(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func generateSensorData(n int) []SensorData {
	sensors := make([]SensorData, n)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	for i := 0; i < n; i++ { // ### Generate n amounts of random sensorData 
		sensors[i] = SensorData{
			SensorID:         i + 1,
			Temperature:      randomIntInRange(TEMP_MIN, TEMP_MAX),
			WindSpeed:        randomIntInRange(WIND_MIN, WIND_MAX),
			RelativeHumidity: randomIntInRange(HUMI_MIN, HUMI_MAX),
			CO2Level:         randomIntInRange(HUMI_MIN, HUMI_MAX),
			Timestamp:        timestamp,
		}
	}
	return sensors
}