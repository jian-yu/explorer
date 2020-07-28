package model

import (
	"time"
)

/*

 */

type BasicInfo struct {
	Avater           string `json:"avater"`
	AKA              string `json:"aka"`
	Jailed           bool   `json:"jailed"`
	ValidatorAddress string `json:"validator_address"`
	VotingPower      struct {
		Amount  float64 `json:"amount"`
		Percent float64 `json:"percent"`
	} `json:"voting_power"`
	Cumulative float64   `json:"cumulative"`
	Commission string    `json:"commission"`
	Uptime     int       `json:"uptime"`
	Time       time.Time `json:"time"`
}

//validator's website,details,identity ,hsn height, address

type ExtraValidatorInfo struct {
	WebSite         string          `json:"web_site"`
	Details         string          `json:"details"`
	Identity        string          `json:"identity"`
	Address         string          `json:"address"`
	Validator       string          `json:"validator"`
	HsnHeight       string          `json:"hsn_height"`
	TotalToken      float64         `json:"total_token"`
	SelfToken       float64         `json:"self_token"`
	OthersToken     float64         `json:"others_token"`
	MissedBlockList []MissBLockData `json:"missed_block_list"`
}
type MissBLockData struct {
	Height int `json:"height"`
	State  int `json:"state"`
}

type Account struct {
	Height string `json:"height"`
	Result struct {
		Type  string `json:"type"`
		Value struct {
			Address string `json:"address"`
			Coins   []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"coins"`
			PublicKey struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"public_key"`
			AccountNumber string `json:"account_number"`
			Sequence      string `json:"sequence"`
		} `json:"value"`
	} `json:"result"`
}

type WithdrawAddress struct {
	Height string `json:"height"`
	Result string `json:"result"`
}

type BaseInfo struct {
	Address       string  `json:"address"`
	RewardAddress string  `json:"reward_address"`
	Amount        float64 `json:"amount"`
	TotalPrice    float64 `json:"total_price"`
	Price         float64 `json:"price"`
}
