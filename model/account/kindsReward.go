package account

type DelegateRewards struct {
	Result struct {
		Total []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"total"`
	} `json:"result"`
}
