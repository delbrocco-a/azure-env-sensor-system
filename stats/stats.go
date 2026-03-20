package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// # stats.go
// Input  : database connection
// Output : calculated statistics per sensor

/*
This is a statistical analysis function for all of the data in the database.
It creates it's own context and hardcoded select query, as well as some basic
error checking. It's important to note, the actual data analysis itself is done
via the SQL query, not the go function. The function then returns this data
is the format presented in the struct bellow. All are in float format for ease 
of use, but the databse only contains integer values.
*/

// ## Statistics Structure ====================================================

type SensorStats struct {
	SensorID    int     `json:"sensor_id"`
	MinTemp     float64 `json:"min_temperature"`
	MaxTemp     float64 `json:"max_temperature"`
	AvgTemp     float64 `json:"avg_temperature"`
	MinWind     float64 `json:"min_wind_speed"`
	MaxWind     float64 `json:"max_wind_speed"`
	AvgWind     float64 `json:"avg_wind_speed"`
	MinHumidity float64 `json:"min_humidity"`
	MaxHumidity float64 `json:"max_humidity"`
	AvgHumidity float64 `json:"avg_humidity"`
	MinCO2      float64 `json:"min_co2"`
	MaxCO2      float64 `json:"max_co2"`
	AvgCO2      float64 `json:"avg_co2"`
}

// ## Query Constants =========================================================

const STATS_QUERY = `
	SELECT 
		sensor_id,
		MIN(temp) as min_temp,
		MAX(temp) as max_temp,
		AVG(CAST(temp AS FLOAT)) as avg_temp,
		MIN(windspeed) as min_wind,
		MAX(windspeed) as max_wind,
		AVG(CAST(windspeed AS FLOAT)) as avg_wind,
		MIN(relative_humidity) as min_humidity,
		MAX(relative_humidity) as max_humidity,
		AVG(CAST(relative_humidity AS FLOAT)) as avg_humidity,
		MIN(co2) as min_co2,
		MAX(co2) as max_co2,
		AVG(CAST(co2 AS FLOAT)) as avg_co2
	FROM SensorReadings
	GROUP BY sensor_id
	ORDER BY sensor_id
`

const TIMEOUT = 30 * time.Second

// ## Statistics Calculation ==================================================

func calcStats(db *sql.DB) ([]SensorStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	// ### The main functionality of this function is here
	rows, err := db.QueryContext(ctx, STATS_QUERY)
	if err != nil {
		return nil, fmt.Errorf("error querying statistics: %v", err)
	}
	defer rows.Close()

	// ### The rest of this is data formatting and error checking
	var stats []SensorStats
	for rows.Next() {
		var s SensorStats
		err := rows.Scan(
			&s.SensorID,
			&s.MinTemp, &s.MaxTemp, &s.AvgTemp,
			&s.MinWind, &s.MaxWind, &s.AvgWind,
			&s.MinHumidity, &s.MaxHumidity, &s.AvgHumidity,
			&s.MinCO2, &s.MaxCO2, &s.AvgCO2,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		stats = append(stats, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return stats, nil
}