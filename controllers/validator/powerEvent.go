package validator

import (
	"explorer/controllers"
	"explorer/model"
	"strings"

	"github.com/astaxie/beego"
)

type PowerEventController struct {
	beego.Controller
	Base *controllers.BaseController
}

type powerEvents struct {
	Height int     `json:"height"`
	Hash   string  `json:"hash"`
	Sign   int     `json:"sign"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
	Time   string  `json:"time"`
}

type powerEventMsg struct {
	Code  string        `json:"code"`
	Data  []powerEvents `json:"data"`
	Msg   string        `json:"msg"`
	Total int           `json:"total"`
}
type powerEventErrMsg struct {
	Code string `json:"code"`
	Data error  `json:"data"`
	Msg  string `json:"msg"`
}

// @Title Get
// @Description get txs (delegate and undelegate)
// @Success code 0
// @Failure code 1
// @router /
func (pwc *PowerEventController) Get() {
	pwc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", pwc.Ctx.Request.Header.Get("Origin"))
	address := pwc.GetString("address", "")
	page, _ := pwc.GetInt("page", 0)
	size, _ := pwc.GetInt("size", 0)
	if address == "" || strings.Index(address, pwc.Base.Bech32PrefixValAddr) != 0 {
		var err powerEventErrMsg
		err.Data = nil
		err.Msg = "Validator address is empty! Or error address!"
		err.Code = "1"
		pwc.Data["json"] = err
	} else {
		//var txs model.Txs
		// get data form mongodb
		// type redelegate,delegate,unbonding ,validator's address == address ,
		txList, total := pwc.Base.Transaction.GetPowerEventInfo(address, page, size)
		//txList, total := txs.GetPowerEventInfo(address, page, size)
		var pe powerEvents
		var msg powerEventMsg
		msg.Total = total
		msg.Code = "0"
		msg.Msg = "OK"
		for _, item := range *txList {
			pe.Time = item.TxTime
			pe.Height = item.Height
			pe.Hash = item.TxHash
			pe.Type = item.Type
			pe.Amount, pe.Sign = getTxValidatorAmountAndSigns(address, item)
			msg.Data = append(msg.Data, pe)
		}
		pwc.Data["json"] = msg
	}
	pwc.ServeJSON()
}
func getTxValidatorAmountAndSigns(address string, item model.Txs) (float64, int) {
	var tempAmount float64
	sing := 0
	if item.Type == "unbonding" {
		sing = 1
	}
	if item.Type == "redelegate" && item.ValidatorAddress[0] != address {
		sing = 1
	}
	if len(item.Amount) == 0 {
		tempAmount = 0.0
		return tempAmount, sing
	}
	for i := range item.ValidatorAddress {
		tempAmount = tempAmount + item.Amount[i]
	}
	return tempAmount, sing
}
