package controllers

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/astaxie/beego"
	"github.com/go-resty/resty/v2"
)

type TxsMsgs interface {
}
type TxDetailControllers struct {
	beego.Controller
	Base *BaseController
}
type ErrorTxInfoBlock struct {
	Code string `json:"code"`
	Data error  `json:"data"`
	Msg  string `json:"msg"`
}
type Msgs struct {
	Code string `json:"code"`
	Data TXD    `json:"data"`
	Msg  string `json:"msg"`
}

type TXD struct {
	Height string `json:"height"`
	Txhash string `json:"txhash"`
	RawLog string `json:"raw_log"`
	Logs   []struct {
		Success bool   `json:"success"`
		Log     string `json:"log"`
	} `json:"logs"`
	GasWanted string `json:"gas_wanted"`
	GasUsed   string `json:"gas_used"`
	Events    []struct {
		Type       string `json:"type"`
		Attributes []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"attributes"`
	} `json:"events"`
	Tx struct {
		Type  string `json:"type"`
		Value struct {
			Msg []interface{} `json:"msg"`
			Fee struct {
				Amount []struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"amount"`
				Gas string `json:"gas"`
			} `json:"fee"`
			Memo string `json:"memo"`
		} `json:"value"`
	} `json:"tx"`
	Timestamp time.Time `json:"timestamp"`
}

// @Title 获取tx detail
// @Description 通过hash 查询 tx详情
// @Success code 0
// @Failure code 1
// @router /
func (td *TxDetailControllers) Get() {
	var httpCli = resty.New()
	var msg Msgs
	var txd TXD

	td.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", td.Ctx.Request.Header.Get("Origin"))
	hash := td.GetString("hash")

	if hash == "" || td.Base.Transaction.CheckHash(hash) == 0 {
		var txd ErrorTxInfoBlock
		txd.Code = "1"
		txd.Msg = "Hash address is empty or Error!"
		txd.Data = nil
		td.Data["json"] = txd
	} else {
		url := td.Base.LcdURL + "/txs/" + hash
		rsp, err := httpCli.R().Get(url)
		if err != nil {
			log.Err(err).Interface(`url`, url).Msg(`Get`)
			td.Abort("500")
			return
		}

		err = json.Unmarshal(rsp.Body(), &txd)
		if err != nil {
			log.Err(err).Interface(`rsp`, rsp).Interface(`url`, url).Msg(`Get`)
			td.Abort("500")
			return
		}
		msg.Data = txd
		msg.Code = "0"
		msg.Msg = "OK"
		td.Data["json"] = msg
	}
	td.ServeJSON()
}
