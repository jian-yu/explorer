package actions

import (
	"encoding/json"
	"explorer/model"

	"github.com/rs/zerolog/log"
)

func (a *action) SetValidatorDelegatorAddress() {
	var genesisValidator model.GenesisFile
	var mappingValidatorDelegator model.ValidatorToDelegatorAddress
	genesisAddress := a.GenesisAddr

	rsp, err := a.R().Get(genesisAddress)
	if err != nil {
		log.Err(err).Interface(`url`, genesisAddress).Msg(`SetValidatorDelegatorAddress`)
		return
	}

	err = json.Unmarshal(rsp.Body(), &genesisValidator)
	if err != nil {
		log.Err(err).Interface(`rsp`, rsp).Msg(`SetValidatorDelegatorAddress`)
		return
	}

	obj := genesisValidator.Result.Genesis.AppState.Genutil.Gentxs
	for _, item := range obj {
		mappingValidatorDelegator.ValidatorAddress = item.Value.Msg[0].Value.ValidatorAddress
		mappingValidatorDelegator.DelegatorAddress = item.Value.Msg[0].Value.DelegatorAddress
		sign, _ := a.Validator.Check(item.Value.Msg[0].Value.ValidatorAddress)
		if sign == 0 {
			a.Validator.SetValidatorToDelegatorAddr(mappingValidatorDelegator)
		}
	}

}
