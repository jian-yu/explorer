package account

import (
	"explorer/controllers"
	"explorer/model"
	"strings"

	"github.com/astaxie/beego"
)

type DeleatorsController struct {
	beego.Controller
	Base *controllers.BaseController
}
type DelegatorerrMsg struct {
	Data error  `json:"data"`
	Msg  string `json:"msg"`
	Code string `json:"code"`
}
type DelegatorMsg struct {
	Data  model.DelegatorExtra `json:"data"`
	Msg   string               `json:"msg"`
	Total int                  `json:"total"`
	Size  int                  `json:"size"`
	Code  string               `json:"code"`
}

// @Title
// @Description
// @Success code 0
// @Failure code 1
//@router /
func (dc *DeleatorsController) Get() {
	dc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", dc.Ctx.Request.Header.Get("Origin"))
	address := dc.GetString("address")
	page, _ := dc.GetInt("page", 0)
	size, _ := dc.GetInt("size", 5)

	if address == "" || strings.Index(address, dc.Base.Bech32PrefixAccAddr) != 0 || strings.Index(address, dc.Base.Bech32PrefixValAddr) == 0 {
		var msg baseInfoerrMsg
		msg.Data = nil
		msg.Msg = "Delegator address is empty Or Error address!"
		msg.Code = "1"
		dc.Data["json"] = msg
	} else {
		//获取验证人账户信息和获取提款地址
		var msg DelegatorMsg
		//var delegators model.DelegatorExtra
		infos := dc.Base.Account.GetDelegator(address)
		//infos := delegators.GetInfo(address)
		msg.Code = "0"
		msg.Size = size
		msg.Total = len(infos.Result)
		msg.Msg = "OK"
		msg.Data = *infos
		//msg.Data = *infos

		/*分页*/
		msg.Data = *delegatorPagination(page, size, msg.Total, infos)
		dc.Data["json"] = msg
	}
	dc.ServeJSON()

}
func delegatorPagination(page, size, totalSize int, infos *model.DelegatorExtra) *model.DelegatorExtra {
	//account.Unbonding
	var tempVar model.DelegatorExtra

	if page*size <= 0 {
		//return first page
		if totalSize < size {
			for i := 0; i < totalSize; i++ {
				tempVar.Result = append(tempVar.Result, infos.Result[i])
			}
		} else {
			for i := 0; i < size; i++ {
				tempVar.Result = append(tempVar.Result, infos.Result[i])
			}
		}
		return &tempVar
	}
	if page*size > 0 && (page+1)*size <= totalSize {
		for i := page * size; i < (page+1)*size; i++ {
			tempVar.Result = append(tempVar.Result, infos.Result[i])
		}
		return &tempVar
	}
	if (page+1)*size > totalSize {
		//return last page
		if totalSize-(page)*size > 0 {
			for i := page * size; i < totalSize; i++ {
				tempVar.Result = append(tempVar.Result, infos.Result[i])
			}
		}
	}
	return &tempVar
}
