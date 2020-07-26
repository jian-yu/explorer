package crawler

import (
	"explorer/common"
	"explorer/crawler/actions"
	"explorer/db"
	"github.com/spf13/viper"
)

type Crawler interface {
	Start()
	Stop()
}

type crawler struct {
	LcdURL       string
	RpcURL       string
	ChainName    string
	Denom        string
	CoinPriceURL string
	VSetCap      int
	GenesisAddr string
	VotingPower float64

	db.MgoOperator
	common.Validator
	common.Custom
}

func New(m db.MgoOperator) Crawler {
	lcdURL := viper.GetString(`LCD.URL`)
	rpcURL := viper.GetString(`RPC.URL`)
	vSetCap := viper.GetInt(`Public.ValidatorSetLimit`)
	chainName := viper.GetString(`Public.ChainName`)
	denom := viper.GetString(`Public.Denom`)
	coinPriceURL := viper.GetString(`Public.CoinPriceURL`)
	genesisAddr := viper.GetString(`Public.CoinPriceURL`)
	votingPower := viper.GetFloat64(`Public.CoinToVoitingPower`)

	return &crawler{
		LcdURL:       lcdURL,
		RpcURL:       rpcURL,
		VSetCap:      vSetCap,
		ChainName:    chainName,
		CoinPriceURL: coinPriceURL,
		Denom:        denom,
		GenesisAddr: genesisAddr,
		VotingPower: votingPower,
		MgoOperator:  m,
	}
}

func (c *crawler) Start() {

	act := actions.NewAction(
		c.MgoOperator,
		c.LcdURL,
		c.RpcURL,
		c.ChainName,
		c.Denom,
		c.CoinPriceURL,
		c.VSetCap,
		c.VotingPower,
		c.GenesisAddr,
		)

	go act.GetPublic()
	go act.GetBlock()
	go act.GetValidators()
	go act.GetValidatorsSet()
	go act.GetDelegations()
	go act.GetDelegatorNums()
	go act.GetTxs()
	go act.GetTxs2()
}

func (c *crawler) Stop() {

}
