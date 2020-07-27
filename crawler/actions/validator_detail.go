package actions

import (
	"encoding/json"
	"errors"
	"explorer/model"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"gopkg.in/mgo.v2/bson"
)

func (a *action) MakeBaseInfo(item model.Result, VS *[]model.ValidatorSet) {
	var baseInfo model.ExtraValidatorInfo
	//var vaada model.ValidatorToDelegatorAddress

	baseInfo.Validator = item.OperatorAddress
	_, baseInfo.Address = a.Validator.Check(item.OperatorAddress)
	baseInfo.Identity = item.Description.Identity
	tempValue, _ := strconv.ParseFloat(a.getSelfToken(item.OperatorAddress, baseInfo.Address), 64)
	baseInfo.SelfToken = tempValue
	baseInfo.TotalToken, _ = strconv.ParseFloat(item.Tokens, 64)
	baseInfo.OthersToken = baseInfo.TotalToken - tempValue
	baseInfo.WebSite = item.Description.Website
	baseInfo.Details = item.Description.Details
	baseInfo.HsnHeight = a.getHeight(item.ConsensusPubkey)
	baseInfo.MissedBlockList = a.getMissBlock(VS, item.ConsensusPubkey)

	Sign := a.ValidatorDetail.Check(baseInfo)
	if Sign == 0 {
		//set
		a.ValidatorDetail.Set(baseInfo)
	} else {
		//update
		a.ValidatorDetail.Update(baseInfo)
	}
}

func (a *action) GetDelegations() {
	/* need validator's opAddress*/
	/*get opAddress form db*/
	sign := 0
	errorCount := 0
	var delegationObj model.DelegatorObj
	var delegations model.Delegators
	for {
		//用于标志delegations信息，删除无用的信息。
		if sign > 100 {
			sign = 0
		}

		//var validators model.ValidatorInfo
		// get validator info
		vaList := a.Validator.GetInfo() //vaList == validators List
		if len(*vaList) == 0 {
			time.Sleep(time.Second * 4)
			errorCount++
			if errorCount >= 5 {
				log.Err(errors.New(`get validator error`)).Interface(`vaList`, vaList).Msg(`GetDelegations`)
				time.Sleep(time.Second * 4 * 2)
			}
			continue
		} else {
			if errorCount > 0 {
				errorCount--
			}
		}
		for _, item := range *vaList {
			address := item.ValidatorAddress
			url := a.LcdURL + "/staking/validators/" + address + "/delegations"
			//log.Debug().Interface(`url`,url).Msg(`GetDelegations`)
			rsp, err := a.R().Get(url)
			if err != nil {
				log.Err(err).Interface(`url`, url).Msg(`GetDelegations`)
				time.Sleep(time.Second * 4)
				continue
			}

			err = json.Unmarshal(rsp.Body(), &delegations)
			if err != nil {
				log.Err(err).Interface(`rsp`, rsp).Interface(`url`, url).Msg(`GetDelegations`)
				time.Sleep(time.Second * 4)
				continue
			}

			for _, item := range delegations.Result {
				delegationObj.Shares, _ = strconv.ParseFloat(item.Balance, 64)
				//delegationObj.Shares, _ = strconv.ParseFloat(item.Balance.Amount, 64)
				delegationObj.DelegatorAddress = item.DelegatorAddress
				delegationObj.Address = item.ValidatorAddress
				delegationObj.Sign = sign
				delegationObj.Time = time.Now()
				a.Delegator.SetInfo(delegationObj)
			}

		}
		a.Delegator.DeleteInfo(sign)
		time.Sleep(time.Second * 4)
		sign++
	}

}

func (a *action) GetDelegatorNums() {
	for {
		utcH := time.Now().UTC().Hour()
		utcM := time.Now().UTC().Minute()
		h, _ := time.ParseDuration(strconv.Itoa(23-utcH) + "h")
		m, _ := time.ParseDuration(strconv.Itoa(50-utcM) + "m")
		time.Sleep(h)
		time.Sleep(m)
		err := a.updateInsertDelegatorData()
		if err != nil {
			//validator list is empty.
			continue
		}
		time.Sleep(time.Hour * 1)
	}
}

func (a *action) updateInsertDelegatorData() error {

	var delegations model.Delegators
	var validatorDelegationNums model.ValidatorDelegatorNums
	// get validator info
	vaList := a.Validator.GetInfo() //vaList == validators List
	if len(*vaList) == 0 {
		time.Sleep(time.Second * 4)
		return errors.New("validator list is empty")
	}

	for _, item := range *vaList {
		address := item.ValidatorAddress
		url := a.LcdURL + "/staking/validators/" + address + "/delegations"

		rsp, err := a.R().Get(url)
		if err != nil {
			log.Err(err).Interface(`url`, url).Msg(`updateInsertDelegatorData`)
			time.Sleep(time.Second * 4)
			continue
		}

		err = json.Unmarshal(rsp.Body(), &delegations)
		if err != nil {
			log.Err(err).Interface(`rsp`, rsp).Interface(`url`, url).Msg(`updateInsertDelegatorData`)
			time.Sleep(time.Second * 4)
			continue
		}

		validatorDelegationNums.ValidatorAddress = address
		validatorDelegationNums.DelegatorNums = len(delegations.Result)
		a.Delegator.SetDelegatorCount(validatorDelegationNums)
	}
	return nil
}

func (a *action) getMissBlock(vs *[]model.ValidatorSet, pbKey string) []model.MissBLockData {
	var blockRecords []model.MissBLockData //记录一百个块中该验证着参与的次数（通过公钥）
	var blockRecord model.MissBLockData
OUTLOOP:
	for _, Sets := range *vs {

		for _, item := range Sets.Validators {
			if item.PubKey == pbKey {
				blockRecord.Height = Sets.Height
				blockRecord.State = 1
				blockRecords = append(blockRecords, blockRecord)
				continue OUTLOOP
			}
		}
		blockRecord.Height = Sets.Height
		blockRecord.State = 0
		blockRecords = append(blockRecords, blockRecord)

	}

	return blockRecords
}

func (a *action) getSelfToken(validatorAddress string, accountAddress string) string {
	//http://172.38.8.89:1317/hsn1zqxayv6qe50w6h3ynfj6tq9pr09r7rtuq565clhsnvaloper1zqxayv6qe50w6h3ynfj6tq9pr09r7rtu4u3wgp
	var delegator model.Delegator

	if accountAddress == "" {
		return ""
	}

	url := a.LcdURL + "/staking/delegators/" + accountAddress + "/delegations/" + validatorAddress
	rsp, err := a.R().Get(url)
	if err != nil {
		log.Err(err).Interface(`url`, url).Msg(`getSelfToken`)
		return ""
	}

	err = json.Unmarshal(rsp.Body(), &delegator)
	if err != nil {
		log.Err(err).Interface(`rsp`, rsp).Interface(`url`, url).Msg(`getSelfToken`)
		return ""
	}
	return delegator.Result.Balance
	//return delegator.Result.Balance.Amount
}

func (a *action) getHeight(pubkey string) string {
	var blockSimpleInfo model.BlocksHeights

	_, hash := a.Proposer.CheckValidator(pubkey)

	//_, hash := mappingRealationship.CheckValidator(pbkey)
	//get block info
	conn := a.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("block").Find(bson.M{"block.header.proposeraddress": hash}).One(&blockSimpleInfo)
	return strconv.Itoa(blockSimpleInfo.IntHeight)

}
