package account

import (
	"explorer/controllers"
	"explorer/model"
	"strings"

	"github.com/astaxie/beego"
)

type UnbondingsController struct {
	beego.Controller
	Base *controllers.BaseController
}
type UnbondingserrMsg struct {
	Data error  `json:"data"`
	Msg  string `json:"msg"`
	Code string `json:"code"`
}
type UnbondingsMsg struct {
	Data  model.Unbonding `json:"data"`
	Msg   string          `json:"msg"`
	Total int             `json:"total"`
	Size  int             `json:"size"`
	Code  string          `json:"code"`
}

// @Title
// @Description
// @Success code 0
// @Failure code 1
//@router /
func (uc *UnbondingsController) Get() {
	uc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", uc.Ctx.Request.Header.Get("Origin"))
	address := uc.GetString("address")
	page, _ := uc.GetInt("page", 0)
	size, _ := uc.GetInt("size", 5)

	if address == "" || strings.Index(address, uc.Base.Bech32PrefixAccAddr) != 0 || strings.Index(address, uc.Base.Bech32PrefixValAddr) == 0 {
		var msg baseInfoerrMsg
		msg.Data = nil
		msg.Msg = "Delegator address is empty Or Error address!"
		msg.Code = "1"
		uc.Data["json"] = msg
	} else {
		var msg UnbondingsMsg
		//var unbonding accountDetail.Unbonding
		infos := uc.Base.Account.GetUnbonding(address)
		//infos := unbonding.GetInfo(address)
		msg.Code = "0"
		msg.Code = "0"
		msg.Size = size
		msg.Total = len(infos.Result)
		msg.Msg = "OK"
		//msg.Data = *infos
		/*分页*/
		msg.Data = *pagination(page, size, msg.Total, infos)
		uc.Data["json"] = msg
	}

	uc.ServeJSON()
}
func pagination(page, size, totalSize int, infos *model.Unbonding) *model.Unbonding {
	//account.Unbonding
	var tempVar model.Unbonding

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
