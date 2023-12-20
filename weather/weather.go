// Copyright (c) 2016 Joonas Kuorilehto
// License: The MIT License (MIT)

// Package weather retrieves local weather from weather.willab.fi.
// This is an unofficial command line client.
package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// DefaultURL is the URL from where to retrieve weather data.
const DefaultURL = "https://weather.willab.fi/weather.json"

// FromURL retrieves weather data from provided HTTP URL to JSON data.
func FromURL(url string) (w Weather, err error) {
	w.WeatherURL = url
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("HTTP request: %s", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatalf("reading response: %s", err)
	}

	err = json.Unmarshal(body, &w)
	if err != nil {
		log.Fatalf("parsing JSON response: %s", err)
	}
	return
}

// Weather response data structure.
type Weather struct {
	WeatherURL string // "https://weather.willab.fi/weather.json",

	TempNow         Temperature   // -17.1,
	TempHi          Temperature   // -12.2,
	TempLo          Temperature   // -19.7,
	DewPoint        Temperature   // -19.9,
	Humidity        Humidity      // 79,
	AirPressure     Pressure      // 1012.6,
	WindSpeed       Windspeed     // 2.6,
	WindSpeedMax    Windspeed     // 4.9,
	WindDir         Winddirection // 63,
	Precipitation1d Precipitation // 0,
	Precipitation1h Precipitation // 0,
	SolarRad        int           // 10,
	WindChill       Temperature   // -23.5,
	Timestamp       Timestamp     // "2016-01-11 18:21 EET"
}

// Temperature value in degrees Celsius.
type Temperature float64

// Humidity value in %.
type Humidity float64

// Windspeed value in m/s.
type Windspeed float64

// Winddirection value in degrees.
type Winddirection int

// Pressure value in hPA.
type Pressure float64

// Precipitation value in millimeters of rain.
type Precipitation float64

// Timestamp value of the measurement.
type Timestamp time.Time

func (t Temperature) String() string   { return fmt.Sprintf("%.1f °C", t) }
func (h Humidity) String() string      { return fmt.Sprintf("%.0f %%", h) }
func (w Windspeed) String() string     { return fmt.Sprintf("%.0f m/s", w) }
func (w Winddirection) String() string { return fmt.Sprintf("%d°", w) }
func (p Pressure) String() string      { return fmt.Sprintf("%.1f hPa", p) }
func (p Precipitation) String() string { return fmt.Sprintf("%.1f mm", p) }

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("Europe/Helsinki")
	if err != nil {
		panic(fmt.Sprintf("Failed to load time location: %s", err))
	}
}

const layout = "2006-01-02 15:04:05 MST"

// UnmarshalJSON is used to decode JSON into Time value.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	ts, err := time.ParseInLocation(layout, s, loc)
	*t = Timestamp(ts)
	return err
}

func (t Timestamp) String() string {
	return time.Time(t).String()
}
