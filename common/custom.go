package common

import (
	"explorer/db"
	"explorer/model"
	"strconv"

	"github.com/shopspring/decimal"
)

type custom struct {
	db.MgoOperator
}

func NewCustom(m db.MgoOperator) Custom {
	return &custom{
		MgoOperator: m,
	}
}

func (c *custom) SetInfo(info model.Information) {
	conn := c.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	collection := conn.C("public")
	err := collection.Insert(&info)
	if err != nil {
		logger.Err(err).Interface(`Information`, info).Msg(`SetInfo`)
	}
}

func (c *custom) GetInfo() model.Information {
	var info model.Information
	conn := c.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	err := conn.C("public").Find(nil).Sort("-height").One(&info)
	if err != nil {
		logger.Err(err).Msg(`GetInfo`)
	}
	return info
}

func (c *custom) GetAllPledgenTokens() decimal.Decimal {
	var info model.Information
	conn := c.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("public").Find(nil).Sort("-height").One(&info)
	tokens := strconv.Itoa(info.PledgeCoin)
	decimalTotalHsn, _ := decimal.NewFromString(tokens)
	return decimalTotalHsn
}
