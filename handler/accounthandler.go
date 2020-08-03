package handler

import (
	"explorer/model"
	"github.com/shopspring/decimal"
)

type AccountHandler struct {
	base *BaseHandler
	*TokensHandler
}

type AccountInfo struct {
	Detail   *model.Account
	BaseInfo *model.BaseInfo
	Tokens   *Tokens
}

type Tokens struct {
	// (amount percentage)
	Available   []decimal.Decimal `json:"available"`
	Delegated   []decimal.Decimal `json:"delegated"`
	Unbonding   []decimal.Decimal `json:"unbonding"`
	Reward      []decimal.Decimal `json:"reward"`
	Commission  []decimal.Decimal `json:"commission"`
	TotalAmount []decimal.Decimal `json:"total_amount"`
}

func NewAccountHandler(
	base *BaseHandler,
	tokens *TokensHandler,
) *AccountHandler {
	return &AccountHandler{
		base:          base,
		TokensHandler: tokens,
	}
}

func (a *AccountHandler) Account(address string) *AccountInfo {
	var baseInfo model.BaseInfo

	tokens := a.AccountTokenInfo(address)

	decimalPrice, _ := decimal.NewFromString(a.base.Custom.GetInfo().Price)
	baseInfo.Address, _ = a.base.Account.GetInfo(address)
	baseInfo.Amount, _ = tokens.TotalAmount[0].Float64()
	baseInfo.RewardAddress = a.base.Account.GetWithDrawAddress(address)
	baseInfo.TotalPrice, _ = tokens.TotalAmount[0].Mul(decimalPrice).Float64()
	baseInfo.Price, _ = decimalPrice.Float64()

	return &AccountInfo{
		BaseInfo: &baseInfo,
		Tokens:   &tokens,
		Detail:   a.base.Account.GetExtraInfo(address),
	}
}

func (a *AccountHandler) AccountTokenInfo(address string) Tokens {
	var tokens Tokens

	reward := a.GetReward(address)
	avail := a.GetAvailableAmount(address)
	com := a.GetTotalCommissionAmount(address)
	unbind := a.GetTotalUnbondingAmount(address)
	delegate := a.GetTotalDelegateAmount(address)

	decimalTotalAmount := reward[0].Add(avail[0].Add(com[0].Add(unbind[0])))

	tokens.Available = append(tokens.Available, avail[0])
	percentage := avail[0].Div(decimalTotalAmount)
	tokens.Available = append(tokens.Available, percentage)

	tokens.Commission = append(tokens.Commission, com[0])
	percentage = com[0].Div(decimalTotalAmount)
	tokens.Commission = append(tokens.Commission, percentage)

	tokens.Unbonding = append(tokens.Unbonding, unbind[0])
	percentage = unbind[0].Div(decimalTotalAmount)
	tokens.Unbonding = append(tokens.Unbonding, percentage)

	tokens.Delegated = append(tokens.Delegated, delegate[0])
	percentage = delegate[0].Div(decimalTotalAmount)
	tokens.Delegated = append(tokens.Delegated, percentage)
	percentage = reward[0].Div(decimalTotalAmount)

	tokens.TotalAmount = append(tokens.TotalAmount, decimalTotalAmount)

	return tokens
}
