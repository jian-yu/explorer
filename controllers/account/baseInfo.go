package account

import (
	"encoding/json"
	"explorer/controllers"
	"explorer/handler"
	"explorer/model"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/shopspring/decimal"
)

type BaseInfoController struct {
	beego.Controller
	Base *controllers.BaseController
	*KindsRewardController
	*handler.AccountHandler
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

type Account struct {
	Address       string          `json:"address"`
	PublicKey     json.RawMessage `json:"public_key"`
	AccountNumber int64           `json:"account_number"`
	Sequence      int64           `json:"sequence"`
	Flags         uint64          `json:"flags"`
	Balances      []AccountToken  `json:"balances"`
}

type AccountToken struct {
	Symbol string `json:"symbol"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
	Frozen string `json:"frozen"`
}

func (bic *BaseInfoController) URLMapping() {
	bic.Mapping("AccountInfo", bic.AccountInfo)
	bic.Mapping("Account", bic.Account)
}

/**/
// @Title
// @Description
// @Success code 0
// @Failure code 1
// @router /accountInfo [get]
func (bic *BaseInfoController) AccountInfo() {
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

/**/
// @Title
// @Description
// @Success code 0
// @Failure code 1
// @router /account/:address [get]
func (bic *BaseInfoController) Account() {
	bic.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", bic.Ctx.Request.Header.Get("Origin"))

	address := bic.Ctx.Input.Param("address")
	if address == "" {
		bic.Data["json"] = &Account{}
		bic.ServeJSON()
		return
	}

	info := bic.AccountHandler.Account(address)

	an, _ := strconv.ParseInt(info.Detail.Result.Value.AccountNumber, 10, 64)
	seq, _ := strconv.ParseInt(info.Detail.Result.Value.Sequence, 10, 64)

	free := info.Tokens.Available[0].String()
	locked := info.Tokens.Unbonding[0].String()
	frozen := info.Tokens.Delegated[0].String()

	account := &Account{
		Address:       info.BaseInfo.Address,
		PublicKey:     nil,
		AccountNumber: an,
		Sequence:      seq,
		Flags:         0,
		Balances: []AccountToken{
			{Symbol: "hst", Free: free, Locked: locked, Frozen: frozen},
		},
	}

	bic.Data["json"] = account
	bic.ServeJSON()
}

func (bic *BaseInfoController) getDeciamlRewardAmount(address string) decimal.Decimal {
	//var delegateReward model.DelegateRewards
	amount := bic.Base.Account.GetDelegateReward(address)
	decimalAmount, _ := decimal.NewFromString(amount)
	return decimalAmount
}
