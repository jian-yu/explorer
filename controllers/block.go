package controllers

import (
	"explorer/model"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type BlockController struct {
	beego.Controller
	Base *BaseController
}

type Blocks struct {
	Code  string        `json:"code"`
	Data  []BlockSimple `json:"data"`
	Total int           `json:"total"`
	Msg   string        `json:"msg"`
}

type BlockSimple struct {
	Height    int    `json:"height"`
	BlockHash string `json:"block_hash"`
	Proposer  string `json:"proposer"`
	AKA       string `json:"aka"`
	Txs       string `json:"txs"`
	Time      string `json:"time"`
}

// @Title 获取区快
// @Description 默认获取after head的20个区块详细信息
// @Success code 0
// @Failure code 1
//@router /
func (bc *BlockController) Get() {
	bc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", bc.Ctx.Request.Header.Get("Origin"))
	head, _ := bc.GetInt("head", 0)
	page, _ := bc.GetInt("page", 0)
	size, _ := bc.GetInt("size", 0)

	conn := bc.Base.GetDBConn()
	defer conn.Session.Close()

	var public model.Information
	_ = conn.C("public").Find(nil).Sort("-height").One(&public)
	if page == 0 {
		// default page
		page = 0
	}
	if size == 0 {
		// default last SIZE
		size = 5
	}
	if head == 0 {
		// default last height
		head = public.Height

	}

	var blocks = make([]model.BlockInfo, size)
	var blockInfoSimples = make([]BlockSimple, size)
	var respJson Blocks
	_ = conn.C("block").Find(bson.M{"intheight": bson.M{
		"$lte": head,}}).Sort("-intheight").Limit(size).Skip(size * page).All(&blocks)
	for i, item := range blocks {
		blockInfoSimples[i].Height = item.IntHeight
		blockInfoSimples[i].BlockHash = item.BlockMeta.BlockId.Hash
		blockInfoSimples[i].Proposer = bc.getProposerAddress(item.Block.Header.ProposerAddress)
		blockInfoSimples[i].AKA = bc.getAKAName(blockInfoSimples[i].Proposer)
		blockInfoSimples[i].Txs = item.Block.Header.NumTxs
		blockInfoSimples[i].Time = item.Block.Header.Time
	}
	respJson.Total = public.Height
	respJson.Data = blockInfoSimples
	respJson.Code = "0"
	respJson.Msg = "OK"
	bc.Data["json"] = respJson
	bc.ServeJSON()
}
func (bc *BlockController) getProposerAddress(hashAddress string) string {
	address := bc.Base.Proposer.GetValidator(hashAddress)
	return address
}

func (bc *BlockController) getAKAName(proposerAddress string) string {
	validator := bc.Base.Validator.GetOne(proposerAddress)
	return validator.AKA
}
