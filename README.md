# willab-weather

Local weather command line client using data from <https://weather.willab.fi/>.
This is unofficial and not endorsed by Willab or VTT.

```
$ go get github.com/joneskoo/willab-weather/...
$ willab-weather
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

As of 2016-11-27 13:07:00 +0200 EET from <https://weather.willab.fi/weather.json>
```

Advanced usage

```
$ willab-weather -help
Usage of willab-weather:
  -template value
        Go template for report, e.g. "{{ .TempNow }}"
  -url value
        URL to get data from (default https://weather.willab.fi/weather.json)
$ willab-weather -template "{{ .TempNow }}
"
-9.6 °C
```
