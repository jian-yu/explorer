package actions

import (
	"encoding/json"
	"explorer/common"
	"explorer/crawler"
	"explorer/db"
	"explorer/model"
	"gopkg.in/resty.v1"
	"strconv"
	"time"
)

func GetBlock() {
	var httpClient = resty.New()
	var block model.BlockInfo
	mgoStore := db.NewMongoStore()
	blockOperator := common.NewBlock(mgoStore)

	for {
		lastBlockHeight, publicHeight := blockOperator.GetAimHeightAndBlockHeight()
		//check the height difference again
		if publicHeight > lastBlockHeight {
			for publicHeight > lastBlockHeight {
				lastBlockHeight = lastBlockHeight + 1

				url := crawler.LcdURL + "/blocks/" + strconv.Itoa(lastBlockHeight)
				rsp, err := httpClient.R().Get(url)
				if err != nil {
					lastBlockHeight = lastBlockHeight - 1
					time.Sleep(time.Second * 2)
					continue
				}
				err = json.Unmarshal(rsp.Body(), &block)
				if err != nil {

				}
				blockOperator.SetBlock(block)
			}
		}
		time.Sleep(time.Second * 1)
	}
}
