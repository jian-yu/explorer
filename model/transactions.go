package model

import (
	"time"
)

//http://172.38.8.89:1317/txs?message.action=send
type Txs struct {
	TxHash                  string    `json:"tx_hash"`
	Page                    int       `json:"page"`
	Type                    string    `json:"type"`
	Result                  bool      `json:"result"`
	Amount                  []float64 `json:"amount"`
	Fee                     float64   `json:"fee"`
	Height                  int       `json:"height"`
	TxTime                  string    `json:"tx_time"`
	Plus                    int       `json:"pluse"`
	ValidatorAddress        []string  `json:"validator_address"`
	DelegatorAddress        []string  `json:"delegator_address"`
	WithDrawRewardAmout     []float64 `json:"with_draw_reward_amout"`
	WithDrawCommissionAmout []float64 `json:"with_draw_commission_amout"`
	FromAddress             []string  `json:"from_address"`
	ToAddress               []string  `json:"to_address"`
	OutPutsAddress          []string  `json:"out_puts_address"`
	InputsAddress           []string  `json:"inputs_address"`
	VoterAddress            []string  `json:"voter_address"`
	Options                 []string  `json:"options"`
	Time                    time.Time `json:"time"`
	Data                    string    `json:"data"`
	RawLog                  string    `json:"raw_log"`
	Memo                    string    `json:"memo"`
}

type MultiSendMsg struct {
	Value struct {
		Inputs []struct {
			Address string `json:"address"`
			Coins   []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"coins"`
		} `json:"inputs"`
		Outputs []struct {
			Address string `json:"address"`
			Coins   []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"coins"`
		} `json:"outputs"`
	} `json:"value"`
}
type GetRewardCommissionMsg struct {
	Type  string `json:"type"`
	Value struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Amount           string `json:"amount"`
	} `json:"value"`
}
type GetRewardMsg struct {
	Type  string `json:"type"`
	Value struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Amount           string `json:"amount"`
	} `json:"value"`
}
type UndelegateMsg struct {
	Type  string `json:"type"`
	Value struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Amount           struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"amount"`
	} `json:"value"`
}
type DelegateMsg struct {
	Type  string `json:"type"`
	Value struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Amount           struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"amount"`
	} `json:"value"`
}
type VoteMsg struct {
	Type  string `json:"type"`
	Value struct {
		ProposalID string `json:"proposal_id"`
		Voter      string `json:"voter"`
		Option     string `json:"option"`
	} `json:"value"`
}
type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
type Value struct {
	FromAddress string   `json:"from_address"`
	ToAddress   string   `json:"to_address"`
	Amount      []Amount `json:"amount"`
}
type Msg struct {
	Type  string `json:"type"`
	Value Value  `json:"value"`
}
type Gas struct {
	Used   string `json:"used"`
	Wanted string `json:"wanted"`
}
