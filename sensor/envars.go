const SERVER = "dsysts.database.windows.net"
const DATABASE = "sensors"

const QUEUE string = "statisticqueue"
const ACCOUNT_NAME_VAR string = "AZURE_STORAGE_ACCOUNT_NAME"
const ACCOUNT_KEY_VAR string = "AZURE_STORAGE_ACCOUNT_KEY"

/* IN ingest.go, modification of these changes implementation
const SQL_INSERT_READINGS string = `INSERT INTO SensorReadings
	         (sensor_id, temp, windspeed, relative_humidity, co2, timestamp)
	         VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`

const TIMEOUT = 30 * time.Second */

// ## Sensor stuff for data simulation

const SENSOR_COUNT int = 20

// ## Random Number Generation Ranges

const TEMP_MIN int = 5
const TEMP_MAX int = 18

const WIND_MIN int = 12
const WIND_MAX int = 24

const HUMI_MIN int = 30
const HUMI_MAX int = 60

const C02_MIN  int = 400
const CO2_MAX  int = 1600