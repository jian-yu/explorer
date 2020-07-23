package actions

import (
	"encoding/json"
	"explorer/common"
	"explorer/conf"
	"explorer/crawler"
	"explorer/crawler/actions/validatorDetails"
	"explorer/db"
	"explorer/model"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

func GetValidators() {
	var httpClient = resty.New()
	mgoStore := db.NewMongoStore()
	validatorOperator := common.NewValidator(mgoStore)

	var coinToVoitingPower  = viper.GetFloat64(`Public.CoinToVoitingPower`)

	for {
		// 获取验证人列表集合 unbonding bonded unbonded
		// http://172.38.8.89:1317/staking/validators?status=unbonding&page=1
		// http://172.38.8.89:1317/staking/validators?status=bonded&page=1
		// http://172.38.8.89:1317/staking/validators?status=unbonded&page=1
		//var validatorList ValidatorList

		var validators model.Validators
		var validatorInfos []model.ValidatorInfo

		ValidatorsSet := validatorOperator.GetValidatorSet(crawler.ValidatorSetCap)

		bondedUrl := crawler.LcdURL + "/staking/validators?status=bonded"
		unbondedUrl := crawler.LcdURL + "/staking/validators?status=unbonded"
		unbondingdUrl := crawler.LcdURL + "/staking/validators?status=unbonding"

		rsp, err := httpClient.R().Get(bondedUrl)
		if err != nil {
			logger.Err(err).Msg(`get bonded validator`)
			time.Sleep(time.Second *2)
			continue
		}
		err = json.Unmarshal(rsp.Body(), &validators)
		if err == nil{
			for _, item := range validators.Result {
				//test
				info := dealWithValidatorList(item, coinToVoitingPower, ValidatorsSet)
				validatorInfos = append(validatorInfos, info)
			}
		}

		rsp, err = httpClient.R().Get(unbondingdUrl)
		if err != nil {
			logger.Err(err).Msg(`get unbonding validator`)
			time.Sleep(time.Second *2)
			continue
		}
		err = json.Unmarshal(rsp.Body(), &validators)
		if err == nil {
			for _, item := range validators.Result {
				//test
				info := dealWithValidatorList(item, coinToVoitingPower, ValidatorsSet)
				validatorInfos = append(validatorInfos, info)
			}
		}

		rsp, err = httpClient.R().Get(unbondedUrl)
		if err != nil {
			logger.Err(err).Msg(`get unbonded validator`)
			time.Sleep(time.Second *2)
			continue
		}
		err = json.Unmarshal(rsp.Body(), &validators)
		if err == nil {
			for _, item := range validators.Result {
				//test
				info := dealWithValidatorList(item, coinToVoitingPower, ValidatorsSet)
				validatorInfos = append(validatorInfos, info)
			}
		}

		for _, info := range validatorInfos {
			validatorOperator.SetInfo(info)
		}

		time.Sleep(time.Second * 4)
	}
}


func getAllPledgenTokens() decimal.Decimal {
	/* GET PLEDGEN TOKENS FROM DB*/
	var Info model.Information
	session := db.NewDBConn()
	defer session.Close()
	dbConn := session.DB(conf.NewConfig().DBName)
	dbConn.C("public").Find(nil).Sort("-height").One(&Info)
	tokens := strconv.Itoa(Info.PledgeCoin)
	decimalTotalHsn, _ := decimal.NewFromString(tokens)
	return decimalTotalHsn
}
func getUptime(vs *[]models.ValidatorsSet, pbKey string) int {
	count := 0 //记录一百个块中该验证着参与的次数（通过公钥）

	for _, Sets := range *vs {
		for _, item := range Sets.Validators {
			if item.PubKey == pbKey {
				count++
			}
		}
	}
	return count
}


func dealWithValidatorList(item models.Result, CoinToVoitingPower float64, VS *[]models.ValidatorsSet) models.ValidatorInfo {
	//time.Sleep(time.Second * 1) // need to fix panic.被除数可能为0
	go validatorDetails.MakeBaseInfo(item, VS, log)
	go SetValidatorHashAddress(item.OperatorAddress, item.ConsensusPubkey, log)
	var validatorInfo models.ValidatorInfo
	validatorInfo.AKA = item.Description.Moniker // get nick name
	validatorInfo.Status = item.Status
	validatorInfo.Avater = ""                    // avater address
	validatorInfo.ValidatorAddress = item.OperatorAddress
	validatorInfo.Jailed = item.Jailed
	validatorInfo.Commission = item.Commission.CommissionRates.Rate
	othersDelegation, _ := decimal.NewFromString(item.Tokens)
	floatAmount := othersDelegation
	floatCoinToVoitingPower := decimal.NewFromFloat(CoinToVoitingPower)
	tempAmount := floatAmount.Div(floatCoinToVoitingPower)
	validatorInfo.VotingPower.Amount, _ = tempAmount.Float64()
	// may be has some problem
	tempPledgenTokens := getAllPledgenTokens()
	if tempPledgenTokens.LessThan(decimal.NewFromFloat(1)) {
		tempPledgenTokens = decimal.NewFromFloat(1.0)
	}
	tempPercent := tempAmount.Div(tempPledgenTokens)
	validatorInfo.VotingPower.Percent, _ = tempPercent.Float64()
	validatorInfo.Uptime = getUptime(VS, item.ConsensusPubkey)
	validatorInfo.Time = time.Now()
	return validatorInfo
}
