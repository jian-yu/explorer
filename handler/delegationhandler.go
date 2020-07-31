package handler

import "explorer/model"

type DelegationsHandler struct {
	base *BaseHandler
}

func NewDelegationsHandler(
	base *BaseHandler,
) *DelegationsHandler {
	return &DelegationsHandler{
		base: base,
	}
}

type DelegationMsg struct {
	TotalDelegations     int          `json:"total_delegation"`
	OneDayAgoDelegations int          `json:"one_day_ago_delegations"`
	Delegations          []*Delegation `json:"delegations"`
}
type Delegation struct {
	Address          string  `json:"address"`
	Amount           float64 `json:"amount"`
	AmountPercentage float64 `json:"share"`
}

func (d *DelegationsHandler) ValidatorDelegations(address string, page int, size int) *DelegationMsg {
	var dMsg DelegationMsg
	var delegations []*Delegation

	validatorBaseInfo := d.base.ValidatorDetail.GetOne(address)
	items, totalDelegations := d.base.Delegator.GetInfo(address, page, size)
	oneDayAgoDelegations := d.base.Delegator.GetDelegatorCount(address)
	for _, item := range items {
		var delegation Delegation
		delegation.Amount = getShares(items, item.DelegatorAddress)
		delegation.Address = item.DelegatorAddress
		delegation.AmountPercentage = getPercentage(delegation.Amount, validatorBaseInfo.TotalToken)
		delegations = append(delegations, &delegation)
	}
	dMsg.TotalDelegations = totalDelegations
	dMsg.OneDayAgoDelegations = totalDelegations - oneDayAgoDelegations
	dMsg.Delegations = delegations

	return &dMsg
}

func getShares(items []*model.DelegatorObj, address string) float64 {
	var amount float64
	for _, item := range items {
		if item.DelegatorAddress == address {
			share := item.Shares
			amount = amount + share
		}
	}
	return amount
}
func getPercentage(amout float64, totalToken float64) float64 {

	return amout / totalToken

}
