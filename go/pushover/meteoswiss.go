package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/kaankivilcim/pushover"
)

type CurrentForecastData struct {
	CurrentVersionDirectory string `json:"currentVersionDirectory"`
}

type Forecast struct {
	Path            string `json:"path"`
	TempHigh        string `json:"temp_high"`
	Name            string `json:"name"`
	TempLow         string `json:"temp_low"`
	WeatherSymbolID string `json:"weather_symbol_id"`
}

func getForecastString(id int) string {
	switch id {
	case 1:
		return "sunny"
	case 2:
		return "mostly sunny, some clouds"
	case 3:
		return "partly sunny, thick passing clouds"
	case 4:
		return "overcast"
	case 5:
		return "very cloudy"
	case 6:
		return "sunny intervals, isolated showers"
	case 7:
		return "sunny intervals, isolated sleet"
	case 8:
		return "sunny intervals, snow showers"
	case 9:
		return "overcast, some rain showers"
	case 10:
		return "overcast, some sleet"
	case 11:
		return "overcast, some snow showers"
	case 12:
		return "sunny intervals, chance of thunderstorms"
	case 13:
		return "sunny intervals and thunderstorms"
	case 14:
		return "very cloudy, light rain"
	case 15:
		return "very cloudy, light sleet"
	case 16:
		return "very cloudy, light snow showers"
	case 17:
		return "very cloudy, intermittent rain"
	case 18:
		return "very cloudy, intermittent sleet"
	case 19:
		return "very cloudy, intermittent snow"
	case 20:
		return "very overcast with rain"
	case 21:
		return "very overcast with frequent sleet"
	case 22:
		return "very overcast with heavy snow"
	case 23:
		return "very overcast, slight chance of storms"
	case 24:
		return "very overcast with storms"
	case 25:
		return "very cloudy, very stormy"
	case 26:
		return "high clouds"
	case 27:
		return "stratus"
	case 28:
		return "fog"
	case 29:
		return "sunny intervals, scattered showers"
	case 30:
		return "sunny intervals, scattered snow showers"
	case 31:
		return "sunny intervals, scattered sleet"
	case 32:
		return "sunny intervals, some showers"
	case 33:
		return "short sunny intervals, frequent rain"
	case 34:
		return "short sunny intervals, frequent snowfalls"
	case 35:
		return "overcast and dry"
	case 36:
		return "partly sunny, slightly stormy"
	case 37:
		return "partly sunny, stormy snow showers"
	case 38:
		return "overcast, thundery showers"
	case 39:
		return "overcast, thundery snow showers"
	case 40:
		return "very cloudly, slightly stormy"
	case 41:
		return "overcast, slightly stormy"
	case 42:
		return "very cloudly, thundery snow showers"
	case 101:
		return "clear"
	case 102:
		return "slightly overcast"
	case 103:
		return "heavy cloud formations"
	case 104:
		return "overcast"
	case 105:
		return "very cloudy"
	case 106:
		return "overcast, scattered showers"
	case 107:
		return "overcast, scattered rain and snow showers"
	case 108:
		return "overcast, snow showers"
	case 109:
		return "overcast, some showers"
	case 110:
		return "overcast, some rain and snow showers"
	case 111:
		return "overcast, some snow showers"
	case 112:
		return "slightly stormy"
	case 113:
		return "storms"
	case 114:
		return "very cloudy, light rain"
	case 115:
		return "very cloudy, light rain and snow showers"
	case 116:
		return "very cloudy, light snowfall"
	case 117:
		return "very cloudy, intermittent rain"
	case 118:
		return "very cloudy, intermittent mixed rain and snowfall"
	case 119:
		return "very cloudy, intermittent snowfall"
	case 120:
		return "very cloudy, constant rain"
	case 121:
		return "very cloudy, frequent rain and snowfall"
	case 122:
		return "very cloudy, heavy snowfall"
	case 123:
		return "very cloudy, slightly stormy"
	case 124:
		return "very cloudy, stormy"
	case 125:
		return "very cloudy, storms"
	case 126:
		return "high cloud"
	case 127:
		return "stratus"
	case 128:
		return "fog"
	case 129:
		return "slightly overcast, scattered showers"
	case 130:
		return "slightly overcast, scattered snowfall"
	case 131:
		return "slightly overcast, rain and snow showers"
	case 132:
		return "slightly overcast, some showers"
	case 133:
		return "overcast, frequent rain showers"
	case 134:
		return "overcast, frequent snow showers"
	case 135:
		return "overcast and dry"
	case 136:
		return "slightly overcast, slightly stormy"
	case 137:
		return "slightly overcast, stormy snow showers"
	case 138:
		return "overcast, thundery showers"
	case 139:
		return "overcast, thundery snow showers"
	case 140:
		return "very cloudly, slightly stormy"
	case 141:
		return "overcast, slightly stormy"
	case 142:
		return "very cloudly, thundery snow showers"
	}
	return "unknown"
}

func getForecastURL() string {
	resp, err := http.Get("https://www.meteoswiss.admin.ch/product/output/weather-widget/forecast/versions.json")
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	var data CurrentForecastData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Panicln(err)
	}

	return fmt.Sprintf("https://www.meteoswiss.admin.ch/product/output/weather-pill/%s/en/800100.json", data.CurrentVersionDirectory)
}

func getPushoverToken() string {
	token := os.Getenv("PUSHOVER_TOKEN")
	if len(token) == 0 {
		log.Panicln("Environment variable PUSHOVER_TOKEN was not set")
	}
	return token
}

func getPushoverTarget() string {
	if len(os.Args[1:]) != 1 {
		log.Panicln("Target user/group key was not supplied")
	}
	return os.Args[1]
}

func main() {
	token := getPushoverToken()
	target := getPushoverTarget()

	url := getForecastURL()
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var forecast Forecast
	if err := json.Unmarshal(body, &forecast); err != nil {
		log.Panicln(err)
	}

	id, err := strconv.Atoi(forecast.WeatherSymbolID)
	if err != nil {
		log.Panicln(err)
	}
	desc := getForecastString(id)
	msg := fmt.Sprintf("%s° - %s° %s", forecast.TempLow, forecast.TempHigh, desc)

	m := pushover.New(token, target)
	_, err = m.Send("Weather forecast", msg)
	if err != nil {
		log.Panicln(err)
	}
}
