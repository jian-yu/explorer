package controllers

import (
	"explorer/handler"
	"github.com/astaxie/beego"
)

type ValidatorsController struct {
	beego.Controller
	*handler.ValidatorHandler
}

func (vc *ValidatorsController) URLMapping() {
	vc.Mapping("Validators", vc.Validators)
}

type MSGS struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
type errMsg struct {
	Code string `json:"code"`
	Data error  `json:"data"`
	Msg  string `json:"msg"`
}

// @Title 获取Validators List
// @Description get validators
// @Success code 0
// @Failure code 1
// @router /validators [get]
func (vc *ValidatorsController) Validators() {
	vc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", vc.Ctx.Request.Header.Get("Origin"))
	typo := vc.GetString("type", "")

	validatorList := vc.ValidatorHandler.Validators(typo)

	var msgs MSGS
	msgs.Code = "0"
	msgs.Msg = "OK"
	if len(validatorList) == 0 {
		var errMsg errMsg
		errMsg.Data = nil
		errMsg.Code = "1"
		errMsg.Msg = "No Record!"
		vc.Data["json"] = errMsg
		vc.ServeJSON()
		return
	}
	msgs.Data = validatorList
	vc.Data["json"] = msgs
	vc.ServeJSON()
}
