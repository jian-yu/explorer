package controllers

import (
	"github.com/astaxie/beego"
)

type BlockTxController struct {
	beego.Controller
	Base *BaseController
}

func (btc *BlockTxController) URLMapping() {
	btc.Mapping("BlockTxs", btc.BlockTxs)
}

// @Title
// @Description
// @Success code 0
// @Failure code 1
// @router /blockTxs [get]
func (btc *BlockTxController) BlockTxs() {
	btc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", btc.Ctx.Request.Header.Get("Origin"))
	var respJSON TxBlocks
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
	respJSON.Code = "0"
	respJSON.Msg = "OK"
	respJSON.Data = txsSet
	respJSON.Total = total
	btc.Data["json"] = respJSON
	btc.ServeJSON()
}
