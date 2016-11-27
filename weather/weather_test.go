package weather_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joneskoo/willab-weather/weather"
)

func TestFromURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handler))
	w, err := weather.FromURL(server.URL)
	if err != nil {
		t.Fatalf("Want no error, got err=%s", err)
	}
	if w.TempNow != -4.7 {
		t.Errorf("Want TempNow = %f, got %f", -4.7, w.TempNow)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	fmt.Fprint(w, data)
}

const data = `{"tempnow":-4.7,"temphi":-2.0,"templo":-5.4,"dewpoint":-7.5,"humidity":81,"airpressure":1003.4,"windspeed":3.0,"windspeedmax":8.5,"winddir":14,"precipitation1d":0.0,"precipitation1h":0.0,"solarrad":-1,"windchill":-9.2,"timestamp":"2016-11-27 13:10 EET"}`
