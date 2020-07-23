package crawler

import (
	"errors"
	"explorer/crawler/actions"
	"explorer/crawler/actions/validatorDetails"
	"github.com/spf13/viper"
)

var LcdURL = ""
var RpcURL = ""
var ValidatorSetCap = 0

func OnStart() {
	LcdURL = viper.GetString(`LCD.URL`)
	RpcURL = viper.GetString(`RPC.URL`)
	ValidatorSetCap = viper.GetInt(`Public.ValidatorSetLimit`)

	// check LCD RPC address.
	if lcdURL != "" && rpcURL != "" {
		go actions.GetPublic()
		go actions.GetBlock()
		go actions.GetValidators(config, log)
		go actions.GetValidatorsSet(config, log)
		go actions.GetTxs2(config, log, config.Public.ChainName)
		go validatorDetails.GetDelegations(config, log)
		go validatorDetails.GetDelegatorNums(config, log)
		go actions.SetValidatorDelegatorAddress(config, log)

		//go validatorBlocksTable()
	} else {
		panic(errors.New("LCD/ RPC address is empty,"))
	}

}
