package controllers

import (
	"explorer/model"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type DrawingDataController struct {
	beego.Controller
	Base *BaseController
}

type Drawing struct {
	Code string `json:"code"`
	Data Items  `json:"data"`
	Msg  string `json:"msg"`
}

type Items struct {
	Price []float64 `json:"price"`
	Token []int     `json:"token"`
}

// @Title Get
// @Description 首页小图
// @Success code 0
// @Failure code 1
//@router /
func (ddc *DrawingDataController) Get() {
	ddc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", ddc.Ctx.Request.Header.Get("Origin"))
	var public model.Information
	var respJson Drawing
	var price []float64
	var token [] int

	conn := ddc.Base.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	now := time.Now()
	for i := 0; i < 10; i++ {
		m, _ := time.ParseDuration("-1m")
		h1 := now.Add(time.Duration(i) * m)
		_ = conn.C("public").Find(bson.M{"time": bson.M{"$lt": h1}}).Sort("-height").One(&public)
		tempPrice, _ := strconv.ParseFloat(public.Price, 64)
		price = append(price, tempPrice)
	}
	for i := 0; i < 24; i++ {
		h, _ := time.ParseDuration("-1h")
		h1 := now.Add(time.Duration(i) * h)
		_ = conn.C("public").Find(bson.M{"time": bson.M{"$lt": h1}}).Sort("-height").One(&public)
		token = append(token, public.PledgeCoin)

	}

	respJson.Data.Price = price
	respJson.Data.Token = token
	respJson.Code = "0"
	respJson.Msg = "OK"
	ddc.Data["json"] = respJson
	ddc.ServeJSON()

}
