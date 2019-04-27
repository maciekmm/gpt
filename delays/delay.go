package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/maciekmm/gpt"
)

const delaysUrl = "http://87.98.237.99:88/delays"

const schema = `CREATE TABLE IF NOT EXISTS delays (
	last_update timestamp NOT NULL,
	"timestamp" timestamp NOT NULL,
	stop_id integer NOT NULL,
	trip integer NOT NULL,
	trip_id integer NOT NULL,
	route_id integer NOT NULL,
	"id" varchar(255) NOT NULL,
	delay_in_seconds integer NOT NULL,
	estimated_time time NOT NULL,
	head_sign varchar(255) NOT NULL,
	"status" varchar(255) NOT NULL,
	theoretical_time time NOT NULL,
	vehicle_code integer NOT NULL,
	vehicle_id integer NOT NULL,
	CONSTRAINT delays_pk PRIMARY KEY(last_update, stop_id, route_id, trip_id, trip)
)`

type Delays struct {
	LastUpdate gpt.UpdateTime `json:"lastUpdate"`
	Delay      []*Delay       `json:"delay"`
}

type Delay struct {
	LastUpdate gpt.UpdateTime            `json:"lastUpdate"`
	Timestamp  gpt.SimpleTimeWithSeconds `json:"timestamp"`
	StopID     int64                     `json:"stopId"`
	Trip       int64                     `json:"trip"`
	TripID     int64                     `json:"tripId"`
	RouteID    int64                     `json:"routeId"`

	ID              string         `json:"id"`
	DelayInSeconds  int            `json:"delayInSeconds"`
	EstimatedTime   gpt.SimpleTime `json:"estimatedTime"`
	HeadSign        string         `json:"headSign"`
	Status          string         `json:"status"`
	TheoreticalTime gpt.SimpleTime `json:"theoreticalTime"`
	VehicleCode     int64          `json:"vehicleCode"`
	VehicleID       int64          `json:"vehicleId"`
}

func getDelays() ([]*Delay, error) {
	var delayMap map[string]*Delays
	resp, err := http.DefaultClient.Get(delaysUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&delayMap); err != nil {
		return nil, err
	}
	flatDelays := make([]*Delay, 0)
	for key, delays := range delayMap {
		stopID, err := strconv.Atoi(key)
		if err != nil {
			log.Println("could not atoi ", key)
			continue
		}
		for _, delay := range delays.Delay {
			delay.LastUpdate = delays.LastUpdate
			delay.StopID = int64(stopID)
			flatDelays = append(flatDelays, delay)
		}
	}
	return flatDelays, nil
}

func main() {
	log.Println("Downloading delays")
	db, err := gpt.InitDatabase()
	if err != nil {
		log.Fatalln("cannot connect to database", err)
	}
	if _, err := db.Exec(schema); err != nil {
		log.Fatalln(err)
	}

	delays, err := getDelays()
	if err != nil {
		log.Fatalln("error getting stop urls", err)
	}
	for _, delay := range delays {
		if _, err := db.NamedExec(`INSERT INTO delays (last_update, "timestamp", stop_id, trip, trip_id, route_id, "id", delay_in_seconds, estimated_time, head_sign, status, theoretical_time, vehicle_code, vehicle_id)
			VALUES (:last_update, :timestamp, :stop_id, :trip, :trip_id, :route_id, :id, :delay_in_seconds, :estimated_time, :head_sign, :status, :theoretical_time, :vehicle_code, :vehicle_id)
			ON CONFLICT (last_update, stop_id, route_id, trip_id, trip) DO NOTHING`, delay); err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("Finished downlaoding delays")
}
