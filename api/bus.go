package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const API_URI = "https://api.opendata.metlink.org.nz/v1/stop-predictions?stop_id="
const API_KEY = "iuoMNXQjzC1PjijgMjKkHhYWPb4ZES2UpaYfgsd0"

func Handler(w http.ResponseWriter, r *http.Request) {
	stop := r.URL.Query().Get("stop")
	services := r.URL.Query().Get("services")
	predictions, err := predict(stop, services)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	jsonstr, err := json.Marshal(predictions)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, string(jsonstr))
}

func contain(str, key string) bool {
	keys := strings.Split(str, ",")
	for _, v := range keys {
		if strings.TrimSpace(v) == strings.TrimSpace(key) {
			return true
		}
	}
	return false
}

func predict(stop, services string) ([]Departure, error) {
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	req, err := http.NewRequest(http.MethodGet, API_URI+stop, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", API_KEY)

	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	prediction := StopPrediction{}
	jsonErr := json.Unmarshal(body, &prediction)
	if jsonErr != nil {
		return nil, jsonErr
	}

	departures := prediction.Departures
	filtered := []Departure{}
	for _, v := range departures {
		if contain(services, v.ServiceID) {
			filtered = append(filtered, v)
		}
	}
	return filtered, nil
}

type Departure struct {
	StopID    string `json:"stop_id"`
	ServiceID string `json:"service_id"`
	Direction string `json:"direction"`
	Operator  string `json:"operator"`
	Origin    struct {
		StopID string `json:"stop_id"`
		Name   string `json:"name"`
	} `json:"origin"`
	Destination struct {
		StopID string `json:"stop_id"`
		Name   string `json:"name"`
	} `json:"destination"`
	Delay     string `json:"delay"`
	VehicleID string `json:"vehicle_id"`
	Name      string `json:"name"`
	Arrival   struct {
		Aimed    time.Time `json:"aimed"`
		Expected time.Time `json:"expected"`
	} `json:"arrival"`
	Departure struct {
		Aimed    time.Time `json:"aimed"`
		Expected time.Time `json:"expected"`
	} `json:"departure"`
	Status               string `json:"status"`
	Monitored            bool   `json:"monitored"`
	WheelchairAccessible bool   `json:"wheelchair_accessible"`
	TripID               string `json:"trip_id"`
}

type StopPrediction struct {
	Farezone   string      `json:"farezone"`
	Closed     bool        `json:"closed"`
	Departures []Departure `json:"departures"`
}
