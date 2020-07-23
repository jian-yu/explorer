package common

import (
	"explorer/db"
	"explorer/model"
	"gopkg.in/mgo.v2/bson"
)

type validatorDetail struct {
	db.MgoOperator
}

func NewValidatorDetail(m db.MgoOperator) ValidatorDetail {
	return &validatorDetail{
		MgoOperator: m,
	}
}

func (v *validatorDetail) Set(info model.ExtraValidatorInfo) {
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	err := conn.C("detailValidatorBase").Insert(&info)
	if err != nil {
		logger.Err(err).Interface(`ExtraValidatorInfo`, info).Msg(`Set`)
	}
}

func (v *validatorDetail) Update(info model.ExtraValidatorInfo) {
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	err := conn.C("detailValidatorBase").Update(bson.M{"validator": info.Validator}, &info)
	if err != nil {
		logger.Err(err).Interface(`ExtraValidatorInfo`, info).Msg(`Update`)
	}
}

func (v *validatorDetail) Check(info model.ExtraValidatorInfo) int {
	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	count, _ := conn.C("detailValidatorBase").Find(bson.M{"validator": info.Validator}).Count()
	return count
}

func (v *validatorDetail) GetOne(address string) *model.ExtraValidatorInfo {
	var info model.ExtraValidatorInfo

	conn := v.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("detailValidatorBase").Find(bson.M{"validator": address}).One(&info)
	return &info
}
