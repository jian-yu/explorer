package crawler

import (
	"explorer/common"
	"explorer/db"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
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

	httpClient *resty.Client

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

	validator := common.NewValidator(m)
	custom := common.NewCustom(m)
	return &crawler{
		LcdURL:       lcdURL,
		RpcURL:       rpcURL,
		VSetCap:      vSetCap,
		ChainName:    chainName,
		CoinPriceURL: coinPriceURL,
		Denom:        denom,
		MgoOperator:  m,
		Validator:    validator,
		Custom:       custom,
	}
}

func (c *crawler) Start() {

}

func (c *crawler) Stop() {

}
