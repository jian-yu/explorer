package common

import (
	"explorer/db"
	"explorer/model"
	"gopkg.in/mgo.v2/bson"
)

type proposer struct {
	db.MgoOperator
}

func NewProposer(m db.MgoOperator) Proposer {
	return &proposer{
		MgoOperator: m,
	}
}

func (p *proposer) SetInfo(vaak model.ValidatorAddressAndKey) {
	conn := p.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	c := conn.C("mapping")
	err := c.Insert(&vaak)
	if err == nil {
		logger.Err(err).Interface(`ValidatorAddressAndKey`, vaak).Msg(`SetInfo`)
	}
}

func (p *proposer) GetInfo(address string) string {
	var vaak model.ValidatorAddressAndKey
	conn := p.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("mapping").Find(bson.M{"operatoraddress": address}).One(&vaak)
	return vaak.ProposerHash
}

func (p *proposer) CheckValidator(pubkey string) (int, string) {
	var vaak model.ValidatorAddressAndKey
	var count int

	conn := p.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("mapping").Find(bson.M{"consensuspubkey": pubkey}).One(&vaak)
	if vaak.ProposerHash != "" {
		count = 1
	}
	return count, vaak.ProposerHash
}

func (p *proposer) GetValidator(hashAddr string) string {
	var vaak model.ValidatorAddressAndKey
	conn := p.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("mapping").Find(bson.M{"proposerhash": hashAddr}).One(&vaak)
	return vaak.OperatorAddress
}
