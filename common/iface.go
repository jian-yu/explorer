package common

import (
	"explorer/model"
	"github.com/rs/zerolog/log"
)

var logger = log.With().Logger()

type Validator interface {
	SetInfo(info model.ValidatorInfo)
	GetInfo() *[]model.ValidatorInfo
	DeleteAllInfo()
	GetOne(address string) *model.ValidatorInfo
	GetValidatorRank(amount float64, jailed bool) int
	SetValidatorSet(vs model.ValidatorSet)
	GetValidatorSet(limit int) *[]model.ValidatorSet
	SetValidatorToDelegatorAddr(v2d model.ValidatorToDelegatorAddress)
	CheckDelegatorAddress(address string) (string, string)
}

type Block interface {
	SetBlock(b model.BlockInfo)
	GetAimHeightAndBlockHeight() (int, int)
	GetBlockListIfHasTx(height int) []model.BlocksHeights
}

type Custom interface {
	SetInfo(info model.Information)
	GetInfo()model.Information
}

type Transaction interface {
	SetInfo(info model.Txs)
	GetInfo(head int, page int, size int)([]model.Txs, int)
	GetDetail(txHash string)model.Txs
	CheckHash(txHash string)int
	GetValidatorsTransactions(address string)
	GetPowerEventInfo(address string, page, size int) (*[]model.Txs, int)
	GetDelegatorTxs(address string, page, size int) (*[]model.Txs, int)
	GetDelegatorCommissionTx(address string) *[]model.Txs
	GetDelegatorRewardTx(address string) *[]model.Txs
	GetSpecifiedHeight(head int, page int, size int) ([]model.Txs, int)
	GetTxHeight(tx model.Txs) int
}

type Proposer interface {
	SetInfo(vaak model.ValidatorAddressAndKey)
	GetInfo(address string) string
	CheckValidator(pubkey string) (int, string)
	GetValidator(address string) string
}

type ValidatorDetail interface {
	Set(info model.ExtraValidatorInfo)
	Update(info model.ExtraValidatorInfo)
	Check(info model.ExtraValidatorInfo)int
	GetOne(address string) *model.ExtraValidatorInfo
}

type Delegator interface {
	GetDelegatorObj(address string, page int, size int) (*[]model.DelegatorObj, int)
	SetDelegatorObj(delegator model.DelegatorObj)
	SetDelegatorCount(vDelegator model.ValidatorDelegatorNums)
	GetDelegatorCount(address string) int
	DeleteInfo(sign int)
}

type Account interface {
	GetInfo(address string)(string,string)
	GetWithDrawAddress(address string) string
	GetDelegator(address string)*model.DelegatorExtra
	GetDelegateReward(address string) string
	GetUnbonding(address string) *model.Unbonding
}