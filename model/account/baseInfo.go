package account

/*
struct and operates,store
*/

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
