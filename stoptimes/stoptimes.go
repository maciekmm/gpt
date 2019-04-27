package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/maciekmm/gpt"
	"gopkg.in/guregu/null.v3"
)

const tripStopTimesUrl = "https://ckan.multimediagdansk.pl/dataset/c24aa637-3619-4dc2-a171-a23eec8f2172/resource/a023ceb0-8085-45f6-8261-02e6fcba7971/download/stoptimes.json"
const schema = `CREATE TABLE IF NOT EXISTS stop_times (
	last_update timestamp NOT NULL,
	stop_id integer NOT NULL,
	route_id integer NOT NULL,
	trip_id integer NOT NULL,
	arrival_time timestamp NOT NULL,
	departure_time timestamp NOT NULL,
	date timestamp NOT NULL,
	agency_id integer NOT NULL,
	topology_version_id integer NOT NULL,
	stop_sequence integer NOT NULL,
	variant_id integer,
	note_symbol varchar(255),
	note_description text,
	bus_service_name varchar(255),
	"order" integer NOT NULL,
	nonpassenger bit,
	ticket_zone_border bit,
	on_demand bit,
	virtual bit,
	islupek int,
	wheelchair_accessible varchar(255),
	stop_short_name varchar(255) not null,
	CONSTRAINT stop_times_pk PRIMARY KEY(last_update, date, stop_id, route_id, trip_id)
)`

type StopTimes struct {
	LastUpdate gpt.UpdateTime `json:"lastUpdate"`
	StopTimes  []*StopTime    `json:"stopTimes"`
}

type StopTime struct {
	LastUpdate gpt.UpdateTime `json:"lastUpdate"`
	Date       gpt.Date       `json:"date"`
	StopID     int            `json:"stopId"`
	RouteID    int            `json:"routeId"`
	TripID     int            `json:"tripId"`

	ArrivalTime          gpt.ArrivalTime `json:"arrivalTime"`
	DepartureTime        gpt.ArrivalTime `json:"departureTime"`
	AgencyID             int             `json:"agencyId"`
	TopologyVersionID    int             `json:"topologyVersionId"`
	StopSequence         int             `json:"stopSequence"`
	VariantID            null.Int        `json:"variantId"`
	NoteSymbol           null.String     `json:"noteSymbol"`
	NoteDescription      null.String     `json:"noteDescription"`
	BusServiceName       string          `json:"busServiceName"`
	Order                int             `json:"order"`
	Nonpassenger         null.Int        `json:"nonpassenger"`
	TicketZoneBorder     null.Int        `json:"ticketZoneBorder"`
	OnDemand             null.Int        `json:"onDemand"`
	Virtual              null.Int        `json:"virtual"`
	Islupek              null.Int        `json:"islupek"`
	WheelchairAccessible null.Int        `json:"wheelchairAccessible"`
	StopShortName        string          `json:"stopShortName"`
}

func getStopUrls() (map[string][]string, error) {
	var urls map[string][]string
	resp, err := http.DefaultClient.Get(tripStopTimesUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&urls); err != nil {
		return nil, err
	}
	return urls, nil
}

func getStopTimes(url string) (*[]*StopTime, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	stopTimes := StopTimes{}
	if err := decoder.Decode(&stopTimes); err != nil {
		return nil, err
	}
	fmt.Println(url, stopTimes.LastUpdate)

	for _, stopTime := range stopTimes.StopTimes {
		stopTime.LastUpdate = stopTimes.LastUpdate
		// Fix their obnoxious format
		arrivalDay := stopTime.Date.Day()
		if stopTime.ArrivalTime.Day() == 31 {
			arrivalDay = stopTime.Date.Day() + 1
		}
		stopTime.ArrivalTime.Time = time.Date(stopTime.Date.Year(), stopTime.Date.Month(), arrivalDay,
			stopTime.ArrivalTime.Hour(), stopTime.ArrivalTime.Minute(), stopTime.ArrivalTime.Second(), stopTime.ArrivalTime.Nanosecond(), stopTime.ArrivalTime.Location())
		departureDay := stopTime.Date.Day()
		if stopTime.DepartureTime.Day() == 31 {
			departureDay = stopTime.Date.Day() + 1
		}
		stopTime.DepartureTime.Time = time.Date(stopTime.Date.Year(), stopTime.Date.Month(), departureDay,
			stopTime.DepartureTime.Hour(), stopTime.DepartureTime.Minute(), stopTime.DepartureTime.Second(), stopTime.DepartureTime.Nanosecond(), stopTime.DepartureTime.Location())

	}

	return &stopTimes.StopTimes, nil
}

func main() {
	db, err := gpt.InitDatabase()
	if err != nil {
		log.Fatalln("cannot connect to database", err)
	}
	defer db.Close()
	if _, err := db.Exec(schema); err != nil {
		log.Fatalln(err)
	}

	tripUrls, err := getStopUrls()
	if err != nil {
		log.Fatalln("error getting stop urls", err)
	}
	for key, urls := range tripUrls {
		if len(urls) < 1 {
			continue
		}
		stopTimes, err := getStopTimes(urls[0])
		if err != nil {
			log.Println("error fetching ", key, " trip", err)
		}
		for _, stopTime := range *stopTimes {
			if _, err := db.NamedExec(`INSERT INTO stop_times (last_update, stop_id, route_id, trip_id, arrival_time, departure_time, date, agency_id, topology_version_id, stop_sequence, variant_id, note_symbol, note_description, bus_service_name, "order", nonpassenger, ticket_zone_border, on_demand, virtual, islupek, wheelchair_accessible, stop_short_name)
			VALUES (:last_update, :stop_id, :route_id, :trip_id, :arrival_time, :departure_time, :date, :agency_id, :topology_version_id, :stop_sequence, :variant_id, :note_symbol, :note_description, :bus_service_name, :order, :nonpassenger, :ticket_zone_border, :on_demand, :virtual, :islupek, :wheelchair_accessible, :stop_short_name)
			ON CONFLICT (last_update, date, stop_id, route_id, trip_id) DO NOTHING`, stopTime); err != nil {
				log.Fatalln(err)
			}

		}

	}
}
