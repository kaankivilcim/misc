package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kaankivilcim/pushover"
)

type Forecast struct {
	Date               string `json:"date"`
	WeatherDescription string `json:"weather_description"`
	TempHigh           int    `json:"temp_high"`
	Weekday            string `json:"weekday"`
	TempLow            int    `json:"temp_low"`
	WeatherSymbolID    int    `json:"weather_symbol_id"`
}

type ForecastLocation []struct {
	Altitude     int         `json:"altitude"`
	Language     interface{} `json:"language,omitempty"`
	CoordX       string      `json:"coord_x"`
	CoordY       string      `json:"coord_y"`
	Version      string      `json:"version,omitempty"`
	LocationID   string      `json:"location_id"`
	LocationType string      `json:"location_type"`
	CityName     string      `json:"city_name"`
	LocationName string      `json:"location_name"`
	MinZoom      int         `json:"min_zoom"`
	Name         string      `json:"name,omitempty"`
	Timestamp    int         `json:"timestamp"`
	Forecasts    []Forecast  `json:"forecasts"`
}

func getForecastURL() string {
	resp, err := http.Get("https://www.meteoswiss.admin.ch/home.html")
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	// Example: https://www.meteoswiss.admin.ch/product/output/forecast-map/version__20220907_2014/en/chmap_20220907.json"
	// Variables: version__YYYYMMDD_HHMM/en/chmap_YYYYMMDD.json
	re := regexp.MustCompile(`product\/output\/forecast-map\/version__\d{8}_(\d{4})`)
	val := re.FindStringSubmatch(string(body))
	if len(val) == 0 {
		log.Panicln("Forecast URL extraction failed")
	}

	today := time.Now().Format("20060102")
	return fmt.Sprintf("https://www.meteoswiss.admin.ch/product/output/forecast-map/version__%s_%s/en/chmap_%s.json", today, val[1], today)
}

func getToken() string {
	token := os.Getenv("PUSHOVER_TOKEN")
	if len(token) == 0 {
		log.Panicln("Environment variable PUSHOVER_TOKEN was not set")
	}
	return token
}

func getTarget() string {
	if len(os.Args[1:]) != 1 {
		log.Panicln("Target user/group key was not supplied")
	}
	return os.Args[1]
}

func main() {
	token := getToken()
	target := getTarget()

	url := getForecastURL()
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Panicln(err)
	}

	req.Header.Set("Referer", "https://www.meteoswiss.admin.ch/home.html?tab=overview")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var forecasts ForecastLocation
	if err := json.Unmarshal(body, &forecasts); err != nil {
		log.Panicln(err)
	}

	var today Forecast
	for _, f := range forecasts {
		if f.CityName == "Zürich" {
			today = f.Forecasts[0]
			break
		}
	}

	desc := strings.ToUpper(today.WeatherDescription[:1]) + today.WeatherDescription[1:]
	msg := fmt.Sprintf("%d° - %d° %s", today.TempLow, today.TempHigh, desc)

	m := pushover.New(token, target)
	_, err = m.Send("Weather forecast", msg)
	if err != nil {
		log.Panicln(err)
	}
}
