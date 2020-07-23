package common

import (
	"encoding/json"
	"explorer/db"
	"explorer/model"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type account struct {
	db.MgoOperator
	Validator
}

func NewAccount(m db.MgoOperator, v Validator) Account {
	return &account{
		MgoOperator: m,
		Validator:   v,
	}
}

func (a *account) getName(validatorAddress string) string {
	validator := a.Validator.GetOne(validatorAddress)
	return validator.AKA
}

func (a *account) getReward(tokenName, baseURl, validatorAddress, delegatorAddress string) string {
	var reward model.DelegatorValidatorReward
	var httpClient = resty.New()

	url := baseURl + "/distribution/delegators/" + delegatorAddress + "/rewards/" + validatorAddress
	rsp, err := httpClient.R().EnableTrace().Get(url)
	if err != nil {
		logger.Err(err).Interface(`url`, url).Msg(`getReward`)
		return ""
	}

	err = json.Unmarshal(rsp.Body(), &reward)
	if err != nil {
		logger.Err(err).Interface(`rsp`, rsp).Msg(`getReward`)
		return ""
	}
	//	仅仅只返回名称与配置文件中名称相符的数据
	for _, item := range reward.Result {
		if item.Denom == tokenName {
			return item.Amount
		}
	}
	return ""
}

func (a *account) GetInfo(address string) (string, string) {
	var account model.Account
	var Token string
	var httpClient = resty.New()
	lcdURL := viper.GetString(`LCD.URL`)
	denom := viper.GetString(`Public.Denom`)

	url := lcdURL + "/auth/accounts/" + address
	rsp, err := httpClient.R().EnableTrace().Get(url)
	if err != nil {
		logger.Err(err).Interface(`url`, url).Msg(`GetInfo`)
		return "", ""
	}

	_ = json.Unmarshal(rsp.Body(), &account)
	amounts := account.Result.Value.Coins
	for _, amount := range amounts {
		if amount.Denom == denom {
			Token = amount.Amount
		}
	}

	return account.Result.Value.Address, Token
}

func (a *account) GetWithDrawAddress(address string) string {
	var withdrawAddress model.WithdrawAddress
	var httpClient = resty.New()
	lcdURL := viper.GetString(`LCD.URL`)

	url := lcdURL + "/distribution/delegators/" + address + "/withdraw_address"
	rsp, err := httpClient.R().EnableTrace().Get(url)
	if err != nil {
		logger.Err(err).Interface(`url`, url).Msg(`GetWithDrawAddress`)
		return ""
	}

	err = json.Unmarshal(rsp.Body(), &withdrawAddress)
	if err != nil {
		logger.Err(err).Interface(`rsp`, rsp).Msg(`GetWithDrawAddress`)
		return ""
	}
	return withdrawAddress.Result
}

func (a *account) GetDelegator(address string) *model.DelegatorExtra {
	var delegators model.DelegatorExtra
	var httpClient = resty.New()
	lcdURL := viper.GetString(`LCD.URL`)
	chainName := viper.GetString(`Public.ChainName`)

	url := lcdURL + "/staking/delegators/" + address + "/delegations"
	rsp, err := httpClient.R().EnableTrace().Get(url)
	if err != nil {
		logger.Err(err).Interface(`url`, url).Msg(`GetDelegator`)
		return nil
	}

	err = json.Unmarshal(rsp.Body(), &delegators)
	if err != nil {
		logger.Err(err).Interface(`rsp`, rsp).Msg(`GetDelegator`)
	}

	for index, item := range delegators.Result {
		delegators.Result[index].Name = a.getName(item.ValidatorAddress)
		delegators.Result[index].Reward = a.getReward(chainName, lcdURL, item.ValidatorAddress, item.DelegatorAddress)
	}
	return &delegators
}

func (a *account) GetDelegateReward(address string) string {
	var delegateRewards model.DelegateRewards
	var httpClient = resty.New()
	lcdURL := viper.GetString(`LCD.URL`)
	denom := viper.GetString(`Public.Denom`)

	url := lcdURL + "/distribution/delegators/" + address + "/rewards"
	rsp, err := httpClient.R().EnableTrace().Get(url)
	if err != nil {
		logger.Err(err).Interface(`url`, url).Msg(`GetDelegateReward`)
		return "0.0"
	}

	err = json.Unmarshal(rsp.Body(), &delegateRewards)
	if err != nil {
		logger.Err(err).Interface(`rsp`, rsp).Msg(`GetDelegateReward`)
	}

	for _, item := range delegateRewards.Result.Total {
		if item.Denom == denom {
			return item.Amount
		}

	}
	return "0.0"
}

func (a *account) GetUnbonding(address string) *model.Unbonding {
	lcdURL := viper.GetString(`LCD.URL`)
	var httpClient = resty.New()
	var unbonding model.Unbonding

	url := lcdURL + "/staking/delegators/" + address + "/unbonding_delegations"
	rsp, err := httpClient.R().EnableTrace().Get(url)
	if err != nil {
		logger.Err(err).Interface(`url`, url).Msg(`GetUnbonding`)
		return nil
	}

	err = json.Unmarshal(rsp.Body(), &unbonding)
	if err != nil {
		logger.Err(err).Interface(`rsp`, rsp).Msg(`GetUnbonding`)
	}

	for index, item := range unbonding.Result {
		unbonding.Result[index].Name = a.getName(item.ValidatorAddress)
	}
	return &unbonding
}
