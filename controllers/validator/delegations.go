package validator

import (
	"explorer/controllers"
	"explorer/model"
	"strings"

	"github.com/astaxie/beego"
)

type DelegationsController struct {
	beego.Controller
	Base *controllers.BaseController
}
type MsgErr struct {
	Code string `json:"code"`
	Data error  `json:"data"`
	Msg  string `json:"msg"`
}
type Msgs struct {
	Code string        `json:"code"`
	Data DelegationMsg `json:"data"`
	Msg  string        `json:"msg"`
}
type DelegationMsg struct {
	TotalDelegations     int           `json:"total_delegation"`
	OneDayAgoDelegations int           `json:"one_day_ago_delegations"`
	Delegations          []Delegations `json:"delegations"`
}
type Delegations struct {
	Address          string  `json:"address"`
	Amount           float64 `json:"amount"`
	AmountPercentage float64 `json:"share"`
}

func (dc *DelegationsController) URLMapping() {
	dc.Mapping("ValidatorDelegations", dc.ValidatorDelegations)
}

// @Title Get
// @Description get delegations
// @Success code 0
// @Failure code 1
// @router /validatorDelegations [get]
func (dc *DelegationsController) ValidatorDelegations() {
	dc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", dc.Ctx.Request.Header.Get("Origin"))
	address := dc.GetString("address")
	page, _ := dc.GetInt("page", 0)
	size, _ := dc.GetInt("size", 0)
	if size == 0 {
		size = 5
	}
	if address == "" || strings.Index(address, dc.Base.Bech32PrefixValAddr) != 0 {
		var errorMessage MsgErr
		errorMessage.Code = "1"
		errorMessage.Data = nil
		errorMessage.Msg = "Validator address is empty! Or error address!"
		dc.Data["json"] = errorMessage
		dc.ServeJSON()
	}
	var msg Msgs
	var respJSON DelegationMsg
	var respJSONDelegations []Delegations
	var respJSONDelegation Delegations
	//var delegations model.DelegatorObj
	//var baseInfo model.ExtraValidatorInfo
	//var validatorDelegationNums model.ValidatorDelegatorNums

	validatorBaseInfo := dc.Base.ValidatorDetail.GetOne(address)
	//validatorBaseInfo := baseInfo.GetOne(address)

	items, totalDelegations := dc.Base.Delegator.GetInfo(address, page, size)
	//items, totalDelegations := delegations.GetInfo(address, page, size)
	oneDayAgoDelegations := dc.Base.Delegator.GetDelegatorCount(address)
	//oneDayAgoDelegations := validatorDelegationNums.GetInfo(address)
	for _, item := range *items {
		respJSONDelegation.Amount = getShares(items, item.DelegatorAddress)
		respJSONDelegation.Address = item.DelegatorAddress
		respJSONDelegation.AmountPercentage = getPercentage(respJSONDelegation.Amount, validatorBaseInfo.TotalToken)
		respJSONDelegations = append(respJSONDelegations, respJSONDelegation)
	}
	respJSON.TotalDelegations = totalDelegations
	respJSON.OneDayAgoDelegations = totalDelegations - oneDayAgoDelegations
	respJSON.Delegations = respJSONDelegations
	msg.Data = respJSON
	msg.Code = "0"
	msg.Msg = "OK"
	dc.Data["json"] = msg
	dc.ServeJSON()
}

func getShares(items *[]model.DelegatorObj, address string) float64 {
	var amount float64
	for _, item := range *items {
		if item.DelegatorAddress == address {
			share := item.Shares
			amount = amount + share
		}
	}
	return amount
}
func getPercentage(amout float64, totalToken float64) float64 {

	return amout / totalToken

}
