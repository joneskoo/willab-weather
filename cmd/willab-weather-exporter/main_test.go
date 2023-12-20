package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joneskoo/willab-weather/weather"
)

func TestExporter(t *testing.T) {
	ex, cleanup := testExporter(t)
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

func testExporter(t *testing.T) (ex *exporter, cleanup func()) {
	testurl, cleanup := testServer(t)
	ex = newExporter(testurl)
	return ex, cleanup
}

func testServer(t *testing.T) (testurl string, cleanup func()) {
	// Start a fake server to mock willab weather
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// update data if golden file is missing
		if _, err := os.Stat("testdata/weather.json"); os.IsNotExist(err) {
			updateWeather(t, "testdata/weather.json")
		}
		http.ServeFile(w, r, "testdata/weather.json")
	}))
	return server.URL, func() { server.Close() }
}

func updateWeather(t *testing.T, datafile string) {
	// download from willab.DefaultURL and save to testdata/weather.json
	res, err := http.Get(weather.DefaultURL)
	if err != nil {
		t.Fatalf("HTTP request: %s", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading response: %s", err)
	}

	res.Body.Close()

	err = os.WriteFile(datafile, body, 0644)
	if err != nil {
		t.Fatalf("writing file: %s", err)
	}
}
