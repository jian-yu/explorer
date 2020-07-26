package controllers

import (
	"github.com/astaxie/beego"
)

type BlockTxController struct {
	beego.Controller
	Base *BaseController
}

// @Title
// @Description
// @Success code 0
// @Failure code 1
// @router /
func (btc *BlockTxController) Get() {
	btc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", btc.Ctx.Request.Header.Get("Origin"))
	var respJson TxBlocks
	head, _ := btc.GetInt("head")
	page, _ := btc.GetInt("page")
	size, _ := btc.GetInt("size")
	if size == 0 {
		size = 5
	}

	var txsSet = make([]txInfo, size)
	list, total := btc.Base.Transaction.GetSpecifiedHeight(head, page, size)
	for i, item := range list {
		txsSet[i].Height = item.Height
		txsSet[i].Hash = item.TxHash
		txsSet[i].Fee = item.Fee
		txsSet[i].Result = item.Result
		txsSet[i].Time = item.TxTime
		txsSet[i].Types = item.Type
		txsSet[i].Nums = item.Plus
		if item.Type == "reward" {
			txsSet[i].Amount = getRewardAmount(item.WithDrawRewardAmout)
		} else {
			txsSet[i].Amount = getAmount(item.Amount)
		}

	}
	respJson.Code = "0"
	respJson.Msg = "OK"
	respJson.Data = txsSet
	respJson.Total = total
	btc.Data["json"] = respJson
	btc.ServeJSON()
}
