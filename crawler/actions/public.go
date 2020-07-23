package actions

import (
	"encoding/json"
	"explorer/common"
	"explorer/crawler"
	"explorer/db"
	"explorer/model"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"strconv"
	"time"
)

var tokenPrice string
var fiveMinAgo time.Time
var logger = log.With().Logger()

/*
	Dashboard info ,

			Price            float32   `json:"price"`
			Height           int       `json:"height"`
			PledgeCoin        float32   `json:"pledge_hsn"`
			TotalCoin         float32   `json:"total_hsn"`
			Inflation        float32   `json:"inflation"`
			TotalValidators  int       `json:"total_validators"`
			OnlineValidators int       `json:"online_validators"`
			BlockTime     float64   `json:"block_time"`

*/

func GetPublic() {

	mgoStore := db.NewMongoStore()
	customOperator := common.NewCustom(mgoStore)

	//info := models.NewInfomation()

	for {
		price := getPriceFormDragonex()
		height, pledgen, total, err := pledgenAndTotalHsn()
		if err != nil {
			logger.Err(err).Msg(`pledgenAndTotalHsn`)
			time.Sleep(time.Second * 4)
			continue
		}
		inflation, err := getInflation()
		if err != nil {
			logger.Err(err).Msg(`getInflation`)
			time.Sleep(time.Second * 4)
			continue
		}
		onlineV, totalV, err := getValidators()
		if err != nil {
			logger.Err(err).Msg(`getValidators`)
			time.Sleep(time.Second * 4)
			continue
		}
		blockTime, err := getBLockTime(height)
		if err != nil {
			logger.Err(err).Msg(`getBLockTime`)
			time.Sleep(time.Second * 4)
			continue
		}

		customOperator.SetInfo(model.Information{
			Price:            price,
			Height:           height,
			PledgeCoin:       pledgen,
			TotalCoin:        total,
			Inflation:        inflation,
			TotalValidators:  totalV,
			OnlineValidators: onlineV,
			BlockTime:        blockTime,
		})

		time.Sleep(time.Second * 4)
	}
}

func getBLockTime(height int) (float64, error) {
	var block model.BlockInfo
	var httpClient = resty.New()
	lastHeightUrl := crawler.LcdURL + "/blocks/" + strconv.Itoa(height)
	aheadHeightUrl := crawler.LcdURL + "/blocks/" + strconv.Itoa(height-1)

	rsp, err := httpClient.R().Get(lastHeightUrl)
	if err != nil {
		return 0.0, err
	}
	err = json.Unmarshal(rsp.Body(), &block)
	if err != nil {
		return 0.0, err
	}

	lastHeightTime := block.Block.Header.Time

	rsp, err = httpClient.R().Get(aheadHeightUrl)
	if err != nil {
		return 0.0, err
	}
	err = json.Unmarshal(rsp.Body(), &block)
	if err != nil {
		return 0.0, err
	}

	aheadHeightTime := block.Block.Header.Time
	t1, _ := time.Parse(time.RFC3339Nano, lastHeightTime)
	t2, _ := time.Parse(time.RFC3339Nano, aheadHeightTime)
	blockTime := t1.Sub(t2).Seconds()
	return blockTime, nil
}

func getPriceFormDragonex() string {

	if tokenPrice == "" {
		tokenPrice = getPrice()
		fiveMinAgo = time.Now()
	} else {
		now := time.Now()
		m, _ := time.ParseDuration("-1m")
		fiveMinAgoFromNow := now.Add(m * 1)
		if fiveMinAgo.Before(fiveMinAgoFromNow) {
			tokenPrice = getPrice()
			fiveMinAgo = time.Now()
		}
	}

	return tokenPrice
}

/*
	30分钟从网站取一次价格
*/
func getPrice() string {
	var price Price
	var coinPriceURL = viper.GetString(`Public.CoinPriceURL`)
	var httpClient = resty.New()

	rsp, err := httpClient.R().Get(coinPriceURL)
	if err != nil {
		logger.Err(err).Interface(`CoinPriceURL`, coinPriceURL).Msg(`getPrice`)
		return ""
	}
	err = json.Unmarshal(rsp.Body(), &price)
	if err != nil {
		logger.Err(err).Interface(`rsp`, rsp).Msg(`getPrice`)
	}

	if len(price.Data) < 1 {
		return tokenPrice
	} else {
		return price.Data[0].ClosePrice
	}

}

func pledgenAndTotalHsn() (int, int, int, error) {
	//return pledge and total http://localhost:1317/staking/pool
	// Cannot specify height
	var httpClient = resty.New()
	var pledgenAndTotalHsn PledgenAndTotalHsn
	url := crawler.LcdURL + "/staking/pool"

	rsp, err := httpClient.R().Get(url)
	if err != nil {
		return 0, 0, 0, err
	}

	err = json.Unmarshal(rsp.Body(), &pledgenAndTotalHsn)
	if err != nil {

		return 0, 0, 0, err
	}

	bonded, _ := strconv.Atoi(pledgenAndTotalHsn.Result.BondedTokens)
	unbonded, _ := strconv.Atoi(pledgenAndTotalHsn.Result.NotBondedTokens)
	total := bonded + unbonded
	height, _ := strconv.Atoi(pledgenAndTotalHsn.Height)
	return height, bonded, total, nil
}

func getInflation() (string, error) {
	// return inflation http://localhost:1317/minting/inflation
	var httpClient = resty.New()
	var inflation Inflation
	url := crawler.LcdURL + "/minting/inflation"

	rsp, err := httpClient.R().Get(url)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(rsp.Body(), &inflation)
	if err != nil {
		return "", err
	}
	result := inflation.Result
	return result, nil

}
func getValidators() (int, int, error) {
	// bonded, 	unbonding  http://172.38.8.89:1317/staking/validators?status=unbonding&page=1
	//http://172.38.8.89:1317/staking/validators?status=bonded&page=1
	var validators model.Validators
	var jailed int
	var total int
	var httpClient = resty.New()
	bondedUrl := crawler.LcdURL + "/staking/validators?status=bonded&page=1"
	unbondedUrl := crawler.LcdURL + "/staking/validators?status=unbonded&page=1"
	unbondingdUrl := crawler.LcdURL + "/staking/validators?status=unbonding&page=1"

	rsp, err := httpClient.R().Get(bondedUrl)
	if err != nil {
		return 0, 0, err
	}
	err = json.Unmarshal(rsp.Body(), &validators)
	if err != nil {
		return 0, 0, err
	}
	for _, item := range validators.Result {
		if item.Jailed {
			jailed += 1
		}
	}
	total += len(validators.Result)

	rsp, err = httpClient.R().Get(unbondingdUrl)
	if err != nil {
		return 0, 0, err
	}
	err = json.Unmarshal(rsp.Body(), &validators)
	if err != nil {
		return 0, 0, err
	}
	for _, item := range validators.Result {
		if item.Jailed {
			jailed += 1
		}
	}
	total += len(validators.Result)

	rsp, err = httpClient.R().Get(unbondedUrl)
	rsp, err = httpClient.R().Get(unbondingdUrl)
	if err != nil {
		return 0, 0, err
	}
	err = json.Unmarshal(rsp.Body(), &validators)
	if err != nil {
		return 0, 0, err
	}
	for _, item := range validators.Result {
		if item.Jailed {
			jailed += 1
		}
	}

	total += len(validators.Result)
	return total - jailed, total, nil
}

type Inflation struct {
	Height string `json:"height"`
	Result string `json:"result"`
}

type PledgenAndTotalHsn struct {
	Height string `json:"height"`
	Result struct {
		NotBondedTokens string `json:"not_bonded_tokens"`
		BondedTokens    string `json:"bonded_tokens"`
	} `json:"result"`
}

type Price struct {
	Ok   bool `json:"ok"`
	Code int  `json:"code"`
	Data []struct {
		ClosePrice      string `json:"close_price"`
		CurrentVolume   string `json:"current_volume"`
		MaxPrice        string `json:"max_price"`
		MinPrice        string `json:"min_price"`
		OpenPrice       string `json:"open_price"`
		PriceBase       string `json:"price_base"`
		PriceChange     string `json:"price_change"`
		PriceChangeRate string `json:"price_change_rate"`
		Timestamp       int    `json:"timestamp"`
		TotalAmount     string `json:"total_amount"`
		TotalVolume     string `json:"total_volume"`
		UsdtAmount      string `json:"usdt_amount"`
		SymbolID        int    `json:"symbol_id"`
	} `json:"data"`
}
