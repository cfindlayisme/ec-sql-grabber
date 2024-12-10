package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cfindlayisme/ec-sql-grabber/dataretriever"
	"github.com/cfindlayisme/ec-sql-grabber/models"
	_ "github.com/lib/pq"
)

const (
	dbConn = "postgres://user:password@localhost:5432/weatherdb?sslmode=disable"
)

var provinceTimeZones = map[string]string{
	"NB": "America/Moncton",
	"NS": "America/Halifax",
	"PE": "America/Halifax",
	"NL": "America/St_Johns",
	"QC": "America/Toronto",
	"ON": "America/Toronto",
	"MB": "America/Winnipeg",
	"SK": "America/Regina",
	"AB": "America/Edmonton",
	"BC": "America/Vancouver",
	"YT": "America/Whitehorse",
	"NT": "America/Yellowknife",
	"NU": "America/Iqaluit",
}

func InsertObservationData(db *sql.DB, observations []models.Observation, province string) error {
	query := `
		INSERT INTO weather_observations (
			station_name, latitude, longitude, timestamp, temperature, dew_point, 
			relative_humidity, wind_speed, wind_direction, wind_gust_speed, wind_chill, 
			mean_sea_level, tendency_amount, tendency_characteristic, present_weather, 
			horizontal_visibility, total_cloud_cover, humidex, observation_date_utc, observation_date_local_time, geom
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, ST_SetSRID(ST_MakePoint($3, $2), 4326))
	`

	for _, obs := range observations {
		latLon := strings.Split(obs.FeatureOfInterest.FeatureCollection.Location.Point.Pos, " ")
		if len(latLon) != 2 {
			continue
		}
		latitude := latLon[0]
		longitude := latLon[1]

		var stationName string
		for _, element := range obs.Metadata.IdentificationElements.Element {
			if element.Name == "station_name" {
				stationName = element.Value
				break
			}
		}

		var observationDateUTC, observationDateLocalTime string
		for _, element := range obs.Metadata.IdentificationElements.Element {
			if element.Name == "observation_date_utc" {
				observationDateUTC = element.Value
			} else if element.Name == "observation_date_local_time" {
				observationDateLocalTime = element.Value
			}
		}

		utcTime, err := time.Parse(time.RFC3339, observationDateUTC)
		if err != nil {
			return fmt.Errorf("failed to parse observation_date_utc: %w", err)
		}

		timeZone, exists := provinceTimeZones[province]
		if !exists {
			return fmt.Errorf("unknown time zone for province: %s", province)
		}

		loc, err := time.LoadLocation(timeZone)
		if err != nil {
			return fmt.Errorf("failed to load time zone: %w", err)
		}

		localTime, err := time.ParseInLocation("2006-01-02T15:04:05.000 MST", observationDateLocalTime, loc)
		if err != nil {
			return fmt.Errorf("failed to parse observation_date_local_time: %w", err)
		}

		var temperature, dewPoint, relativeHumidity, windSpeed, windDirection, windGustSpeed, windChill, meanSeaLevel, tendencyAmount, tendencyCharacteristic, presentWeather, horizontalVisibility, totalCloudCover, humidex string
		for _, element := range obs.Result.Elements.Element {
			switch element.Name {
			case "air_temperature":
				temperature = element.Value
			case "dew_point":
				dewPoint = element.Value
			case "relative_humidity":
				relativeHumidity = element.Value
			case "wind_speed":
				windSpeed = element.Value
			case "wind_direction":
				windDirection = element.Value
			case "wind_gust_speed":
				windGustSpeed = element.Value
			case "wind_chill":
				windChill = element.Value
			case "mean_sea_level":
				meanSeaLevel = element.Value
			case "tendency_amount":
				tendencyAmount = element.Value
			case "tendency_characteristic":
				tendencyCharacteristic = element.Value
			case "present_weather":
				presentWeather = element.Value
			case "horizontal_visibility":
				horizontalVisibility = element.Value
			case "total_cloud_cover":
				totalCloudCover = element.Value
			case "humidex":
				humidex = element.Value
			}
		}

		_, err = db.Exec(query, stationName, latitude, longitude,
			obs.SamplingTime.TimeInstant.TimePosition, temperature, dewPoint, relativeHumidity,
			windSpeed, windDirection, windGustSpeed, windChill, meanSeaLevel, tendencyAmount,
			tendencyCharacteristic, presentWeather, horizontalVisibility, totalCloudCover, humidex,
			utcTime, localTime)

		log.Println("Log for "+stationName, latitude, longitude, "at time "+utcTime.Format(time.RFC3339)+" inserted into database")
		if err != nil {
			return fmt.Errorf("failed to insert observation: %w", err)
		}
	}
	return nil
}

func main() {
	// Connect to the database
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	provinces := []string{"NB", "NS", "PE", "NL", "QC", "ON", "MB", "SK", "AB", "BC", "YT", "NT", "NU"}
	for _, province := range provinces {
		xmlData, err := dataretriever.FetchLatestObservationsXMLByProvince(province)
		if err != nil {
			log.Printf("Error fetching XML for province %s: %v", province, err)
			continue
		}
		defer xmlData.Close()

		var weatherData models.WeatherData
		if err := xml.NewDecoder(xmlData).Decode(&weatherData); err != nil {
			log.Printf("Error parsing XML for province %s: %v", province, err)
			continue
		}

		if err := InsertObservationData(db, weatherData.Observation, province); err != nil {
			log.Printf("Error inserting observation data for province %s: %v", province, err)
			continue
		}
		log.Printf("Weather data for province %s successfully inserted.", province)
	}
}
