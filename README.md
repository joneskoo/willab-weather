# willab-weather

Local weather command line client using data from <http://weather.willab.fi/>.

```
$ go get -u github.com/joneskoo/willab-weather
$ $GOPATH/bin/willab-weather -help
Usage of willab-weather:
  -template value
    	Go template for report, e.g. "{{ .TempNow }}"
  -url value
    	URL to get data from (default http://weather.willab.fi/weather.json)

$ willab-weather -template "{{ .TempNow }}"
-4.7 °C%
$ $GOPATH/bin/willab-weather
Current weather in Oulu, Linnanmaa

    Temperature:    -4.7 °C

                    24 hour low -5.4 °C / high -2.0 °C

    Wind chill:     -9.8 °C
    Dew point:      -7.5 °C
    Humidity:       81 %
    Air pressure:   1003.4 hPa
    Wind speed:     4 m/s   (gusts 8 m/s)
    Wind direction: 9°
    Precipitation:  past hour 0.0 mm
                    past day  0.0 mm

As of 2016-11-27 13:07:00 +0200 EET from <http://weather.willab.fi/weather.html>
```

There is no relationship between this client and Willab or VTT.

## Todo

- [ ] Cache result by timestamp
- [x] Support different output formats
- [ ] Example zsh configuration
- [ ] Timeout for request
- [x] Better error handling
