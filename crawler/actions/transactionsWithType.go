package actions

import (
	"explorer/model"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"
)

// get txs list from the listed urls
//http://172.38.8.89:1317/txs?message.action=send
//http://172.38.8.89:1317/txs?message.action=delegate
//http://172.38.8.89:1317/txs?message.action=vote
//http://172.38.8.89:1317/txs?message.action=begin_unbonding
//http://172.38.8.89:1317/txs?message.action=withdraw_delegator_reward
//http://172.38.8.89:1317/txs?message.action=withdraw_validator_commission
//http://172.38.8.89:1317/txs?message.action=multisend
// .. Unfinished

const (
	SendURL             = "/txs?message.action=send"
	DelegateURL         = "/txs?message.action=delegate"
	VoteURL             = "/txs?message.action=vote"
	UnDelegateURL       = "/txs?message.action=begin_unbonding"
	RewardURL           = "/txs?message.action=withdraw_delegator_reward"
	RewardCommissionURL = "/txs?message.action=withdraw_validator_commission"
	MultiSendURL        = "/txs?message.action=multisend"
	ReDelegateURL       = "/txs?message.action=begin_redelegate"
	CreateValidatorURL  = "/txs?message.action=create_validator"
	EditValidatorURL    = "/txs?message.action=edit_validator"
	EditAddressURL      = "/txs?message.action=set_withdraw_address"
	SubmitProposalURL   = "/txs?message.action=submit_proposal"
	DepositURL          = "/txs?message.action=deposit"
)

//type TxMsgField struct {
//	fee                                                                                                        float64
//	realAmount, withDrawRewardAmount, withDrawCommissionAmount                                                 []float64
//	delegatorList, validatorList, fromAddress, toAddress, outputsAddress, inputsAddress, voterAddress, options []string
//}
//
//var MsgTypeHandle = map[string]func(json *simplejson.Json,event *simplejson.Json) *TxMsgField{
//	"cosmos-sdk/MsgSend":                        msgSendHandle,
//	"cosmos-sdk/MsgMultiSend":                   msgMultiSendHandle,
//	"cosmos-sdk/MsgVote":                        msgVoteHandle,
//	"cosmos-sdk/MsgWithdrawValidatorCommission": msgWithdrawValidatorCommissionHandle,
//	"cosmos-sdk/MsgWithdrawDelegationReward":    msgWithdrawDelegationRewardHandle,
//	"cosmos-sdk/MsgDelegate":                    msgDelegateHandle,
//	"cosmos-sdk/MsgBeginRedelegate":             msgBeginRedelegateHandle,
//	"cosmos-sdk/MsgUndelegate":                  msgUndelegateHandle,
//	"cosmos-sdk/MsgCreateValidator":             msgCreateValidatorHandle,
//	"cosmos-sdk/MsgModifyWithdrawAddress":       msgModifyWithdrawAddressHandle,
//	"cosmos-sdk/MsgSubmitProposal":              msgSubmitProposal,
//	"cosmos-sdk/MsgDeposit":                     msgDepositHandle,
//}

var ChainName = ""

func (a *action) GetTxs() {

	ChainName = a.ChainName
	//get the transaction judge whether it is stored in the database
	for {
		a.getTxs(a.LcdURL+SendURL, a.ChainName, "send")
		a.getTxs(a.LcdURL+DelegateURL, a.ChainName, "delegate")
		a.getTxs(a.LcdURL+RewardCommissionURL, a.ChainName, "commission")
		a.getTxs(a.LcdURL+RewardURL, a.ChainName, "reward")
		a.getTxs(a.LcdURL+VoteURL, a.ChainName, "vote")
		a.getTxs(a.LcdURL+UnDelegateURL, a.ChainName, "unbonding")
		a.getTxs(a.LcdURL+MultiSendURL, a.ChainName, "multisend")
		a.getTxs(a.LcdURL+ReDelegateURL, a.ChainName, "redelegate")
		a.getTxs(a.LcdURL+CreateValidatorURL, a.ChainName, "createValidator")
		a.getTxs(a.LcdURL+EditValidatorURL, a.ChainName, "editValidator")
		a.getTxs(a.LcdURL+EditAddressURL, a.ChainName, "editAddress")
		a.getTxs(a.LcdURL+SubmitProposalURL, a.ChainName, "submitProposal")
		a.getTxs(a.LcdURL+DepositURL, a.ChainName, "deposit")
		time.Sleep(time.Second * 4) //Avoid frequent request api
	}

}

//func (a *action) getTx(url string, types string) {
//	var txsInfo model.Txs
//	page := a.getPage(types)
//	if page == 0 {
//		page = 1
//	}
//
//	httpCli := resty.New()
//	rsp, err := httpCli.R().Get(url + "&page=" + strconv.Itoa(page))
//	if err != nil {
//		log.Err(err).Interface(`url`, url).Msg(`getTxs`)
//		time.Sleep(time.Second * 4)
//		return
//	}
//
//	jsonObj, err := simplejson.NewJson(rsp.Body())
//	if err != nil {
//		log.Err(err).Interface(`rsp`, rsp).Interface(`url`, url).Msg(`getTxs`)
//	}
//	jsonTxs, err := jsonObj.Get("txs").Array()
//
//	if err != nil {
//		return
//	}
//
//	txsLen := len(jsonTxs)
//	if txsLen == 0 {
//		return
//	}
//
//	txsError, err := jsonObj.Get("error").String()
//	if err != nil {
//		return
//	}
//
//	if txsError != "" {
//		return
//	}
//
//	txsObj := jsonObj.Get("txs")
//
//	for index := 0; index < txsLen; index++ {
//		//var fee float64
//		txObj := txsObj.GetIndex(index)
//
//		txHash, _ := txObj.Get("txhash").String()
//		if isTxExist := a.Transaction.CheckHash(txHash); isTxExist == 0 {
//			continue
//		}
//		height, _ := txsObj.Get("height").String()
//		status, _ := txsObj.Get("logs").GetIndex(0).Get("success").Bool()
//		txTime, _ := txsObj.Get("timestamp").String()
//		logs, _ := txsObj.Get("logs").Array()
//		pluse := len(logs)
//
//		txObj = txsObj.Get("tx")
//
//		tmf := txMsgHandle(txObj)
//
//		if tmf != nil {
//			txsInfo.Height, _ = strconv.Atoi(height)
//			txsInfo.TxHash = txHash
//			txsInfo.Result = status
//			txsInfo.Page = page
//			txsInfo.Amount = tmf.realAmount
//			txsInfo.Plus = pluse
//			txsInfo.Fee = tmf.fee
//			txsInfo.Type = types
//			txsInfo.Time = time.Now()
//			txsInfo.TxTime = txTime //string to time
//			txsInfo.ValidatorAddress = tmf.validatorList
//			txsInfo.DelegatorAddress = tmf.delegatorList
//			txsInfo.WithDrawCommissionAmout = tmf.withDrawCommissionAmount
//			txsInfo.WithDrawRewardAmout = tmf.withDrawRewardAmount
//			txsInfo.FromAddress = tmf.fromAddress
//			txsInfo.ToAddress = tmf.toAddress
//			txsInfo.OutPutsAddress = tmf.outputsAddress
//			txsInfo.InputsAddress = tmf.inputsAddress
//			txsInfo.VoterAddress = tmf.voterAddress
//			txsInfo.Options = tmf.options
//			a.Transaction.SetInfo(txsInfo)
//		}
//	}
//
//}

/**
  "tx": {
      "type": "cosmos-sdk/StdTx",
      "value": {
          "msg": [
              {
                  "type": "cosmos-sdk/MsgDelegate",
                  "value": {
                      "delegator_address": "hsn1j4yux0ytemqjmcd6z7dej7ermuw2hp9mgwu04a",
                      "validator_address": "hsnvaloper1zqxayv6qe50w6h3ynfj6tq9pr09r7rtu4u3wgp",
                      "amount": {
                          "denom": "hsn",
                          "amount": "8465216965"
                      }
                  }
              }
          ],
          "fee": {
              "amount": [],
              "gas": "200000"
          },
          "signatures": [
              {
                  "pub_key": {
                      "type": "tendermint/PubKeySecp256k1",
                      "value": "Awc87jJrm2k5JKdmDDDvvfrvsTm+nn4MF3V2KY2hByDw"
                  },
                  "signature": "eCUibWhipX+pMq3Uc/Y9lQo76BnhSEilHtu+frMbNighPwR5EUAS86EBzFhEdPKdu/gRp/PXmYAdXiW7NYYqRg=="
              }
          ],
          "memo": ""
      }
*/
//func txMsgHandle(txObj *simplejson.Json) *TxMsgField {
//	msgArrObj := txObj.Get("value").Get("msg")
//	msgArr, _ := msgArrObj.Array()
//	msgArrLen := len(msgArr)
//	if msgArrLen == 0 {
//		return nil
//	}
//	for index := 0; index < msgArrLen; index++ {
//		msgObj := msgArrObj.GetIndex(index)
//		typo, _ := msgObj.Get("type").String()
//		if handle, ok := MsgTypeHandle[typo]; ok {
//			return handle(msgObj)
//		}
//	}
//	return nil
//}

func (a *action) getTxs(url string, chainName string, types string) {
	page := a.getPage(types)
	if page == 0 {
		page = 1
	}
	for {
		tempURL := url
		url := tempURL + "&page=" + strconv.Itoa(page)

		httpCli := resty.New()
		rsp, err := httpCli.R().Get(url)
		if err != nil {
			log.Err(err).Interface(`url`, url).Msg(`getTxs`)
			time.Sleep(time.Second * 4)
			continue
		}
		var txsInfo model.Txs
		jsonObj, err := simplejson.NewJson(rsp.Body())
		if err != nil {
			log.Err(err).Interface(`rsp`, rsp).Interface(`url`, url).Msg(`getTxs`)
		}
		jsonTxs, _ := jsonObj.Get("txs").Array()
		txsError, _ := jsonObj.Get("error").String()
		if txsError != "" {
			break
		}

		lenTxs := len(jsonTxs)

		for i := 0; i < lenTxs; i++ {
			hash, _ := jsonObj.Get("txs").GetIndex(i).Get("txhash").String()
			flage := a.Transaction.CheckHash(hash)
			if flage == 0 {
				height, _ := jsonObj.Get("txs").GetIndex(i).Get("height").String()
				status, _ := jsonObj.Get("txs").GetIndex(i).Get("logs").GetIndex(0).Get("success").Bool()
				txTime, _ := jsonObj.Get("txs").GetIndex(i).Get("timestamp").String()
				feeArray, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("fee").Get("amount").Array()
				var fee float64
				// get fee

				for index := 0; index < len(feeArray); index++ {
					demo, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("fee").Get("amount").GetIndex(index).Get("denom").String()

					if demo == chainName {
						strFee, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("fee").Get("amount").GetIndex(index).Get("amount").String()
						floatFee, _ := strconv.ParseFloat(strFee, 64)
						fee = fee + floatFee
					}
				}
				types := types
				logs, _ := jsonObj.Get("txs").GetIndex(i).Get("logs").Array()
				pluse := len(logs)
				//, amount , validator ,delegator,from ,to
				msgArray, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").Array()
				var realAmount, withDrawRewardAmout, withDrawCommissionAmout []float64
				var delegatorList, validatorList, fromAddress, toAddress, outputsAddress, inputsAddress, voterAddress, options []string
				for index := 0; index < len(msgArray); index++ {
					msgType, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("type").String()
					switch msgType {
					case "cosmos-sdk/MsgSend":
						// get amount,from address
						from, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("from_address").String()
						to, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("to_address").String()
						fromAddress = append(fromAddress, from)
						toAddress = append(toAddress, to)

						amount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").Array()
						for index := 0; index < len(amount); index++ {
							denom, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").GetIndex(index).Get("denom").String()
							log.Debug().Interface(`denom`, denom).Interface(`chainName`, chainName).Msg(`getTxs2`)
							if denom == chainName {
								strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").GetIndex(index).Get("amount").String()
								floatAmount, _ := strconv.ParseFloat(strAmount, 64)
								realAmount = append(realAmount, floatAmount)
							}

						}
					case "cosmos-sdk/MsgMultiSend":
						//Get input calculation amount,output address
						outputsArray, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("outputs").Array()
						for outputIndex := 0; outputIndex < len(outputsArray); outputIndex++ {
							output, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("outputs").GetIndex(outputIndex).Get("address").String()
							outputsAddress = append(outputsAddress, output)
						}

						inputsArray, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("inputs").Array()
						for index2 := 0; index2 < len(inputsArray); index2++ {
							input, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("inputs").GetIndex(index2).Get("address").String()
							inputsAddress = append(inputsAddress, input)

							coinArray, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("inputs").GetIndex(index2).Get("coins").Array()
							for innerIndex := 0; innerIndex < len(coinArray); innerIndex++ {
								denom, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("inputs").GetIndex(index2).Get("coins").GetIndex(innerIndex).Get("denom").String()
								if denom == chainName {
									strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("inputs").GetIndex(index2).Get("coins").GetIndex(innerIndex).Get("amount").String()
									floatAmount, _ := strconv.ParseFloat(strAmount, 64)
									realAmount = append(realAmount, floatAmount)
								}

							}
						}
					case "cosmos-sdk/MsgVote":
						// No amount attribute get voter,options
						voter, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("voter").String()
						option, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("option").String()
						voterAddress = append(voterAddress, voter)
						options = append(options, option)
					case "cosmos-sdk/MsgWithdrawValidatorCommission":
						delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
						validator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_address").String()
						eventsArrery, _ := jsonObj.Get("txs").GetIndex(i).Get("events").Array()
						for iEvent := 0; iEvent < len(eventsArrery); iEvent++ {
							eventsType, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("type").String()
							if eventsType == "withdraw_commission" {
								attributesArrery, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("attributes").Array()
								for iAttributes := 0; iAttributes < len(attributesArrery); iAttributes++ {
									value, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("attributes").GetIndex(iAttributes).Get("value").String()

									lenChanName := len(chainName)
									if len(value) > lenChanName && value[len(value)-lenChanName:] == chainName {
										floatAmout, _ := strconv.ParseFloat(value[0:len(value)-lenChanName], 64)
										withDrawCommissionAmout = append(withDrawCommissionAmout, floatAmout)
									}
									//if len(value) > 3 && value[len(value)-3:len(value)] == chainName {
									//	floatAmout, _ := strconv.ParseFloat(value[0:len(value)-3], 64)
									//	withDrawCommissionAmout = append(withDrawCommissionAmout, floatAmout)
									//}
								}
							}
						}
						if validator != "" {
							validatorList = append(validatorList, validator)
						}
						if delegator != "" {
							delegatorList = append(delegatorList, delegator)
						}
					case "cosmos-sdk/MsgWithdrawDelegationReward":
						delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
						validator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_address").String()
						eventsArrery, _ := jsonObj.Get("txs").GetIndex(i).Get("events").Array()
						for iEvent := 0; iEvent < len(eventsArrery); iEvent++ {
							eventsType, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("type").String()
							if eventsType == "withdraw_rewards" {
								attributesArrery, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("attributes").Array()
								for iAttributes := 0; iAttributes < len(attributesArrery); iAttributes++ {
									key, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("attributes").GetIndex(iAttributes).Get("key").String()
									amountValue, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("attributes").GetIndex(iAttributes).Get("value").String()
									lenChanName := len(chainName)
									//need to fix bug ,
									if len(amountValue) > lenChanName && amountValue[len(amountValue)-lenChanName:] == chainName && key == "amount" {
										floatAmout, _ := strconv.ParseFloat(amountValue[0:len(amountValue)-lenChanName], 64)
										withDrawRewardAmout = append(withDrawRewardAmout, floatAmout)
									}
									// 									need to fix bug ,
									//									if len(amountValue) > 3 && amountValue[len(amountValue)-3:len(amountValue)] == chainName && key == "amount" {
									//										floatAmout, _ := strconv.ParseFloat(amountValue[0:len(amountValue)-3], 64)
									//										withDrawRewardAmout = append(withDrawRewardAmout, floatAmout)
									//									}
								}
							}
						}
						if validator != "" {
							validatorList = append(validatorList, validator)
						}
						if delegator != "" {
							delegatorList = append(delegatorList, delegator)
						}
					case "cosmos-sdk/MsgDelegate":
						delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
						validator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_address").String()
						strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").Get("amount").String()
						floatAmount, _ := strconv.ParseFloat(strAmount, 64)
						realAmount = append(realAmount, floatAmount)
						validatorList = append(validatorList, validator)
						delegatorList = append(delegatorList, delegator)
					case "cosmos-sdk/MsgBeginRedelegate":
						delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
						validatorDst, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_dst_address").String()
						validatorSrc, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_src_address").String()
						strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").Get("amount").String()
						floatAmount, _ := strconv.ParseFloat(strAmount, 64)
						realAmount = append(realAmount, floatAmount, 0)
						validatorList = append(validatorList, validatorDst, validatorSrc)
						delegatorList = append(delegatorList, delegator)
					case "cosmos-sdk/MsgUndelegate":
						// get amount ,delegatorList, validatorList
						delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
						validator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_address").String()
						strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").Get("amount").String()
						floatAmount, _ := strconv.ParseFloat(strAmount, 64)
						realAmount = append(realAmount, floatAmount)
						validatorList = append(validatorList, validator)
						delegatorList = append(delegatorList, delegator)
					case "cosmos-sdk/MsgCreateValidator":
						var mappingValidatorDelegator model.ValidatorToDelegatorAddress
						mappingValidatorDelegator.ValidatorAddress, _ = jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_address").String()
						mappingValidatorDelegator.DelegatorAddress, _ = jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
						sign, _ := a.Validator.Check(mappingValidatorDelegator.ValidatorAddress)
						if sign == 0 {
							a.Validator.SetValidatorToDelegatorAddr(mappingValidatorDelegator)
						}
						strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("value").Get("amount").String()
						floatAmount, _ := strconv.ParseFloat(strAmount, 64)
						realAmount = append(realAmount, floatAmount)
						delegatorList = append(delegatorList, mappingValidatorDelegator.DelegatorAddress)
					case "cosmos-sdk/MsgModifyWithdrawAddress":
						//hash problem
						delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
						delegatorList = append(delegatorList, delegator)
					case "cosmos-sdk/MsgSubmitProposal":
						delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("proposer").String()

						coinArray, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("initial_deposit").Array()
						for innerIndex := 0; innerIndex < len(coinArray); innerIndex++ {
							denom, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("initial_deposit").GetIndex(innerIndex).Get("denom").String()
							if denom == chainName {
								strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("initial_deposit").GetIndex(innerIndex).Get("amount").String()
								floatAmount, _ := strconv.ParseFloat(strAmount, 64)
								realAmount = append(realAmount, floatAmount)
							}

						}

						delegatorList = append(delegatorList, delegator)
					case "cosmos-sdk/MsgDeposit":
						delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("depositor").String()

						coinArray, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").Array()
						for innerIndex := 0; innerIndex < len(coinArray); innerIndex++ {
							denom, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").GetIndex(innerIndex).Get("denom").String()
							if denom == chainName {
								strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").GetIndex(innerIndex).Get("amount").String()
								floatAmount, _ := strconv.ParseFloat(strAmount, 64)
								realAmount = append(realAmount, floatAmount)
							}

						}

						delegatorList = append(delegatorList, delegator)

					}

				}

				txsInfo.Height, _ = strconv.Atoi(height)
				txsInfo.TxHash = hash
				txsInfo.Result = status
				txsInfo.Page = page
				txsInfo.Amount = realAmount
				txsInfo.Plus = pluse
				txsInfo.Fee = fee
				txsInfo.Type = types
				txsInfo.Time = time.Now()
				txsInfo.TxTime = txTime //string to time
				txsInfo.ValidatorAddress = validatorList
				txsInfo.DelegatorAddress = delegatorList
				txsInfo.WithDrawCommissionAmout = withDrawCommissionAmout
				txsInfo.WithDrawRewardAmout = withDrawRewardAmout
				txsInfo.FromAddress = fromAddress
				txsInfo.ToAddress = toAddress
				txsInfo.OutPutsAddress = outputsAddress
				txsInfo.InputsAddress = inputsAddress
				txsInfo.VoterAddress = voterAddress
				txsInfo.Options = options
				a.Transaction.SetInfo(txsInfo)
				//fmt.Println(fromAddress)
				//fmt.Println(toAddress)
				//fmt.Println(outputsAddress)
				//fmt.Println(inputsAddress)
				//fmt.Println(voterAddress,options)
			}
		}
		page++
		//time.Sleep(time.Second * 2)
	}

}

func (a *action) getPage(types string) int {
	var txsInfo model.Txs

	conn := a.MgoOperator.GetDBConn()
	defer conn.Session.Close()
	_ = conn.C("Txs").Find(bson.M{"type": types}).Sort("-height").One(&txsInfo)
	return txsInfo.Page
}

//
//func feeCollection(feeObj *simplejson.Json, chainName string) float64 {
//	var fee float64
//	feeArray, err := feeObj.Array()
//	if err != nil {
//		return 0
//	}
//
//	for index := 0; index < len(feeArray); index++ {
//		demo, err := feeObj.GetIndex(index).Get("denom").String()
//		if err != nil {
//			continue
//		}
//		if demo == chainName {
//			strFee, _ := feeObj.GetIndex(index).Get("amount").String()
//			floatFee, _ := strconv.ParseFloat(strFee, 64)
//			fee = fee + floatFee
//		}
//	}
//	return fee
//}
//
//func msgSendHandle(msgJSONObj *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	// get amount,from address
//	var fromAddress, toAddress []string
//	var realAmount []float64
//
//	from, _ := msgJSONObj.Get("value").Get("from_address").String()
//	to, _ := msgJSONObj.Get("value").Get("to_address").String()
//	fromAddress = append(fromAddress, from)
//	toAddress = append(toAddress, to)
//
//	amountObj := msgJSONObj.Get("value").Get("amount")
//	amountArr, _ := msgJSONObj.Get("value").Get("amount").Array()
//	for index := 0; index < len(amountArr); index++ {
//		denom, _ := amountObj.GetIndex(index).Get("denom").String()
//		if denom == ChainName {
//			strAmount, _ := msgJSONObj.Get("amount").String()
//			floatAmount, _ := strconv.ParseFloat(strAmount, 64)
//			realAmount = append(realAmount, floatAmount)
//		}
//
//	}
//	return &TxMsgField{fromAddress: fromAddress, toAddress: toAddress, realAmount: realAmount}
//}
//
//func msgMultiSendHandle(msgJSON *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	var outputsAddress, inputsAddress []string
//	var realAmount []float64
//
//	outputsObj := msgJSON.Get("value").Get("outputs")
//	outputsArray, _ := msgJSON.Get("value").Get("outputs").Array()
//	for outputIndex := 0; outputIndex < len(outputsArray); outputIndex++ {
//		output, _ := outputsObj.GetIndex(outputIndex).Get("address").String()
//		outputsAddress = append(outputsAddress, output)
//	}
//
//	inputsObj := msgJSON.Get("value").Get("inputs")
//	inputsArray, _ := msgJSON.Get("value").Get("inputs").Array()
//	for index2 := 0; index2 < len(inputsArray); index2++ {
//		input, _ := inputsObj.GetIndex(index2).Get("address").String()
//		inputsAddress = append(inputsAddress, input)
//
//		coinObj := inputsObj.GetIndex(index2).Get("coins")
//		coinArray, _ := inputsObj.GetIndex(index2).Get("coins").Array()
//		for innerIndex := 0; innerIndex < len(coinArray); innerIndex++ {
//			denom, _ := coinObj.GetIndex(innerIndex).Get("denom").String()
//			if denom == ChainName {
//				strAmount, _ := coinObj.GetIndex(innerIndex).Get("amount").String()
//				floatAmount, _ := strconv.ParseFloat(strAmount, 64)
//				realAmount = append(realAmount, floatAmount)
//			}
//		}
//	}
//	return &TxMsgField{outputsAddress: outputsAddress, inputsAddress: inputsAddress, realAmount: realAmount}
//}
//func msgVoteHandle(msgJSON *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	// No amount attribute get voter,options
//	var voterAddress, options []string
//
//	voter, _ := msgJSON.Get("value").Get("voter").String()
//	option, _ := msgJSON.Get("value").Get("option").String()
//	voterAddress = append(voterAddress, voter)
//	options = append(options, option)
//	return &TxMsgField{voterAddress: voterAddress, options: options}
//}
//func msgWithdrawValidatorCommissionHandle(msgJSON *simplejson.Json,event *simplejson.Json) *TxMsgField {
//	delegator, _ := msgJSON.Get("value").Get("delegator_address").String()
//	validator, _ := msgJSON.Get("value").Get("validator_address").String()
//
//	eventObj := event
//	eventsArr, _ := jsonObj.Get("txs").GetIndex(i).Get("events").Array()
//	for iEvent := 0; iEvent < len(eventsArrery); iEvent++ {
//		eventsType, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("type").String()
//		if eventsType == "withdraw_commission" {
//			attributesArrery, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("attributes").Array()
//			for iAttributes := 0; iAttributes < len(attributesArrery); iAttributes++ {
//				value, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("attributes").GetIndex(iAttributes).Get("value").String()
//
//				lenChanName := len(chainName)
//				if len(value) > lenChanName && value[len(value)-lenChanName:] == chainName {
//					floatAmout, _ := strconv.ParseFloat(value[0:len(value)-lenChanName], 64)
//					withDrawCommissionAmout = append(withDrawCommissionAmout, floatAmout)
//				}
//			}
//		}
//	}
//	if validator != "" {
//		validatorList = append(validatorList, validator)
//	}
//	if delegator != "" {
//		delegatorList = append(delegatorList, delegator)
//	}
//	return nil
//}
//
//func msgWithdrawDelegationRewardHandle(msgJSON *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	return nil
//}
//
//func msgDelegateHandle(msgJSON *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	var realAmount []float64
//	var validatorList, delegatorList []string
//
//	delegator, _ := msgJSON.Get("value").Get("delegator_address").String()
//	validator, _ := msgJSON.Get("value").Get("validator_address").String()
//	strAmount, _ := msgJSON.Get("value").Get("amount").Get("amount").String()
//	floatAmount, _ := strconv.ParseFloat(strAmount, 64)
//	realAmount = append(realAmount, floatAmount)
//	validatorList = append(validatorList, validator)
//	delegatorList = append(delegatorList, delegator)
//	return &TxMsgField{validatorList: validatorList, delegatorList: delegatorList, realAmount: realAmount}
//}
//
//func msgBeginRedelegateHandle(msgJSON *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	var realAmount []float64
//	var validatorList, delegatorList []string
//
//	delegator, _ := msgJSON.Get("value").Get("delegator_address").String()
//	validatorDst, _ := msgJSON.Get("value").Get("validator_dst_address").String()
//	validatorSrc, _ := msgJSON.Get("value").Get("validator_src_address").String()
//	strAmount, _ := msgJSON.Get("value").Get("amount").Get("amount").String()
//	floatAmount, _ := strconv.ParseFloat(strAmount, 64)
//	realAmount = append(realAmount, floatAmount, 0)
//	validatorList = append(validatorList, validatorDst, validatorSrc)
//	delegatorList = append(delegatorList, delegator)
//	return &TxMsgField{validatorList: validatorList, delegatorList: delegatorList, realAmount: realAmount}
//}
//
//func msgUndelegateHandle(msgJSON *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	// get amount ,delegatorList, validatorList
//	var realAmount []float64
//	var validatorList, delegatorList []string
//
//	delegator, _ := msgJSON.Get("value").Get("delegator_address").String()
//	validator, _ := msgJSON.Get("value").Get("validator_address").String()
//	strAmount, _ := msgJSON.Get("value").Get("amount").Get("amount").String()
//	floatAmount, _ := strconv.ParseFloat(strAmount, 64)
//	realAmount = append(realAmount, floatAmount)
//	validatorList = append(validatorList, validator)
//	delegatorList = append(delegatorList, delegator)
//	return &TxMsgField{validatorList: validatorList, delegatorList: delegatorList, realAmount: realAmount}
//}
//
//func msgCreateValidatorHandle(msgJSON *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	var realAmount []float64
//	var delegatorList []string
//	var mappingValidatorDelegator model.ValidatorToDelegatorAddress
//
//	mappingValidatorDelegator.ValidatorAddress, _ = msgJSON.Get("value").Get("validator_address").String()
//	mappingValidatorDelegator.DelegatorAddress, _ = msgJSON.Get("value").Get("delegator_address").String()
//	//sign, _ := a.Validator.Check(mappingValidatorDelegator.ValidatorAddress)
//	//if sign == 0 {
//	//	a.Validator.SetValidatorToDelegatorAddr(mappingValidatorDelegator)
//	//}
//	strAmount, _ := msgJSON.Get("value").Get("value").Get("amount").String()
//	floatAmount, _ := strconv.ParseFloat(strAmount, 64)
//	realAmount = append(realAmount, floatAmount)
//	delegatorList = append(delegatorList, mappingValidatorDelegator.DelegatorAddress)
//	return &TxMsgField{realAmount: realAmount, delegatorList: delegatorList}
//}
//func msgModifyWithdrawAddressHandle(msgJSON *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	//hash problem
//	var delegatorList []string
//
//	delegator, _ := msgJSON.Get("value").Get("delegator_address").String()
//	delegatorList = append(delegatorList, delegator)
//	return &TxMsgField{delegatorList: delegatorList}
//}
//func msgDepositHandle(msgJSON *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	var realAmount []float64
//	var delegatorList []string
//
//	delegator, _ := msgJSON.Get("value").Get("proposer").String()
//
//	coinObj := msgJSON.Get("value").Get("initial_deposit")
//	coinArray, _ := msgJSON.Get("value").Get("initial_deposit").Array()
//	for innerIndex := 0; innerIndex < len(coinArray); innerIndex++ {
//		denom, _ := coinObj.GetIndex(innerIndex).Get("denom").String()
//		if denom == ChainName {
//			strAmount, _ := coinObj.GetIndex(innerIndex).Get("amount").String()
//			floatAmount, _ := strconv.ParseFloat(strAmount, 64)
//			realAmount = append(realAmount, floatAmount)
//		}
//	}
//	delegatorList = append(delegatorList, delegator)
//	return &TxMsgField{delegatorList: delegatorList, realAmount: realAmount}
//}
//func msgSubmitProposal(msgJSON *simplejson.Json,event ...*simplejson.Json) *TxMsgField {
//	var realAmount []float64
//	var delegatorList []string
//
//	delegator, _ := msgJSON.Get("value").Get("depositor").String()
//
//	coinObj := msgJSON.Get("value").Get("amount")
//	coinArray, _ := msgJSON.Get("value").Get("amount").Array()
//	for innerIndex := 0; innerIndex < len(coinArray); innerIndex++ {
//		denom, _ := coinObj.GetIndex(innerIndex).Get("denom").String()
//		if denom == ChainName {
//			strAmount, _ := coinObj.GetIndex(innerIndex).Get("amount").String()
//			floatAmount, _ := strconv.ParseFloat(strAmount, 64)
//			realAmount = append(realAmount, floatAmount)
//		}
//	}
//	delegatorList = append(delegatorList, delegator)
//	return &TxMsgField{delegatorList: delegatorList, realAmount: realAmount}
//}
