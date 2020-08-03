package web

import (
	"encoding/json"
	"time"
)

type (
	Account struct {
		Address       string          `json:"address"`
		PublicKey     json.RawMessage `json:"public_key"`
		AccountNumber int64           `json:"account_number"`
		Sequence      int64           `json:"sequence"`
		Flags         uint64          `json:"flags"`
		Balances      []AccountToken  `json:"balances"`
	}

	AccountToken struct {
		Symbol string `json:"symbol"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
		Frozen string `json:"frozen"`
	}

	Validator struct {
		AccountAddress     string          `json:"account_address" sql:",notnull, unique"`
		OperatorAddress    string          `json:"operator_address" sql:",notnull, unique"`
		ConsensusPubKey    json.RawMessage `json:"consensus_pubkey" sql:",notnull, unique"`
		ConsensusAddress   string          `json:"consensus_address" sql:",notnull, unique"`
		Jailed             bool            `json:"jailed"`
		Status             string          `json:"status"`
		Tokens             string          `json:"tokens"`
		Power              int64           `json:"power"`
		DelegatorShares    string          `json:"delegator_shares"`
		Description        Description     `json:"description"`
		BondHeight         int64           `json:"bond_height"`
		BondIntraTxCounter int64           `json:"bond_intra_tx_counter"`
		UnbondingHeight    int64           `json:"unbonding_height"`
		UnbondingTime      time.Time       `json:"unbonding_time"`
		Commission         Commission      `json:"commission"`
	}

	// Description wraps validator's description information
	Description struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}

	// Commission wraps validator's commission information
	Commission struct {
		Rate          string    `json:"rate"`
		MaxRate       string    `json:"max_rate"`
		MaxChangeRate string    `json:"max_change_rate"`
		UpdateTime    time.Time `json:"update_time"`
	}

	Token struct {
		Name           string `json:"name"`
		Symbol         string `json:"symbol"`
		OriginalSymbol string `json:"original_symbol"`
		TotalSupply    string `json:"total_supply"`
		Owner          string `json:"owner"`
		Mintable       bool   `json:"mintable"`
	}

	// TxData wraps tx data
	TxData struct {
		ID         int32       `json:"id,omitempty"`
		Height     int64       `json:"height"`
		Result     bool        `json:"result"`
		TxHash     string      `json:"tx_hash"`
		Messages   []Message   `json:"messages"`
		Signatures []Signature `json:"signatures"`
		Memo       string      `json:"memo"`
		Code       uint32      `json:"code"`
		Timestamp  time.Time   `json:"timestamp"`
	}

	// Signature wraps tx signature
	Signature struct {
		Pubkey        string `json:"pubkey"`
		Address       string `json:"address"`
		Sequence      string `json:"sequence"`
		Signature     string `json:"signature"`
		AccountNumber string `json:"account_number"`
	}

	// Message wraps tx message
	Message struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}

	AccountTxs struct {
		TxNums  int          `json:"txNums"`
		TxArray []*AccountTx `json:"txArray"`
	}

	AccountTx struct {
		TxHash        string  `json:"txHash"`
		BlockHeight   int64   `json:"blockHeight"`
		TxType        string  `json:"txType"`
		TimeStamp     int64   `json:"timeStamp"`
		FromAddr      string  `json:"fromAddr"`
		ToAddr        string  `json:"toAddr"`
		Value         float64 `json:"value"`
		TxAsset       string  `json:"txAsset"`
		TxQuoteAsset  string  `json:"txQuoteAsset"`
		TxFee         float64 `json:"txFee"`
		TxAge         int64   `json:"txAge"`
		OrderID       string  `json:"orderId"`
		Data          string  `json:"data,omitempty"`
		Code          int64   `json:"code"`
		Log           string  `json:"log"`
		ConfirmBlocks int64   `json:"confirmBlocks"`
		Memo          string  `json:"memo"`
		Source        int64   `json:"source"`
		HasChildren   int64   `json:"hasChildren"`
	}

	BlockData struct {
		Height        int64     `json:"height"`
		Proposer      string    `json:"proposer"`
		Moniker       string    `json:"moniker"`
		BlockHash     string    `json:"block_hash"`
		ParentHash    string    `json:"parent_hash"`
		NumPrecommits int64     `json:"num_pre_commits" sql:",notnull"`
		NumTxs        int64     `json:"num_txs" sql:"default:0"`
		TotalTxs      int64     `json:"total_txs" sql:"default:0"`
		Txs           []Txs     `json:"txs"`
		Timestamp     time.Time `json:"timestamp" sql:"default:now()"`
	}

	Txs struct {
		Height    int64     `json:"height"`
		Result    bool      `json:"result"`
		TxHash    string    `json:"tx_hash"`
		Messages  []Message `json:"messages"`
		Memo      string    `json:"memo"`
		Code      uint32    `json:"code"`
		Timestamp time.Time `json:"timestamp"`
	}
)

type (
	Param struct {
		Before    int    `form:"before" `
		After     int    `form:"after" `
		Limit     int    `form:"limit"`
		Asset     string `form:"txAsset"`
		Type      string `form:"type"`
		StartTime int64  `form:"starttime"`
		EndTime   int64  `form:"endtime"`
	}
)
