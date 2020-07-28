package model

import (
	"time"
)

type Delegators struct {
	Address string    `json:"address"`
	Height  string    `json:"height"`
	Result  []DResult `json:"result"`
	Time    time.Time `json:"time"`
}

type DResult struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	Balance          string `json:"balance"`
}

//type Balance struct {
//	Denom  string `json:"denom"`
//	Amount string `json:"amount"`
//}
type DelegatorObj struct {
	Address          string    `json:"address"`
	Height           string    `json:"height"`
	DelegatorAddress string    `json:"delegator_address"`
	Shares           float64   `json:"shares"`
	Sign             int       `json:"sign"`
	Time             time.Time `json:"time"`
}
type ValidatorDelegatorNums struct {
	ValidatorAddress string `json:"validator_address"`
	DelegatorNums    int    `json:"delegator_nums"`
}

type Delegator struct {
	Result DResult `json:"result"`
}

type DelegatorExtra struct {
	Height string `json:"height"`
	Result []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Shares           string `json:"shares"`
		Balance          struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
		Reward string `json:"reward"`
		Name   string `json:"name"`
	} `json:"result"`
}

type DelegatorValidatorReward struct {
	Height string `json:"height"`
	Result []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"result"`
}

type DelegateRewards struct {
	Result struct {
		Total []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"total"`
	} `json:"result"`
}

type Unbonding struct {
	Height string `json:"height"`
	Result []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Entries          []struct {
			CreationHeight string    `json:"creation_height"`
			CompletionTime time.Time `json:"completion_time"`
			InitialBalance string    `json:"initial_balance"`
			Balance        string    `json:"balance"`
		} `json:"entries"`
		Name string `json:"name"`
	} `json:"result"`
}
