package controllers

import (
	"explorer/model"

	"github.com/astaxie/beego"
)

type PublicController struct {
	beego.Controller
	Base *BaseController
}

type Public struct {
	Code string            `json:"code"`
	Data model.Information `json:"data"`
	Msg  string            `json:"msg"`
}

func (pb *PublicController) URLMapping() {
	pb.Mapping("Public", pb.Public)
}

// @Title Get
// @Description public Item
// @Success code 0
// @Failure code 1
// @router /public [get]
func (pb *PublicController) Public() {
	pb.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", pb.Ctx.Request.Header.Get("Origin"))
	var public model.Information
	var respJSON Public

	conn := pb.Base.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("public").Find(nil).Sort("-height").One(&public)
	respJSON.Data = public
	respJSON.Code = "0"
	respJSON.Msg = "OK"
	pb.Data["json"] = respJSON
	pb.ServeJSON()

}
