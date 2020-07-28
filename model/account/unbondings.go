package account

import "time"

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
