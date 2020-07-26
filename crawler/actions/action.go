package actions

import (
	"encoding/json"
	"explorer/common"
	"explorer/db"
	"explorer/model"
	"explorer/utils"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

var logger = log.With().Logger()

var tokenPrice string
var fiveMinAgo time.Time

type action struct {
	LcdURL            string
	RpcURL            string
	ChainName         string
	Denom             string
	CoinPriceURL      string
	VSetCap           int
	CoinToVotingPower float64
	GenesisAddr       string

	*resty.Client
	db.MgoOperator
	common.Validator
	common.Custom
	common.Block
	common.ValidatorDetail
	common.Proposer
	common.Delegator
	common.Transaction
}

func NewAction(
	m db.MgoOperator,
	lcdUrl string,
	rcpURL string,
	chainName string,
	denom string,
	coinPriceURL string,
	vSetCap int,
	coinToVotingPower float64,
	genesisAddr string,
) Action {
	cli := resty.New()
	validator := common.NewValidator(m)
	custom := common.NewCustom(m)
	block := common.NewBlock(m)
	validatorDetail := common.NewValidatorDetail(m)
	proposer := common.NewProposer(m)
	delegator := common.NewDelegator(m)
	transaction := common.NewTransaction(m)

	return &action{
		MgoOperator:       m,
		LcdURL:            lcdUrl,
		RpcURL:            rcpURL,
		ChainName:         chainName,
		Denom:             denom,
		CoinPriceURL:      coinPriceURL,
		VSetCap:           vSetCap,
		GenesisAddr:       genesisAddr,
		Validator:         validator,
		Custom:            custom,
		Block:             block,
		ValidatorDetail:   validatorDetail,
		Proposer:          proposer,
		Delegator:         delegator,
		Transaction:       transaction,
		Client:            cli,
		CoinToVotingPower: coinToVotingPower,
	}
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

func (a *action) getBLockTime(height int) (float64, error) {
	var block model.BlockInfo

	lastHeightUrl := a.LcdURL + "/blocks/" + strconv.Itoa(height)
	aheadHeightUrl := a.LcdURL + "/blocks/" + strconv.Itoa(height-1)

	rsp, err := a.Client.R().Get(lastHeightUrl)
	if err != nil {
		return 0.0, err
	}
	err = json.Unmarshal(rsp.Body(), &block)
	if err != nil {
		return 0.0, err
	}

	lastHeightTime := block.Block.Header.Time

	rsp, err = a.Client.R().Get(aheadHeightUrl)
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

func (a *action) getAllPledgenTokens() decimal.Decimal {
	/* GET PLEDGEN TOKENS FROM DB*/
	var Info model.Information
	conn := a.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("public").Find(nil).Sort("-height").One(&Info)
	tokens := strconv.Itoa(Info.PledgeCoin)
	total, _ := decimal.NewFromString(tokens)
	return total
}

func (a *action) getUptime(vs *[]model.ValidatorSet, pbKey string) int {
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

func (a *action) setValidatorHashAddress(operaAddress string, pubkey string) {
	/*
		Get the validators public key .
		Creat hash mapping between public key and validator and hash values.
	*/
	var validatorAddressAndKey model.ValidatorAddressAndKey

	sign, _ := a.CheckValidator(pubkey)

	if sign == 0 {
		validatorAddressAndKey.ConsensusPubkey = pubkey
		validatorAddressAndKey.OperatorAddress = operaAddress
		validatorAddressAndKey.ProposerHash = utils.GenHexAddrFromPubKey(pubkey)
		a.Proposer.SetInfo(validatorAddressAndKey)
	}
}

func (a *action) dealWithValidatorList(item model.Result, CoinToVoitingPower float64, VS *[]model.ValidatorSet) model.ValidatorInfo {
	a.MakeBaseInfo(item, VS)
	a.setValidatorHashAddress(item.OperatorAddress, item.ConsensusPubkey)

	var validatorInfo model.ValidatorInfo
	validatorInfo.AKA = item.Description.Moniker // get nick name
	validatorInfo.Status = item.Status
	validatorInfo.Avater = "" // avater address
	validatorInfo.ValidatorAddress = item.OperatorAddress
	validatorInfo.Jailed = item.Jailed
	validatorInfo.Commission = item.Commission.CommissionRates.Rate
	othersDelegation, _ := decimal.NewFromString(item.Tokens)
	floatAmount := othersDelegation
	floatCoinToVoitingPower := decimal.NewFromFloat(CoinToVoitingPower)
	tempAmount := floatAmount.Div(floatCoinToVoitingPower)
	validatorInfo.VotingPower.Amount, _ = tempAmount.Float64()
	// may be has some problem
	tempPledgenTokens := a.getAllPledgenTokens()
	if tempPledgenTokens.LessThan(decimal.NewFromFloat(1)) {
		tempPledgenTokens = decimal.NewFromFloat(1.0)
	}
	tempPercent := tempAmount.Div(tempPledgenTokens)
	validatorInfo.VotingPower.Percent, _ = tempPercent.Float64()
	validatorInfo.Uptime = a.getUptime(VS, item.ConsensusPubkey)
	validatorInfo.Time = time.Now()
	return validatorInfo
}

func (a *action) getInflation() (string, error) {
	// return inflation http://localhost:1317/minting/inflation
	var inflation Inflation
	url := a.LcdURL + "/minting/inflation"

	rsp, err := a.R().Get(url)
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

func (a *action) getPriceFormDragonex() string {
	if tokenPrice == "" {
		tokenPrice = a.getPrice()
		fiveMinAgo = time.Now()
	} else {
		now := time.Now()
		m, _ := time.ParseDuration("-1m")
		fiveMinAgoFromNow := now.Add(m * 1)
		if fiveMinAgo.Before(fiveMinAgoFromNow) {
			tokenPrice = a.getPrice()
			fiveMinAgo = time.Now()
		}
	}

	return tokenPrice
}

func (a *action) getPrice() string {
	/*
		30分钟从网站取一次价格
	*/
	var price Price

	rsp, err := a.R().Get(a.CoinPriceURL)

	if err != nil {
		logger.Err(err).Interface(`url`, a.CoinPriceURL).Msg(`getPrice`)
		return tokenPrice
	}

	_ = json.Unmarshal(rsp.Body(), &price)
	if len(price.Data) < 1 {
		// empty
		return tokenPrice
	} else {
		return price.Data[0].ClosePrice
	}
}

func (a *action) pledgenAndTotal() (int, int, int, error) {
	//return pledge and total http://localhost:1317/staking/pool
	// Cannot specify height
	var pledgenAndTotalHsn PledgenAndTotalHsn
	url := a.LcdURL + "/staking/pool"

	rsp, err := a.R().Get(url)
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

func (a *action) getValidatorState() (int, int, error) {
	var validators model.Validators
	bondedUrl := a.LcdURL + "/staking/validators?status=bonded&page=1"
	unbondedUrl := a.LcdURL + "/staking/validators?status=unbonded&page=1"
	unbondingdUrl := a.LcdURL + "/staking/validators?status=unbonding&page=1"
	var jailed = 0
	var total = 0

	rsp, err := a.R().Get(bondedUrl)
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

	rsp, err = a.R().Get(unbondingdUrl)
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

	rsp, err = a.R().Get(unbondedUrl)
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

func (a *action) GetPublic() {
	for {
		price := a.getPriceFormDragonex()
		height, pledgen, total, err := a.pledgenAndTotal()
		if err != nil {
			logger.Err(err).Msg(`pledgenAndTotalHsn`)
			time.Sleep(time.Second * 4)
			continue
		}
		inflation, err := a.getInflation()
		if err != nil {
			logger.Err(err).Msg(`getInflation`)
			time.Sleep(time.Second * 4)
			continue
		}
		onlineV, totalV, err := a.getValidatorState()
		if err != nil {
			logger.Err(err).Msg(`getValidators`)
			time.Sleep(time.Second * 4)
			continue
		}
		blockTime, err := a.getBLockTime(height)
		if err != nil {
			logger.Err(err).Msg(`getBLockTime`)
			time.Sleep(time.Second * 4)
			continue
		}

		a.Custom.SetInfo(model.Information{
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

func (a *action) GetBlock() {
	var block model.BlockInfo
	for {
		lastBlockHeight, publicHeight := a.Block.GetAimHeightAndBlockHeight()
		//check the height difference again
		if publicHeight > lastBlockHeight {
			for publicHeight > lastBlockHeight {
				lastBlockHeight = lastBlockHeight + 1

				url := a.LcdURL + "/blocks/" + strconv.Itoa(lastBlockHeight)
				rsp, err := a.Client.R().Get(url)
				if err != nil {
					lastBlockHeight = lastBlockHeight - 1
					time.Sleep(time.Second * 2)
					continue
				}
				err = json.Unmarshal(rsp.Body(), &block)
				if err != nil {

				}
				a.Block.SetBlock(block)
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func (a *action) GetValidators() {
	for {
		// 获取验证人列表集合 unbonding bonded unbonded
		// http://172.38.8.89:1317/staking/validators?status=unbonding&page=1
		// http://172.38.8.89:1317/staking/validators?status=bonded&page=1
		// http://172.38.8.89:1317/staking/validators?status=unbonded&page=1
		//var validatorList ValidatorList

		var validators model.Validators
		var validatorInfos []model.ValidatorInfo

		ValidatorsSet := a.Validator.GetValidatorSet(a.VSetCap)

		bondedUrl := a.LcdURL + "/staking/validators?status=bonded"
		unbondedUrl := a.LcdURL + "/staking/validators?status=unbonded"
		unbondingdUrl := a.LcdURL + "/staking/validators?status=unbonding"

		rsp, err := a.Client.R().Get(bondedUrl)
		if err != nil {
			logger.Err(err).Msg(`get bonded validator`)
			time.Sleep(time.Second * 2)
			continue
		}
		err = json.Unmarshal(rsp.Body(), &validators)
		if err == nil {
			for _, item := range validators.Result {
				//test
				info := a.dealWithValidatorList(item, a.CoinToVotingPower, ValidatorsSet)
				validatorInfos = append(validatorInfos, info)
			}
		}

		rsp, err = a.Client.R().Get(unbondingdUrl)
		if err != nil {
			logger.Err(err).Msg(`get unbonding validator`)
			time.Sleep(time.Second * 2)
			continue
		}
		err = json.Unmarshal(rsp.Body(), &validators)
		if err == nil {
			for _, item := range validators.Result {
				//test
				info := a.dealWithValidatorList(item, a.CoinToVotingPower, ValidatorsSet)
				validatorInfos = append(validatorInfos, info)
			}
		}

		rsp, err = a.Client.R().Get(unbondedUrl)
		if err != nil {
			logger.Err(err).Msg(`get unbonded validator`)
			time.Sleep(time.Second * 2)
			continue
		}
		err = json.Unmarshal(rsp.Body(), &validators)
		if err == nil {
			for _, item := range validators.Result {
				//test
				info := a.dealWithValidatorList(item, a.CoinToVotingPower, ValidatorsSet)
				validatorInfos = append(validatorInfos, info)
			}
		}

		for _, info := range validatorInfos {
			a.Validator.SetInfo(info)
		}

		time.Sleep(time.Second * 4)
	}
}
