package account

import (
	"explorer/controllers"
	"explorer/handler"
	"strings"

	"github.com/astaxie/beego"
)

type DelegatorTxController struct {
	beego.Controller
	Base *controllers.BaseController
	*handler.TransactionHandler
}
type txInfo struct {
	Height int     `json:"height"`
	Hash   string  `json:"hash"`
	Types  string  `json:"types"`
	Result bool    `json:"result"`
	Amount float64 `json:"amount"`
	Fee    float64 `json:"fee"`
	Nums   int     `json:"nums"`
	Time   string  `json:"time"`
}
type TxBlocks struct {
	Code  string   `json:"code"`
	Data  []txInfo `json:"data"`
	Total int      `json:"total"`
	Msg   string   `json:"msg"`
}

type AccountTxs struct {
	TxNums  int          `json:"txNums"`
	TxArray []*AccountTx `json:"txArray"`
}

type AccountTx struct {
	TxHash        string  `json:"txHash"`
	BlockHeight   int64   `json:"blockHeight"`
	TxType        string  `json:"txType"`
	TimeStamp     int64   `json:"timeStamp"`
	FromAddr      string  `json:"fromAddr"`
	ToAddr        string  `json:"toAddr"`
	Value         float64 `json:"value"`
	TxAsset       string  `json:"txAsset"`
	TxQuoteAsset  string  `json:"txQuoteAsset"`
	TxFee         float64 `json:"txFee"`
	TxAge         int64   `json:"txAge"`
	OrderID       string  `json:"orderId"`
	Data          string  `json:"data,omitempty"`
	Code          int64   `json:"code"`
	Log           string  `json:"log"`
	ConfirmBlocks int64   `json:"confirmBlocks"`
	Memo          string  `json:"memo"`
	Source        int64   `json:"source"`
	HasChildren   int64   `json:"hasChildren"`
}

func (dtc *DelegatorTxController) URLMapping() {
	dtc.Mapping("DelegatorTxs", dtc.DelegatorTxs)
}

// @Title
// @Description
// @Success code 0
// @Failure code 1
// @router /delegatorTxs [get]
func (dtc *DelegatorTxController) DelegatorTxs() {
	dtc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", dtc.Ctx.Request.Header.Get("Origin"))
	address := dtc.GetString("address")
	page, _ := dtc.GetInt("page", 0)
	size, _ := dtc.GetInt("size", 5)
	if address == "" || strings.Index(address, dtc.Base.Bech32PrefixAccAddr) != 0 || strings.Index(address, dtc.Base.Bech32PrefixValAddr) == 0 {
		var errMsg TxBlocks
		errMsg.Data = nil
		errMsg.Code = "1"
		errMsg.Msg = "Delegator address is empty Or Error address!"
		dtc.Data["json"] = errMsg
	} else {
		//var txList models.Txs
		var txsSet = make([]txInfo, size)
		var respJSON TxBlocks
		list, count := dtc.Base.Transaction.GetDelegatorTxs(address, page, size)
		//list, count := txList.GetDelegatorTxs(address, page, size)
		for i, item := range list {
			txsSet[i].Height = item.Height
			txsSet[i].Hash = item.TxHash
			txsSet[i].Fee = item.Fee
			txsSet[i].Result = item.Result
			txsSet[i].Time = item.TxTime
			txsSet[i].Types = item.Type
			txsSet[i].Nums = item.Plus

			if item.Type == "reward" {
				txsSet[i].Amount = getRewardAmount(item.WithDrawRewardAmout)
			} else {
				txsSet[i].Amount = getAmount(item.Amount)
			}

		}
		respJSON.Code = "0"
		respJSON.Msg = "OK"
		respJSON.Total = count
		respJSON.Data = txsSet
		dtc.Data["json"] = respJSON

	}

	dtc.ServeJSON()
}

// @Title
// @Description
// @Success code 0
// @Failure code 1
// @router /txs [get]
func (dtc *DelegatorTxController) Txs() {
	dtc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", dtc.Ctx.Request.Header.Get("Origin"))
	address := dtc.GetString("address", "")
	if address == "" {
		dtc.Abort("400")
		return
	}

	page, _ := dtc.GetInt("page", 0)
	if page <= 1 {
		page = 0
	}
	rows, _ := dtc.GetInt("rows", 1)
	if rows <= 0 {
		rows = 1
	}

	txs, count := dtc.TransactionHandler.DelegatorTxs(address, page, rows)

	var accTxs []*AccountTx
	for _, tx := range txs {
		fromAddrs := strings.Join(tx.FromAddress, ",")
		toAddrs := strings.Join(tx.ToAddress, ",")

		accTx := &AccountTx{
			TxHash:        tx.TxHash,
			BlockHeight:   int64(tx.Height),
			TxType:        tx.Type,
			TimeStamp:     tx.Time.Unix(),
			FromAddr:      fromAddrs,
			ToAddr:        toAddrs,
			Value:         getAmount(tx.Amount),
			TxAsset:       "",
			TxQuoteAsset:  "",
			TxFee:         tx.Fee,
			TxAge:         0,
			OrderID:       "",
			Data:          tx.Data,
			Code:          0,
			Log:           tx.RawLog,
			ConfirmBlocks: int64(tx.Height),
			Memo:          tx.Memo,
			Source:        0,
			HasChildren:   0,
		}
		accTxs = append(accTxs, accTx)
	}

	accTxArr := &AccountTxs{
		TxNums:  count,
		TxArray: accTxs,
	}

	dtc.Data["json"] = accTxArr
	dtc.ServeJSON()
}

func getAmount(amounts []float64) float64 {
	var totalAmount float64
	if len(amounts) <= 0 {
		return 0.0
	}
	for i := 0; i < len(amounts); i++ {
		totalAmount = totalAmount + amounts[i]
	}

	return totalAmount
}

func getRewardAmount(amounts []float64) float64 {
	if len(amounts) == 1 {
		return amounts[0]
	}
	return 0.0
}
