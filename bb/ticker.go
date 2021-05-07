package bb

import (
	"encoding/json"
	"log"
	"time"
)

// See bybit REST api referene
// https://bybit-exchange.github.io/docs/inverse/#t-latestsymbolinfo
type TickerRec struct {
	Symbol               string      `json:"symbol"`                 //"BTCUSD",
	BitPrice             json.Number `json:"bid_price"`              // "7230",
	AskPrice             json.Number `json:"ask_price"`              // "7230.5",
	LastPrice            json.Number `json:"last_price"`             // "7230.00",
	LastTickDirection    string      `json:"last_tick_direction"`    // "ZeroMinusTick",
	PrevPrice24h         json.Number `json:"prev_price_24h"`         // "7163.00",
	Price24hPercent      json.Number `json:"price_24h_pcnt"`         // "0.009353",
	HighPrice24h         json.Number `json:"high_price_24h"`         // "7267.50",
	LowPrice24h          json.Number `json:"low_price_24h"`          // "7067.00",
	PrevPrice1h          json.Number `json:"prev_price_1h"`          // "7209.50",
	Price1hPercent       json.Number `json:"price_1h_pcnt"`          // "0.002843",
	MarkPrice            json.Number `json:"mark_price"`             // "7230.31",
	IndexPrice           json.Number `json:"index_price"`            // "7230.14",
	OpenInterest         json.Number `json:"open_interest"`          // 117860186,
	OpenValue            json.Number `json:"open_value"`             // "16157.26",
	TotalTurnover        json.Number `json:"total_turnover"`         // "3412874.21",
	Turnover24h          json.Number `json:"turnover_24h"`           // "10864.63",
	TotalVolume          json.Number `json:"total_volume"`           // 28291403954,
	Volume24h            json.Number `json:"volume_24h"`             // 78053288,
	FundingRate          json.Number `json:"funding_rate"`           // "0.0001",
	PredictedFundingRate json.Number `json:"predicted_funding_rate"` // "0.0001",
	NextFundingTime      string      `json:"next_funding_time"`      // "2019-12-28T00:00:00Z",
	CountdownHour        json.Number `json:"countdown_hour"`         // 2,
	// DeliveryFeeRate        json.Number `json:"delivery_fee_rate"`        // "0",
	// PredictedDeliveryPrice json.Number `json:"predicted_delivery_price"` // "0.00",
	// DeliveryTime           json.Number `json:"delivery_time"`            // ""
}

func TickerRequest() (body string, time time.Time, err error) {
	url := "https://api.bybit.com//v2/public/tickers?symbol=BTCUSD"

	body, time, err = RestRequest(url)

	return body, time, err
}

func TickerMessage(message string) (ticker []TickerRec, err error) {
	err = json.Unmarshal([]byte(message), &ticker)

	if err != nil {
		log.Println(err)
		return ticker, err // error case
	}
	return ticker, nil
}
