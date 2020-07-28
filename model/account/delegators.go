package account

type Delegators struct {
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
