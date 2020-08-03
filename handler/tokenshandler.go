package handler

import (
	"explorer/model"
	"github.com/shopspring/decimal"
)

type TokensHandler struct {
	base *BaseHandler
}

type Kinds struct {
	// (amount percentage)
	Available   []decimal.Decimal `json:"available"`
	Delegated   []decimal.Decimal `json:"delegated"`
	Unbonding   []decimal.Decimal `json:"unbonding"`
	Reward      []decimal.Decimal `json:"reward"`
	Commission  []decimal.Decimal `json:"commission"`
	TotalAmount []decimal.Decimal `json:"total_amount"`
}

type ValidatorToken struct {
	Validator *model.ValidatorInfo
	Owner     string
	Extra     *model.ExtraValidatorInfo
}

func NewTokenHandler(base *BaseHandler) *TokensHandler {
	return &TokensHandler{base: base}
}

func (t *TokensHandler) GetAvailableAmount(address string) []decimal.Decimal {
	var availableAmount []decimal.Decimal
	_, strAmount := t.base.Account.GetInfo(address)
	decimalAmount, _ := decimal.NewFromString(strAmount)
	availableAmount = append(availableAmount, decimalAmount)
	return availableAmount
}

// from LCD interface.
//The amount is so large that there is a negative number in the results.
func (t *TokensHandler) GetReward(address string) []decimal.Decimal {
	//var delegateReward model.DelegateRewards
	var decimalAmount decimal.Decimal
	var rewards []decimal.Decimal
	amount := t.base.Account.GetDelegateReward(address)
	//amount := delegateReward.GetDelegateReward(address)
	decimalAmount, _ = decimal.NewFromString(amount)
	rewards = append(rewards, decimalAmount)
	return rewards
}
func (t *TokensHandler) GetTotalDelegateAmount(address string) []decimal.Decimal {
	//var delegators model.Delegators
	var amount decimal.Decimal
	var delegate []decimal.Decimal
	infos := t.base.Account.GetDelegator(address)
	//infos := delegators.GetInfo(address)
	for _, item := range infos.Result {
		decimalAmount, _ := decimal.NewFromString(item.Balance.Amount)
		amount = amount.Add(decimalAmount)
	}
	delegate = append(delegate, amount)
	return delegate
}
func (t *TokensHandler) GetTotalUnbondingAmount(address string) []decimal.Decimal {
	var amount decimal.Decimal
	var unbond []decimal.Decimal
	infos := t.base.Account.GetUnbonding(address)
	for _, item := range infos.Result {
		for _, entrie := range item.Entries {
			decimalAmount, _ := decimal.NewFromString(entrie.Balance)
			amount = amount.Add(decimalAmount)
		}
	}
	unbond = append(unbond, amount)
	return unbond
}
func (t *TokensHandler) GetTotalCommissionAmount(address string) []decimal.Decimal {
	var commission []decimal.Decimal
	var decimalCommissionAmount decimal.Decimal
	commissionTxs := t.base.Transaction.GetDelegatorCommissionTx(address)
	if len(*commissionTxs) == 0 {
		decimalCommissionAmount, _ = decimal.NewFromString("0.0")
	} else {
		for _, item := range *commissionTxs {
			for index, delegator := range item.DelegatorAddress {
				if delegator == address {
					decimalWithDrawCommissionAmount := decimal.NewFromFloat(item.WithDrawCommissionAmout[index])
					decimalCommissionAmount = decimalCommissionAmount.Add(decimalWithDrawCommissionAmount)
				}
			}
		}
	}

	commission = append(commission, decimalCommissionAmount)
	return commission
}

func (t *TokensHandler) GetValidatorsToken() []*ValidatorToken {
	var validatorTokens []*ValidatorToken

	validators := t.base.Validator.GetInfo()

	for _, val := range validators {
		_, owner := t.base.Validator.Check(val.ValidatorAddress)

		valDetail := t.base.ValidatorDetail.GetOne(val.ValidatorAddress)

		valToken := &ValidatorToken{
			Validator: val,
			Owner:     owner,
			Extra:     valDetail,
		}
		validatorTokens = append(validatorTokens, valToken)
	}
	return validatorTokens
}
