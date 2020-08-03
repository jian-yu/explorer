package actions

import (
	"encoding/json"
	"errors"
	"explorer/model"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bitly/go-simplejson"
	"github.com/go-resty/resty/v2"
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

func (a *action) GetTxs2() {
	//for loop
	// get tx height ,last checked height
	// get aim height, public height
	//for loop
	//get height list ( if block.tx >0 ),height > nowHeight ,
	// get txs
	// sleep 30s
	var tx model.Txs
	nowHeight := a.Transaction.GetTxHeight(tx)
	for {
		var index = 0
		blockHeightList := a.Block.GetBlockListIfHasTx(nowHeight)
		lenBlockList := len(blockHeightList)
		if lenBlockList != 0 {
			for index < lenBlockList {
				nowHeight = blockHeightList[index].IntHeight
				err := a.getTxs2(nowHeight, a.LcdURL, a.ChainName)
				if err == nil {
					index++
				} else {
					logger.Err(err).Msg(`GetTxs2`)
				}
			}
		}
		time.Sleep(time.Second * 4)
	}

}
func (a *action) getTxs2(height int, lcdURL string, chainName string) error {
	var httpCli = resty.New()
	var txsInfo model.Txs

	msgsType := ""
	tempURL := lcdURL + "/txs?tx.height="
	strHeight := strconv.Itoa(height)

	url := tempURL + strHeight
	rsp, err := httpCli.R().Get(url)
	if err != nil {
		return err
	}
	jsonObj, err := simplejson.NewJson(rsp.Body())
	if err != nil || jsonObj == nil {
		log.Err(err).Interface(`rsp`, rsp).Interface(`url`, url).Msg(`getTxs`)
		return err
	}

	if jsonObj.Get("txs") == nil {
		return errors.New(`txs is nil`)
	}

	jsonTxs, _ := jsonObj.Get("txs").Array()
	lenTxs := len(jsonTxs)

	if lenTxs == 0 {
		return errors.New(`txs len == 0`)
	}

	jsonErr, _ := jsonObj.Get("error").String()
	if jsonErr != "" {
		return errors.New(jsonErr)
	}

	for i := 0; i < lenTxs; i++ {
		hash, _ := jsonObj.Get("txs").GetIndex(i).Get("txhash").String()
		flage := a.Transaction.CheckHash(hash)
		if flage == 0 {
			height, _ := jsonObj.Get("txs").GetIndex(i).Get("height").String()
			status, _ := jsonObj.Get("txs").GetIndex(i).Get("logs").GetIndex(0).Get("success").Bool()
			txTime, _ := jsonObj.Get("txs").GetIndex(i).Get("timestamp").String()
			rawLog, _ := jsonObj.Get("txs").GetIndex(i).Get("raw_log").String()
			feeArray, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("fee").Get("amount").Array()
			gasWanted, _ := jsonObj.Get("txs").GetIndex(i).Get("gas_wanted").String()
			gasUsed, _ := jsonObj.Get("txs").GetIndex(i).Get("gas_used").String()
			msg, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").Encode()
			sign, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("signatures").Encode()
			var fee float64
			var memo string
			// get fee
			for index := 0; index < len(feeArray); index++ {
				demo, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("fee").Get("amount").GetIndex(index).Get("denom").String()

				if demo == chainName {
					strFee, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("fee").Get("amount").GetIndex(index).Get("amount").String()
					floatFee, _ := strconv.ParseFloat(strFee, 64)
					fee = fee + floatFee
				}
			}
			//logs, _ := jsonObj.Get("txs").GetIndex(i).Get("logs").Array()

			//, amount , validator ,delegator,from ,to
			msgArray, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").Array()
			pluse := len(msgArray)
			var realAmount, withDrawRewardAmout, withDrawCommissionAmout []float64
			var delegatorList, validatorList, fromAddress, toAddress, outputsAddress, inputsAddress, voterAddress, options []string
			var msgRawType string
			for index := 0; index < len(msgArray); index++ {
				msgRawType, _ = jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("type").String()
				switch msgRawType {
				case "cosmos-sdk/MsgSend":
					// get amount,from address
					msgsType = "send"
					from, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("from_address").String()
					to, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("to_address").String()
					fromAddress = append(fromAddress, from)
					toAddress = append(toAddress, to)
					amount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").Array()
					for index := 0; index < len(amount); index++ {
						demo, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").GetIndex(index).Get("denom").String()
						if demo == chainName {
							strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").GetIndex(index).Get("amount").String()
							floatAmount, _ := strconv.ParseFloat(strAmount, 64)
							realAmount = append(realAmount, floatAmount)
						}

					}
				case "cosmos-sdk/MsgMultiSend":
					//Get input calculation amount,output address
					msgsType = "multisend"
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
							log.Debug().Interface(`denom`, denom).Interface(`chainName`, chainName).Msg(`getTxs2`)
							if denom == chainName {
								strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("inputs").GetIndex(index2).Get("coins").GetIndex(innerIndex).Get("amount").String()
								floatAmount, _ := strconv.ParseFloat(strAmount, 64)
								realAmount = append(realAmount, floatAmount)
							}

						}
					}
				case "cosmos-sdk/MsgVote":
					// No amount attribute get voter,options
					msgsType = "vote"
					voter, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("voter").String()
					option, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("option").String()
					voterAddress = append(voterAddress, voter)
					options = append(options, option)
				case "cosmos-sdk/MsgWithdrawValidatorCommission":
					if msgsType == "" || msgsType == "reward" {
						msgsType = "commission"
					}
					delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
					validator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_address").String()
					eventsArrery, _ := jsonObj.Get("txs").GetIndex(i).Get("events").Array()
					if len(withDrawCommissionAmout) == 0 {
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
					}
					if validator != "" {
						validatorList = append(validatorList, validator)
					}
					if delegator != "" {
						delegatorList = append(delegatorList, delegator)
					}
				case "cosmos-sdk/MsgWithdrawDelegationReward":
					if msgsType == "" {
						msgsType = "reward"
					}
					delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
					validator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_address").String()
					eventsArrery, _ := jsonObj.Get("txs").GetIndex(i).Get("events").Array()
					if len(withDrawRewardAmout) == 0 {
						for iEvent := 0; iEvent < len(eventsArrery); iEvent++ {
							eventsType, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("type").String()
							if eventsType == "withdraw_rewards" {
								attributesArrery, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("attributes").Array()
								for iAttributes := 0; iAttributes < len(attributesArrery); iAttributes++ {
									key, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("attributes").GetIndex(iAttributes).Get("key").String()
									amountValue, _ := jsonObj.Get("txs").GetIndex(i).Get("events").GetIndex(iEvent).Get("attributes").GetIndex(iAttributes).Get("value").String()
									lenChanName := len(chainName)
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
					}
					if validator != "" {
						validatorList = append(validatorList, validator)
					}
					if delegator != "" {
						delegatorList = append(delegatorList, delegator)
					}
				case "cosmos-sdk/MsgDelegate":
					msgsType = "delegate"
					delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
					validator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_address").String()
					strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").Get("amount").String()
					floatAmount, _ := strconv.ParseFloat(strAmount, 64)
					realAmount = append(realAmount, floatAmount)
					validatorList = append(validatorList, validator)
					delegatorList = append(delegatorList, delegator)
				case "cosmos-sdk/MsgBeginRedelegate":
					msgsType = "redelegate"
					delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
					validatorDst, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_dst_address").String()
					validatorSrc, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_src_address").String()
					strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").Get("amount").String()
					floatAmount, _ := strconv.ParseFloat(strAmount, 64)
					realAmount = append(realAmount, floatAmount, 0)
					validatorList = append(validatorList, validatorDst, validatorSrc)
					delegatorList = append(delegatorList, delegator)
				case "cosmos-sdk/MsgUndelegate":
					msgsType = "unbonding"
					// get amount ,delegatorList, validatorList
					delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
					validator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("validator_address").String()
					strAmount, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("amount").Get("amount").String()
					floatAmount, _ := strconv.ParseFloat(strAmount, 64)
					realAmount = append(realAmount, floatAmount)
					validatorList = append(validatorList, validator)
					delegatorList = append(delegatorList, delegator)
				case "cosmos-sdk/MsgCreateValidator":
					msgsType = "createValidator"
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
					msgsType = "editAddress"
					delegator, _ := jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("value").Get("delegator_address").String()
					delegatorList = append(delegatorList, delegator)
				case "cosmos-sdk/MsgSubmitProposal":
					msgsType = "submitProposal"
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
					msgsType = "deposit"
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
				case "cosmos-sdk/MsgEditValidator":
					msgsType = "editValidator"
				case "cosmos-sdk/MsgUnjail":
					msgsType = "unjail"
				default:
					msgsType = msgRawType
				}
				memo, _ = jsonObj.Get("txs").GetIndex(i).Get("tx").Get("value").Get("msg").GetIndex(index).Get("memo").String()
			}

			url := a.LcdURL + "/txs/%s"
			rsp, err := a.Client.R().Get(fmt.Sprintf(url, hash))
			if err != nil {
				continue
			}

			var txHashMap map[string]interface{}
			err = json.Unmarshal(rsp.Body(), &txHashMap)
			if err != nil {
				continue
			}
			var txData string
			if v, ok := txHashMap["data"]; ok {
				txData = v.(string)
			}

			t, _ := time.ParseInLocation(time.RFC3339Nano, txTime,time.Local)

			txsInfo.Height, _ = strconv.Atoi(height)
			txsInfo.TxHash = hash
			txsInfo.Result = status
			txsInfo.RawType = msgRawType
			txsInfo.Page = 0
			txsInfo.Amount = realAmount
			txsInfo.Plus = pluse
			txsInfo.Fee = fee
			txsInfo.Type = msgsType
			txsInfo.Time = time.Now()
			txsInfo.TxTime = t.Unix() //string to time
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
			txsInfo.RawLog = rawLog
			txsInfo.Memo = memo
			txsInfo.Data = txData
			txsInfo.GasWanted = gasWanted
			txsInfo.GasUsed = gasUsed
			txsInfo.Message = msg
			txsInfo.Sign = sign
			a.Transaction.SetInfo(txsInfo)
		}
	}
	//page++
	return nil
}
