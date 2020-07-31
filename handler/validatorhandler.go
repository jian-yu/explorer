package handler

import (
	"explorer/model"
)

type ValidatorHandler struct {
	base       *BaseHandler
	delegation *DelegationsHandler
}

type ValidatorTypeList struct {
	Jailed    []*model.ValidatorInfo `json:"jailed"`
	Active    []*model.ValidatorInfo `json:"active"`
	Candidate []*model.ValidatorInfo `json:"candidate"`
}

func NewValidatorHandler(
	base *BaseHandler,
	delegation *DelegationsHandler,
) *ValidatorHandler {
	return &ValidatorHandler{
		base:       base,
		delegation: delegation,
	}
}

func (v *ValidatorHandler) PublicInfo() model.Information {
	return v.base.Custom.GetInfo()
}

func (v *ValidatorHandler) Validators(typo string) []*model.ValidatorInfo {
	var vtl ValidatorTypeList
	var normallyValidatorList []*model.ValidatorInfo
	var validatorList []*model.ValidatorInfo
	var count int
	list := v.base.Validator.GetInfo()

	if len(list) == 0 {
		return nil
	}
	for _, item := range list {
		if item.Jailed {
			lenJailedList := len(vtl.Jailed)
			if lenJailedList == 0 {
				item.Cumulative = item.VotingPower.Percent
			} else {
				aheadItemInJailedList := lenJailedList - 1
				item.Cumulative = vtl.Jailed[aheadItemInJailedList].Cumulative + item.VotingPower.Percent
			}
			vtl.Jailed = append(vtl.Jailed, item)
			validatorList = append(validatorList, item)
		} else {
			lenNormallyValidatorList := len(normallyValidatorList)

			if len(normallyValidatorList) > 0 {
				aheadItemInNormallyList := lenNormallyValidatorList - 1
				item.Cumulative = normallyValidatorList[aheadItemInNormallyList].Cumulative + item.VotingPower.Percent
			} else {
				item.Cumulative = item.VotingPower.Percent
			}
			normallyValidatorList = append(normallyValidatorList, item)
			if count < 100 {
				vtl.Active = append(vtl.Active, item)
				validatorList = append(validatorList, item)
				count++
			} else {
				vtl.Candidate = append(vtl.Candidate, item)
				validatorList = append(validatorList, item)
			}
		}
	}

	switch typo {
	case "jailed":
		return vtl.Jailed
	case "active":
		return vtl.Active
	case "candidate":
		return vtl.Candidate
	default:
		return validatorList
	}
}

func (v *ValidatorHandler) ValidatorAccount(address string) string {
	count, accountAddr := v.base.Validator.Check(address)
	if count > 0 {
		return accountAddr
	}
	return ""
}

func (v *ValidatorHandler) ValidatorDetail(address string) *model.ExtraValidatorInfo {
	return v.base.ValidatorDetail.GetOne(address)
}

func (v *ValidatorHandler) ValidatorByAsset(asset string) *model.ValidatorInfo {
	validators := v.Validators("")
	for _, val := range validators {
		if asset == val.AKA {
			return val
		}
	}
	return nil
}

func (v *ValidatorHandler) Delegations(address string, page int, size int) *DelegationMsg {
	return v.delegation.ValidatorDelegations(address, page, size)
}
