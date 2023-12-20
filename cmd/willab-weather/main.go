// Copyright (c) 2016 Joonas Kuorilehto
// License: The MIT License (MIT)

// The command willab-weather shows local weather from weather.willab.fi.
// This is an unofficial command line client.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/joneskoo/willab-weather/pkg/flags"
	"github.com/joneskoo/willab-weather/weather"
)

const defaultTemplate = `Current weather in Oulu, Linnanmaa

    Temperature:    {{ .TempNow }}

                    24 hour low {{ .TempLo }} / high {{ .TempHi }}

    Wind chill:     {{ .WindChill }}
    Dew point:      {{ .DewPoint }}
    Humidity:       {{ .Humidity }}
    Air pressure:   {{ .AirPressure }}
    Wind speed:     {{ .WindSpeed }}   (gusts {{ .WindSpeedMax }})
    Wind direction: {{ .WindDir }}
    Precipitation:  past hour {{ .Precipitation1h }}
                    past day  {{ .Precipitation1d }}

As of {{ .Timestamp }} from <{{ .WeatherURL }}>
`

func main() {
	var (
		weatherURL flags.URLFlag
		tmpl       flags.TemplateFlag
	)
	// Defaults
	weatherURL.Set(weather.DefaultURL)
	tmpl.Set(defaultTemplate)
	flag.Var(&weatherURL, "url", "URL to get data from")
	flag.Var(&tmpl, "template", `Go template for report, e.g. "{{ .TempNow }}"`)
	flag.Parse()

	// Get weather data
	weatherData, err := weather.FromURL(weatherURL.String())
	if err != nil {
		log.Fatalf("Failed to get weather data: %s", err)
	}

	// Print report using customizable template
	if err := tmpl.Execute(os.Stdout, weatherData); err != nil {
		log.Fatalf("Failed to execute template: %s", err)
	}

}
