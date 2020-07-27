package account

import (
	"explorer/controllers"
	"explorer/model"
	"strings"

	"github.com/astaxie/beego"
	"github.com/shopspring/decimal"
)

type BaseInfoController struct {
	beego.Controller
	Base *controllers.BaseController
	*KindsRewardController
}
type baseInfoerrMsg struct {
	Data error  `json:"data"`
	Msg  string `json:"msg"`
	Code string `json:"code"`
}
type baseInfoMsg struct {
	Data model.BaseInfo `json:"data"`
	Msg  string         `json:"msg"`
	Code string         `json:"code"`
}

/**/
// @Title
// @Description
// @Success code 0
// @Failure code 1
// @router /
func (bic *BaseInfoController) Get() {
	bic.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", bic.Ctx.Request.Header.Get("Origin"))
	address := bic.GetString("address")
	if address == "" || strings.Index(address, bic.Base.Bech32PrefixAccAddr) != 0 || strings.Index(address, bic.Base.Bech32PrefixValAddr) == 0 {
		var msg baseInfoerrMsg
		msg.Data = nil
		msg.Msg = "Delegator address is empty Or Error address!"
		msg.Code = "1"
		bic.Data["json"] = msg
	} else {
		//获取验证人账户信息和获取提款地址
		var baseInfo model.BaseInfo
		//var account model.Account
		//var withdrawAddress model.WithdrawAddress
		//var price model.Information
		decimalPrice, _ := decimal.NewFromString(bic.Base.Custom.GetInfo().Price)
		var msg baseInfoMsg
		baseInfo.Address, _ = bic.Base.Account.GetInfo(address)
		decimalTotalAmount := bic.KindsRewardController.GetAllKindsAmount(address).Data.TotalAmount[0]
		baseInfo.Amount, _ = decimalTotalAmount.Float64()
		baseInfo.RewardAddress = bic.Base.Account.GetWithDrawAddress(address)
		baseInfo.TotalPrice, _ = decimalTotalAmount.Mul(decimalPrice).Float64()
		baseInfo.Price, _ = decimalPrice.Float64()
		msg.Data = baseInfo
		msg.Code = "0"
		msg.Msg = "OK"
		bic.Data["json"] = msg
	}
	bic.ServeJSON()

}
func (bic *BaseInfoController) getDeciamlRewardAmount(address string) decimal.Decimal {
	//var delegateReward model.DelegateRewards
	amount := bic.Base.Account.GetDelegateReward(address)
	decimalAmount, _ := decimal.NewFromString(amount)
	return decimalAmount
}
