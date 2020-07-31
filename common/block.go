package common

import (
	"explorer/db"
	"explorer/model"

	"gopkg.in/mgo.v2/bson"
)

type block struct {
	db.MgoOperator
}

func NewBlock(m db.MgoOperator) Block {
	return &block{
		MgoOperator: m,
	}
}

func (b *block) SetBlock(block *model.BlockInfo) {
	conn := b.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	c := conn.C("block")
	err := c.Insert(&block)
	if err != nil {
		logger.Err(err).Interface(`BlockInfo`, block).Msg(`SetBlock`)
	}
}

func (b *block) GetAimHeightAndBlockHeight() (int, int) {
	var block model.BlockInfo
	var public model.Information

	conn := b.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	err := conn.C("public").Find(nil).Sort("-height").One(&public)
	if err != nil {
		logger.Err(err).Interface(`findby`, `height`).Msg(`GetAimHeightAndBlockHeight`)
	}
	err = conn.C("block").Find(nil).Sort("-intheight").One(&block)
	if err != nil {
		logger.Err(err).Interface(`findby`, `intheight`).Msg(`GetAimHeightAndBlockHeight`)
	}

	lastBlockHeight := block.IntHeight
	publicHeight := public.Height

	return lastBlockHeight, publicHeight
}

func (b *block) GetBlockListIfHasTx(height int) []model.BlocksHeights {
	var heights []model.BlocksHeights
	conn := b.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("block").Find(bson.M{"intheight": bson.M{"$gte": height}, "block.header.numtxs": bson.M{"$ne": "0"}}).All(&heights)
	return heights
}

func (b *block) GetLastBlockHeight() int {
	var block model.BlockInfo
	conn := b.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("block").Find(nil).Sort("-intheight").One(&block)

	return block.IntHeight
}
