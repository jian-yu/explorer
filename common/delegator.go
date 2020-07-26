package common

import (
	"explorer/db"
	"explorer/model"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type delegator struct {
	db.MgoOperator
}

func NewDelegator(m db.MgoOperator) Delegator {
	return &delegator{
		MgoOperator: m,
	}
}

func (d *delegator) GetInfo(address string, page int, size int) (*[]model.DelegatorObj, int) {
	conn := d.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	if size <= 0 {
		size = 5
		page = 0
	}
	var tempObj = make([]model.DelegatorObj, size)

	_ = conn.C("delegations").Find(
		bson.M{
			"address": address}).Sort("time").Sort("-shares").Skip(page * size).Limit(size).All(&tempObj)

	inOneIntervalDelegations, _ := conn.C("delegations").Find(
		bson.M{
			"address": address}).Count()

	return &tempObj, inOneIntervalDelegations
}
func (d *delegator) SetInfo(delegator model.DelegatorObj) {
	conn := d.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	delegator.Time = time.Now()
	// 检查数据，如果验证人地址和账户地址都一致，则使用upsert。
	_, err := conn.C("delegations").Upsert(bson.D{{"address", delegator.Address}, {"delegatoraddress", delegator.DelegatorAddress}}, &delegator)
	//err := dbConn.C("delegations").Insert(&d)
	if err != nil {
		logger.Err(err).Interface(`DelegatorObj`, delegator).Msg(`SetDelegatorObj`)
	}
}
func (d *delegator) SetDelegatorCount(vDelegator model.ValidatorDelegatorNums) {
	conn := d.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_, err := conn.C("validatorDelegatorNums").Upsert(bson.D{{"validatoraddress", vDelegator.ValidatorAddress}}, &vDelegator)
	if err != nil {
		logger.Err(err).Interface(`ValidatorDelegatorNums`, vDelegator).Msg(`SetDelegatorCount`)
	}
}
func (d *delegator) GetDelegatorCount(address string) int {
	var vdn model.ValidatorDelegatorNums
	conn := d.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("validatorDelegatorNums").Find(bson.M{"validatoraddress": address}).One(&vdn)
	return vdn.DelegatorNums
}

func (d *delegator) DeleteInfo(sign int) {
	conn := d.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_, _ = conn.C("delegations").RemoveAll(bson.M{"sign": bson.M{"$ne": sign}})
}
