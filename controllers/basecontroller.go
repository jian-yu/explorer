package controllers

import (
	"explorer/common"
	"explorer/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var logger = log.With().Logger()

type BaseController struct {
	LcdURL              string
	RpcURL              string
	Bech32PrefixAccAddr string
	Bech32PrefixValAddr string

	db.MgoOperator
	common.Validator
	common.Custom
	common.Block
	common.ValidatorDetail
	common.Proposer
	common.Delegator
	common.Transaction
	common.Account
}

func NewBaseController(m db.MgoOperator) *BaseController {
	lcdURL := viper.GetString(`LCD.URL`)
	rpcURL := viper.GetString(`RPC.URL`)
	bech32PrefixAccAddr := viper.GetString(`Public.Bech32PrefixAccAddr`)
	bech32PrefixValAddr := viper.GetString(`Public.Bech32PrefixValAddr`)

	validator := common.NewValidator(m)
	custom := common.NewCustom(m)
	block := common.NewBlock(m)
	validatorDetail := common.NewValidatorDetail(m)
	proposer := common.NewProposer(m)
	delegator := common.NewDelegator(m)
	transaction := common.NewTransaction(m)
	account := common.NewAccount(m, validator)

	return &BaseController{
		LcdURL:              lcdURL,
		RpcURL:              rpcURL,
		Bech32PrefixAccAddr: bech32PrefixAccAddr,
		Bech32PrefixValAddr: bech32PrefixValAddr,
		MgoOperator:         m,
		Validator:           validator,
		Custom:              custom,
		Block:               block,
		ValidatorDetail:     validatorDetail,
		Proposer:            proposer,
		Delegator:           delegator,
		Transaction:         transaction,
		Account:             account,
	}
}
