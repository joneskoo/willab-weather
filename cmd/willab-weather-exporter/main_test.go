package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExporter(t *testing.T) {
	ex, cleanup := testExporter()
	defer cleanup()
	h := ex.Handler()

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/metrics", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	expectedLines := []string{
		"# HELP willab_airpressure_hpa Current air pressure in Oulu, Linnanmaa\n",
		"# TYPE willab_airpressure_hpa gauge\n",
		"willab_airpressure_hpa 983.7\n",
		"# HELP willab_dewpoint_celsius Current dew point in Oulu, Linnanmaa\n",
		"# TYPE willab_dewpoint_celsius gauge\n",
		"willab_dewpoint_celsius -13\n",
		"# HELP willab_humidity_ratio Current humidity in Oulu, Linnanmaa\n",
		"# TYPE willab_humidity_ratio gauge\n",
		"willab_humidity_ratio 0.745\n",
		"# HELP willab_precipitation_mm_per_hour Current wind chill in Oulu, Linnanmaa\n",
		"# TYPE willab_precipitation_mm_per_hour gauge\n",
		"willab_precipitation_mm_per_hour 0\n",
		"# HELP willab_request_duration_seconds Duration of requests made to Willab API\n",
		"# TYPE willab_request_duration_seconds histogram\n",
		"willab_request_duration_seconds_bucket{le=\"0.02\"} 1\n",
		"willab_request_duration_seconds_bucket{le=\"0.04\"} 1\n",
		"willab_request_duration_seconds_bucket{le=\"0.08\"} 1\n",
		"willab_request_duration_seconds_bucket{le=\"0.16\"} 1\n",
		"willab_request_duration_seconds_bucket{le=\"0.32\"} 1\n",
		"willab_request_duration_seconds_bucket{le=\"0.64\"} 1\n",
		"willab_request_duration_seconds_bucket{le=\"1.28\"} 1\n",
		"willab_request_duration_seconds_bucket{le=\"2.56\"} 1\n",
		"willab_request_duration_seconds_bucket{le=\"+Inf\"} 1\n",
		// Duration sum is not deterministic, so we just check that it's there
		"willab_request_duration_seconds_sum 0.", // omitted
		"willab_request_duration_seconds_count 1\n",
		"# HELP willab_temperature_celsius Current temperature in Oulu, Linnanmaa\n",
		"# TYPE willab_temperature_celsius gauge\n",
		"willab_temperature_celsius -9.3\n",
		"# HELP willab_windchill_celsius Current wind chill in Oulu, Linnanmaa\n",
		"# TYPE willab_windchill_celsius gauge\n",
		"willab_windchill_celsius -9.3\n",
		"# HELP willab_winddirection_degrees Current wind direction degrees in Oulu, Linnanmaa\n",
		"# TYPE willab_winddirection_degrees gauge\n",
		"willab_winddirection_degrees 147\n",
		"# HELP willab_windspeed_meters_per_second Current wind speed in Oulu, Linnanmaa\n",
		"# TYPE willab_windspeed_meters_per_second gauge\n",
		"willab_windspeed_meters_per_second 0.3\n",
	}
	gotErr := false
	for _, line := range expectedLines {
		if !strings.Contains(rec.Body.String(), line) {
			t.Errorf("expected line %q", line)
			gotErr = true
		}
	}
	if gotErr {
		t.Logf("got:\n%s", rec.Body.String())
	}
}

func testExporter() (ex *exporter, cleanup func()) {
	// Start a fake server to mock willab weather
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		w.Write([]byte(`{"tempnow":-9.3,"templo":-9.4,"temphi":-1.2,"airpressure":983.7,"humidity":74.5,"precipitation1h":0.0,"precipitation1d":0.0,"precipitation1w":0.0,"solarrad":-1,"windspeed":0.3,"windspeedmax":1.1,"winddir":147,"timestamp":"2023-12-20 20:34:47 UTC","windchill":-9.3,"dewpoint":-13.0}`))
	}))
	cleanup = func() {
		server.Close()
	}
	ex = newExporter(server.URL)
	return ex, cleanup
}
