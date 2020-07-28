package model

import (
	"time"
)

type ValidatorSet struct {
	BlockHeight string    `json:"block_height"`
	Height      int       `json:"height"`
	Time        time.Time `json:"time"`
	Validators  []struct {
		Address     string `json:"address"`
		PubKey      string `json:"pub_key"`
		VotingPower string `json:"voting_power"`
	} `json:"validators"`
}