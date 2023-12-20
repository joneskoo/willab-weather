package weather_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joneskoo/willab-weather/weather"
)

func TestFromURL(t *testing.T) {
	testurl, cleanup := testServer()
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

func testServer() (testurl string, cleanup func()) {
	// Start a fake server to mock willab weather
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		w.Write([]byte(`{"tempnow":-9.3,"templo":-9.4,"temphi":-1.2,"airpressure":983.7,"humidity":74.5,"precipitation1h":0.0,"precipitation1d":0.0,"precipitation1w":0.0,"solarrad":-1,"windspeed":0.3,"windspeedmax":1.1,"winddir":147,"timestamp":"2023-12-20 20:34:47 UTC","windchill":-9.3,"dewpoint":-13.0}`))
	}))
	return server.URL, func() { server.Close() }
}

func mustParseTime(t *testing.T, s string) time.Time {
	tt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("time.Parse(%q): %v", s, err)
	}
	return tt
}
