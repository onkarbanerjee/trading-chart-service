package entities

import "time"

type aggTradeMessage struct {
	Symbol   string `json:"s"`
	Price    string `json:"p"`
	Quantity string `json:"q"`
	Time     int64  `json:"T"`
}

type Tick struct {
	Stream string          `json:"stream"`
	Data   aggTradeMessage `json:"data"`
}

type Candle struct {
	Symbol    string
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}
