package validator

import (
	"explorer/controllers"
	"explorer/model"
	"strings"

	"github.com/astaxie/beego"
)

type VaBaseInfoController struct {
	beego.Controller
	Base *controllers.BaseController
}
type MSGS struct {
	Code string        `json:"code"`
	Data ValidatorInfo `json:"data"`
	Msg  string        `json:"msg"`
}
type MSGError struct {
	Code string `json:"code"`
	Data error  `json:"data"`
	Msg  string `json:"msg"`
}
type ValidatorInfo struct {
	Jailed      bool   `json:"jailed"`
	Avater      string `json:"avater"`
	Rank        int    `json:"rank"`
	AKA         string `json:"aka"`
	Address     string `json:"address"`
	Validator   string `json:"validator"`
	WebSite     string `json:"web_site"`
	Commission  string `json:"commission"`
	Uptime      int    `json:"uptime"`
	VotingPower struct {
		Amount  float64 `json:"amount"`
		Percent float64 `json:"percent"`
	} `json:"voting_power"`
	HsnHeight       string                `json:"hsn_height"`
	Details         string                `json:"details"`
	TotalToken      float64               `json:"total_token"`
	SelfToken       float64               `json:"self_token"`
	OthersToken     float64               `json:"others_token"`
	MissedBlockList []model.MissBLockData `json:"missed_block_list"`
}

func (vbic *VaBaseInfoController) URLMapping() {
	vbic.Mapping("ValidatorBase", vbic.ValidatorBase)
}

// @Title 获取validator detail
// @Description 通过validator address 查询 validator detail详情
// @Success code 0
// @Failure code 1
// @router /validatorBase [get]
func (vbic *VaBaseInfoController) ValidatorBase() {
	vbic.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", vbic.Ctx.Request.Header.Get("Origin"))
	address := vbic.GetString("address")
	//var baseInfo validatorsDetail.ExtraValidatorInfo
	//var validatorInfo models.ValidatorInfo
	var respObj ValidatorInfo
	if address == "" || strings.Index(address, vbic.Base.Bech32PrefixValAddr) != 0 {
		var errMsg MSGError
		errMsg.Code = "1"
		errMsg.Data = nil
		errMsg.Msg = "Validator address is empty! Or error address!"
		vbic.Data["json"] = errMsg
		vbic.ServeJSON()
	}
	//objBase
	objBase := vbic.Base.ValidatorDetail.GetOne(address)
	//objBase := baseInfo.GetOne(address)
	objValidator := vbic.Base.Validator.GetOne(address)
	//objValidator := validatorInfo.GetOne(address)
	respObj.Jailed = objValidator.Jailed
	respObj.Avater = objValidator.Avater
	respObj.Rank = vbic.Base.Validator.GetValidatorRank(objValidator.VotingPower.Amount, objValidator.Jailed)
	//respObj.Rank = validatorInfo.GetValidatorRank(objValidator.VotingPower.Amount,objValidator.Jailed)
	respObj.AKA = objValidator.AKA
	respObj.Address = objBase.Address
	respObj.Validator = objBase.Validator
	respObj.WebSite = objBase.WebSite
	respObj.Commission = objValidator.Commission
	respObj.Uptime = objValidator.Uptime
	respObj.VotingPower = objValidator.VotingPower
	respObj.HsnHeight = objBase.HsnHeight
	respObj.Details = objBase.Details

	respObj.TotalToken = objBase.TotalToken
	respObj.SelfToken = objBase.SelfToken
	respObj.OthersToken = objBase.OthersToken
	respObj.MissedBlockList = objBase.MissedBlockList
	var msg MSGS
	msg.Data = respObj
	msg.Msg = "OK"
	msg.Code = "0"
	vbic.Data["json"] = msg
	vbic.ServeJSON()
}
