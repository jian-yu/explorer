package model

import (
	"time"
)

type Information struct {
	Price            string    `json:"price"`
	Height           int       `json:"height"`
	PledgeCoin       int       `json:"pledge_hsn"`
	TotalCoin        int       `json:"total_hsn"`
	Inflation        string    `json:"inflation"`
	TotalValidators  int       `json:"total_validators"`
	OnlineValidators int       `json:"online_validators"`
	BlockTime        float64   `json:"block_time"`
	Time             time.Time `json:"time"`
}
