package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Linnanmaan sääasema JSON weather
// http://weather.willab.fi/
const weatherURL = "http://weather.willab.fi/weather.json"

// Weather response JSON structure
type Weather struct {
	TempNow         float64 // -17.1,
	TempHi          float64 // -12.2,
	TempLo          float64 // -19.7,
	DewPoint        float64 // -19.9,
	Humidity        float64 // 79,
	AirPressure     float64 // 1012.6,
	WindSpeed       float64 // 2.6,
	WindSpeedMax    float64 // 4.9,
	WindDir         int     // 63,
	Precipitation1d float64 // 0,
	Precipitation1h float64 // 0,
	SolarRad        int     // 10,
	WindChill       float64 // -23.5,
	Timestamp       string  // "2016-01-11 18:21 EET"
}

func main() {
	res, err := http.Get(weatherURL)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Printf("%+v\n", weather)
	fmt.Printf("Oulu %.1f°C\n", weather.TempNow)
}
