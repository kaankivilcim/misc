package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kaankivilcim/pushover"
)

type ExchangeRate struct {
	Change  bool   `json:"change"`
	EndDate string `json:"end_date"`
	Quotes  struct {
		Chfaud struct {
			Change    float64 `json:"change"`
			ChangePct float64 `json:"change_pct"`
			EndRate   float64 `json:"end_rate"`
			StartRate float64 `json:"start_rate"`
		} `json:"CHFAUD"`
	} `json:"quotes"`
	Source    string `json:"source"`
	StartDate string `json:"start_date"`
	Success   bool   `json:"success"`
}

func getCurrencyDataURL() string {
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("https://api.apilayer.com/currency_data/change?start_date=%s&end_date=%s&currencies=AUD&source=CHF", yesterday, today)
}

func getAPILayerKey() string {
	token := os.Getenv("APILAYER_KEY")
	if len(token) == 0 {
		log.Panicln("Environment variable APILAYER_KEY was not set")
	}
	return token
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
	key := getAPILayerKey()
	token := getPushoverToken()
	target := getPushoverTarget()

	url := getCurrencyDataURL()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicln(err)
	}

	req.Header.Set("apikey", key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var rate ExchangeRate
	if err := json.Unmarshal(body, &rate); err != nil {
		log.Panicln(err)
	}

	msg := fmt.Sprintf("CHF to AUD is %0.3f (%0.3f%% change)", rate.Quotes.Chfaud.EndRate, rate.Quotes.Chfaud.ChangePct)
	m := pushover.New(token, target)
	_, err = m.Send("Exchange rate", msg)
	if err != nil {
		log.Panicln(err)
	}
}
