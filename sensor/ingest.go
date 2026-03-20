
package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Ingestion Function (ingest.go)
// Input  : Sensor data, and a live db pointer
// Output : Error code depending on the outcome of SQL insert

/*
As it says on the tin: given some data and a db pointer, the function executes
an SQL insert of the data provided, creating it's own context and failsafes.
Note that it also enqueues a message for the statistics function.
*/

// ## SQL Query ===============================================================
// ### In this file, as they are specific to the implementation

const SQL_INSERT_READINGS string = `INSERT INTO SensorReadings
	         (sensor_id, temp, windspeed, relative_humidity, co2, timestamp)
	         VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`

const TIMEOUT = 30 * time.Second

// ## Data Insertion/Ingestion ================================================
// ### Note that ingestions add to a message queue

func ingestReadings(data []SensorData, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	for _, sensor := range data {
		_, err := db.ExecContext(
			ctx, SQL_INSERT_READINGS,
			sensor.SensorID,
			sensor.Temperature,
			sensor.WindSpeed,
			sensor.RelativeHumidity,
			sensor.CO2Level,
			sensor.Timestamp,
		)
		if err != nil {
			return fmt.Errorf("error inserting sensor %d: %v", sensor.SensorID, err)
		}
	}

	// ### Push to the message queue that new data been ingested (for cswk)
	err := enqueueMessage("New rows appended to %s on %s", DATABASE, SERVER)
	if err != nil {
		return fmt.Errorf("failed to enqueue queue message: %v", err)
	}

	return nil
}