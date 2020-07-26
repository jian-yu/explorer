package common

import (
	"explorer/db"
	"explorer/model"
	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"
)

type validator struct {
	db.MgoOperator
}

//NewValidator vali
func NewValidator(m db.MgoOperator) Validator {
	return &validator{
		MgoOperator: m,
	}
}

func (v *validator) SetInfo(info model.ValidatorInfo) {
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	c := conn.C("validators")
	_, err := c.Upsert(bson.D{{"validatoraddress", info.ValidatorAddress}}, &info)
	if err != nil {
		log.Err(err).Interface(`validator address`, info).Msg(`SetInfo`)
	}
}

func (v *validator) GetInfo() *[]model.ValidatorInfo {
	var list []model.ValidatorInfo
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("validators").Find(nil).Sort("-votingpower.amount").All(&list)
	return &list
}

func (v *validator) DeleteAllInfo() {
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("validators").DropCollection()
}

func (v *validator) GetOne(address string) *model.ValidatorInfo {
	var info model.ValidatorInfo
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("validators").Find(bson.M{"validatoraddress": address}).One(&info)
	return &info
}

func (v *validator) GetValidatorRank(amount float64, jailed bool) int {
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	if !jailed {
		rank, _ := conn.C("validators").Find(bson.M{
			"votingpower.amount": bson.M{"$gte": amount},
			"jailed":             jailed}).Count()
		return rank
	}
	rank, _ := conn.C("validators").Find(bson.M{
		"votingpower.amount": bson.M{"$gte": amount},
		"jailed":             jailed}).Count()
	return rank
}

func (v *validator) SetValidatorSet(vs model.ValidatorSet) {
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	c := conn.C("validatorsSet")
	err := c.Insert(&vs)
	if err != nil {
		log.Err(err).Interface(`ValidatorSet`, vs).Msg(`SetValidatorSet`)
	}
}

func (v *validator) GetValidatorSet(limit int) *[]model.ValidatorSet {
	var ValidatorsSets []model.ValidatorSet
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	err := conn.C("validatorsSet").Find(nil).Sort("-height").Limit(limit).All(&ValidatorsSets)
	if err != nil {
		log.Err(err).Msg(`GetValidatorSet`)
		return nil
	}
	return &ValidatorsSets
}

func (v *validator) SetValidatorToDelegatorAddr(v2d model.ValidatorToDelegatorAddress) {
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	c := conn.C("mapping")
	err := c.Insert(&v2d)
	if err != nil {
		log.Err(err).Msg(`SetValidatorToDelegatorAddr`)
	}
}

func (v *validator) Check(address string) (int, string) {
	var tempValue model.ValidatorToDelegatorAddress
	var count = 0

	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	conn.C("mapping").Find(bson.M{"validatoraddress": address}).One(&tempValue)
	if tempValue.DelegatorAddress != "" {
		count = 1
	}
	return count, tempValue.DelegatorAddress
}

func (v *validator) CheckDelegatorAddress(address string) (string, string) {
	var tempValue model.ValidatorToDelegatorAddress
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("mapping").Find(bson.M{"delegatoraddress": address}).One(&tempValue)
	return tempValue.ValidatorAddress, tempValue.DelegatorAddress
}
