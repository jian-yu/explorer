package handler

import (
	"encoding/json"
	"explorer/model"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
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

type ValidatorSet struct {
	Height string `json:"height"`
	Result struct {
		BlockHeight string `json:"block_height"`
		Validators  []struct {
			Address          string `json:"address"`
			PubKey           string `json:"pub_key"`
			ProposerPriority string `json:"proposer_priority"`
			VotingPower      string `json:"voting_power"`
		} `json:"validators"`
	} `json:"result"`
}

type StakeValidator struct {
	Validator        model.Validators
	Owner            string
	Power            string
	ProposerPriority string
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

func (v *ValidatorHandler) StakeValidators() []*StakeValidator {
	var stakeValidators []*StakeValidator

	cli := resty.New()
	url := v.base.LcdURL + "/validatorsets/latest"
	resp, err := cli.R().Get(url)
	if err != nil {
		log.Err(err).Interface(`url`, url).Msg("StakeValidators")
		return stakeValidators
	}

	var validatorSet ValidatorSet
	err = json.Unmarshal(resp.Body(), &validatorSet)
	if err != nil {
		return stakeValidators
	}

	for i, val := range validatorSet.Result.Validators {
		resp, err := cli.R().Get(v.base.LcdURL + fmt.Sprintf("/staking/validators?page=%d&limit=1", i+1))
		if err != nil {
			continue
		}

		var stakingVal model.Validators
		err = json.Unmarshal(resp.Body(), &stakingVal)
		if err != nil {
			continue
		}

		_, owner := v.base.Validator.Check(val.Address)
		stakeVal := &StakeValidator{
			Validator:        stakingVal,
			Owner:            owner,
			Power:            val.VotingPower,
			ProposerPriority: val.ProposerPriority,
		}
		stakeValidators = append(stakeValidators, stakeVal)
	}
	return stakeValidators
}
