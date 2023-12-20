package weather_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joneskoo/willab-weather/weather"
)

func TestFromURL(t *testing.T) {
	testurl, cleanup := testServer(t)
	defer cleanup()

	w, err := weather.FromURL(testurl)
	if err != nil {
		t.Fatalf("Want no error, got err=%s", err)
	}

	tt := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"TempNow", w.TempNow, weather.Temperature(-9.3)},
		{"TempHi", w.TempHi, weather.Temperature(-1.2)},
		{"TempLo", w.TempLo, weather.Temperature(-9.4)},
		{"DewPoint", w.DewPoint, weather.Temperature(-13.0)},
		{"Humidity", w.Humidity, weather.Humidity(74.5)},
		{"AirPressure", w.AirPressure, weather.Pressure(983.7)},
		{"WindSpeed", w.WindSpeed, weather.Windspeed(0.3)},
		{"WindSpeedMax", w.WindSpeedMax, weather.Windspeed(1.1)},
		{"WindDir", w.WindDir, weather.Winddirection(147)},
		{"Precipitation1d", w.Precipitation1d, weather.Precipitation(0.0)},
		{"Precipitation1h", w.Precipitation1h, weather.Precipitation(0.0)},
		{"SolarRad", w.SolarRad, -1},
		{"WindChill", w.WindChill, weather.Temperature(-9.3)},
		{"Timestamp", w.Timestamp, weather.Timestamp(mustParseTime(t, "2023-12-20T20:34:47Z"))},
	}
	for _, tc := range tt {
		if tc.got != tc.want {
			t.Errorf("Want %s=%v, got %s=%v", tc.name, tc.want, tc.name, tc.got)
		}
	}
}

func mustParseTime(t *testing.T, s string) time.Time {
	tt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("time.Parse(%q): %v", s, err)
	}
	return tt
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
