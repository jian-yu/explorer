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

// @Title Get
// @Description public Item
// @Success code 0
// @Failure code 1
//@router / [get]
func (pb *PublicController) Get() {
	pb.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", pb.Ctx.Request.Header.Get("Origin"))
	var public model.Information
	var respJson Public

	conn := pb.Base.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("public").Find(nil).Sort("-height").One(&public)
	respJson.Data = public
	respJson.Code = "0"
	respJson.Msg = "OK"
	pb.Data["json"] = respJson
	pb.ServeJSON()

}
