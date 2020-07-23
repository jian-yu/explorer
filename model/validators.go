package model

import (
	"time"
)

type Validators struct {
	Height string   `json:"height"`
	Result []Result `json:"result"`
}
type Result struct {
	OperatorAddress string `json:"operator_address"`
	ConsensusPubkey string `json:"consensus_pubkey"`
	Jailed          bool   `json:"jailed"`
	Status          int    `json:"status"`
	Tokens          string `json:"tokens"`
	DelegatorShares string `json:"delegator_shares"`
	Description     struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	} `json:"description"`
	UnbondingHeight string    `json:"unbonding_height"`
	UnbondingTime   time.Time `json:"unbonding_time"`
	Commission      struct {
		CommissionRates struct {
			Rate          string `json:"rate"`
			MaxRate       string `json:"max_rate"`
			MaxChangeRate string `json:"max_change_rate"`
		} `json:"commission_rates"`
		UpdateTime time.Time `json:"update_time"`
	} `json:"commission"`
	MinSelfDelegation string `json:"min_self_delegation"`
}
type ValidatorInfo struct {
	Avater           string `json:"avater"`
	AKA              string `json:"aka"`
	Status           int    `json:"status"`
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


