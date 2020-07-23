package model

import (
	"time"
)

type ValidatorSet struct {
	BlockHeight string    `json:"block_height"`
	Height      int       `json:"height"`
	Time        time.Time `json:"time"`
	Validators  []struct {
		Address     string `json:"address"`
		PubKey      string `json:"pub_key"`
		VotingPower string `json:"voting_power"`
	} `json:"validators"`
}

// func (vs ValidatorsSet) SetInfo(log zap.Logger) {
// 	conn := v.MgoOperator.GetDBConn()
// 	defer conn.Session.Close()

// 	conn := session.DB(conf.NewConfig().DBName)
// 	c := dbConn.C("validatorsSet")
// 	err := c.Insert(&vs)
// 	if err != nil {
// 		log.Error("ValidatorsSet insert error", zap.String("error", err.Error()))
// 	} else {
// 		log.Info("ValidatorsSet insert success", zap.Int("height", vs.Height))
// 	}

// }

// // limit conf.go -> ValidatorsSetLimit
// // 获取最新的100个区块验证着集合 limit默认是  conf.go -> ValidatorsSetLimit
// func (vs ValidatorsSet) GetInfo(limit int) *[]ValidatorsSet {
// 	var ValidatorsSets []ValidatorsSet
// 	session := db.NewDBConn()
// 	defer session.Close()
// 	dbConn := session.DB(conf.NewConfig().DBName)
// 	dbConn.C("validatorsSet").Find(nil).Sort("-height").Limit(limit).All(&ValidatorsSets)
// 	return &ValidatorsSets
// }
