package model

type (
	Supply struct {
		Height string         `json:"height"`
		Result []SupplyResult `json:"result"`
	}
	SupplyResult struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
)
