package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/joneskoo/willab-weather/pkg/flags"
	"github.com/joneskoo/willab-weather/weather"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	name = "willab-weather-exporter"
)

var (
	ttl    = 60 * time.Second
	listen = ":8080"
)

var (
	reqDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Subsystem: "willab",
		Name:      "request_duration_seconds",
		Help:      "Duration of requests made to Willab API",
		Buckets:   prometheus.ExponentialBuckets(0.02, 2, 8),
	})

	temp = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "willab",
		Name:      "temperature_celsius",
		Help:      "Current temperature in Oulu, Linnanmaa",
	})
	windchill = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "willab",
		Name:      "windchill_celsius",
		Help:      "Current wind chill in Oulu, Linnanmaa",
	})
	dewpoint = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "willab",
		Name:      "dewpoint_celsius",
		Help:      "Current dew point in Oulu, Linnanmaa",
	})
	humidity = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "willab",
		Name:      "humidity_ratio",
		Help:      "Current humidity in Oulu, Linnanmaa",
	})
	airpressure = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "willab",
		Name:      "airpressure_hpa",
		Help:      "Current air pressure in Oulu, Linnanmaa",
	})
	windspeed = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "willab",
		Name:      "windspeed_meters_per_second",
		Help:      "Current wind speed in Oulu, Linnanmaa",
	})
	winddirection = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "willab",
		Name:      "winddirection_degrees",
		Help:      "Current wind direction degrees in Oulu, Linnanmaa",
	})
	precipitation = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "willab",
		Name:      "precipitation_mm_per_hour",
		Help:      "Current wind chill in Oulu, Linnanmaa",
	})
)

type updateHandler struct {
	dataURL string
	ticker  <-chan time.Time

	http.Handler
}

func (h updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	select {
	case <-h.ticker:
		refresh(h.dataURL)
	default:
		// use cached data
	}

	h.Handler.ServeHTTP(w, r)
}

func main() {
	var weatherURL flags.URLFlag
	_ = weatherURL.Set(weather.DefaultURL)
	flag.Var(&weatherURL, "url", "URL to get data from")
	flag.StringVar(&listen, "listen", listen, "HTTP listen port for Prometheus metrics")
	flag.DurationVar(&ttl, "ttl", ttl, "Minimum TTL for caching weather data")
	flag.Parse()

	dataURL := weatherURL.URL.String()
	refresh(dataURL)

	log.Printf("%v listening on %v", name, listen)
	server := http.Server{Addr: listen}
	http.Handle("/metrics", updateHandler{
		dataURL: dataURL,
		ticker:  time.Tick(ttl),
		Handler: promhttp.Handler(),
	})
	log.Fatal(server.ListenAndServe())
}

func refresh(dataURL string) {
	start := time.Now()
	w, err := weather.FromURL(dataURL)
	if err != nil {
		log.Printf("Error retrieving weather: %v", err)
		return
	}
	reqDuration.Observe(time.Since(start).Seconds())
	temp.Set(float64(w.TempNow))
	windchill.Set(float64(w.WindChill))
	dewpoint.Set(float64(w.DewPoint))
	humidity.Set(float64(w.Humidity) / 100)
	airpressure.Set(float64(w.AirPressure))
	windspeed.Set(float64(w.WindSpeed))
	winddirection.Set(float64(w.WindDir))
	precipitation.Set(float64(w.Precipitation1h))
}
