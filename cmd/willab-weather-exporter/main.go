package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/joneskoo/willab-weather/pkg/flags"
	"github.com/joneskoo/willab-weather/weather"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// name of the program
const name = "willab-weather-exporter"

func main() {
	var weatherURL flags.URLFlag
	_ = weatherURL.Set(weather.DefaultURL)
	flag.Var(&weatherURL, "url", "URL to get data from")
	listenAddr := flag.String("listen", ":8080", "HTTP listen port for Prometheus metrics")
	flag.Parse()

	ex := newExporter(weatherURL.URL.String())

	log.Printf("%v listening on %v", name, *listenAddr)
	server := http.Server{Addr: *listenAddr}
	http.Handle("/metrics", ex.Handler())
	log.Fatal(server.ListenAndServe())
}

type exporter struct {
	willabWeatherURL string

	reqDuration   prometheus.Histogram
	temp          prometheus.Gauge
	windchill     prometheus.Gauge
	dewpoint      prometheus.Gauge
	humidity      prometheus.Gauge
	airpressure   prometheus.Gauge
	windspeed     prometheus.Gauge
	winddirection prometheus.Gauge
	precipitation prometheus.Gauge
}

func newExporter(weatherURL string) *exporter {
	ex := &exporter{
		willabWeatherURL: weatherURL,

		reqDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Subsystem: "willab",
			Name:      "request_duration_seconds",
			Help:      "Duration of requests made to Willab API",
			Buckets:   prometheus.ExponentialBuckets(0.02, 2, 8),
		}),

		temp: prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: "willab",
			Name:      "temperature_celsius",
			Help:      "Current temperature in Oulu, Linnanmaa",
		}),
		windchill: prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: "willab",
			Name:      "windchill_celsius",
			Help:      "Current wind chill in Oulu, Linnanmaa",
		}),
		dewpoint: prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: "willab",
			Name:      "dewpoint_celsius",
			Help:      "Current dew point in Oulu, Linnanmaa",
		}),
		humidity: prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: "willab",
			Name:      "humidity_ratio",
			Help:      "Current humidity in Oulu, Linnanmaa",
		}),
		airpressure: prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: "willab",
			Name:      "airpressure_hpa",
			Help:      "Current air pressure in Oulu, Linnanmaa",
		}),
		windspeed: prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: "willab",
			Name:      "windspeed_meters_per_second",
			Help:      "Current wind speed in Oulu, Linnanmaa",
		}),
		winddirection: prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: "willab",
			Name:      "winddirection_degrees",
			Help:      "Current wind direction degrees in Oulu, Linnanmaa",
		}),
		precipitation: prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: "willab",
			Name:      "precipitation_mm_per_hour",
			Help:      "Current wind chill in Oulu, Linnanmaa",
		}),
	}
	return ex
}

func (ex *exporter) Handler() http.Handler {
	reg := prometheus.NewRegistry()
	reg.MustRegister(ex.reqDuration, ex.temp, ex.windchill, ex.dewpoint, ex.humidity, ex.airpressure, ex.windspeed, ex.winddirection, ex.precipitation)

	//reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	//reg.MustRegister(collectors.NewGoCollector())
	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ex.refresh()
		h.ServeHTTP(w, r)
	})

}

func (ex *exporter) refresh() {
	start := time.Now()
	w, err := weather.FromURL(ex.willabWeatherURL)
	if err != nil {
		log.Printf("Error retrieving weather: %v", err)
		return
	}
	ex.reqDuration.Observe(time.Since(start).Seconds())
	ex.temp.Set(float64(w.TempNow))
	ex.windchill.Set(float64(w.WindChill))
	ex.dewpoint.Set(float64(w.DewPoint))
	ex.humidity.Set(float64(w.Humidity) / 100)
	ex.airpressure.Set(float64(w.AirPressure))
	ex.windspeed.Set(float64(w.WindSpeed))
	ex.winddirection.Set(float64(w.WindDir))
	ex.precipitation.Set(float64(w.Precipitation1h))
}
