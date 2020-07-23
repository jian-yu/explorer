package model

import (
	"time"
)

type Delegators struct {
	Address string    `json:"address"`
	Height  string    `json:"height"`
	Result  []DResult  `json:"result"`
	Time    time.Time `json:"time"`
}

type DResult struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Shares           string  `json:"shares"`
	Balance          Balance `json:"balance"`
}
type Balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
type DelegatorObj struct {
	Address          string    `json:"address"`
	Height           string    `json:"height"`
	DelegatorAddress string    `json:"delegator_address"`
	Shares           float64   `json:"shares"`
	Sign             int       `json:"sign"`
	Time             time.Time `json:"time"`
}
type ValidatorDelegatorNums struct {
	ValidatorAddress string `json:"validator_address"`
	DelegatorNums    int    `json:"delegator_nums"`
}

type Delegator struct {
	Result Result `json:"result"`
}

type DelegatorExtra struct {
	Height string `json:"height"`
	Result []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Shares           string `json:"shares"`
		Balance          struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
		Reward string `json:"reward"`
		Name   string `json:"name"`
	} `json:"result"`
}

type DelegatorValidatorReward struct {
	Height string `json:"height"`
	Result []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"result"`
}


type DelegateRewards struct {
	Result struct {
		Total []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"total"`
	} `json:"result"`
}


type Unbonding struct {
	Height string `json:"height"`
	Result []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Entries          []struct {
			CreationHeight string    `json:"creation_height"`
			CompletionTime time.Time `json:"completion_time"`
			InitialBalance string    `json:"initial_balance"`
			Balance        string    `json:"balance"`
		} `json:"entries"`
		Name string `json:"name"`
	} `json:"result"`
}


//func (d *DelegatorObj) GetInfo(address string, page int, size int) (*[]DelegatorObj, int) {
//
//	session := db.NewDBConn()
//	dbConn := session.DB(conf.NewConfig().DBName)
//	defer session.Close()
//	if size <= 0 {
//		size = 5
//		page = 0
//	}
//	var tempObj = make([]DelegatorObj, size)
//
//	//_ = dbConn.C("delegations").Find(
//	//	bson.M{
//	//		"address": address,
//	//		"time": bson.M{
//	//			"$gte": nowTime.Add(inInterval)}}).Sort("time").Skip(page * size).Limit(size).All(&tempObj)
//	//oneDayAgoDelegations, _ := dbConn.C("delegations").Find(
//	//	bson.M{
//	//		"address": address,
//	//		"time": bson.M{
//	//			"$gte": standardTime.Add(inInterval), "$lte": standardTime}}).Count()
//	//inOneIntervalDelegations, _ := dbConn.C("delegations").Find(
//	//	bson.M{
//	//		"address": address,
//	//		"time": bson.M{
//	//			"$gte": time.Now().Add(inInterval)}}).Count()
//	_ = dbConn.C("delegations").Find(
//		bson.M{
//			"address": address}).Sort("time").Sort("-shares").Skip(page * size).Limit(size).All(&tempObj)
//
//	inOneIntervalDelegations, _ := dbConn.C("delegations").Find(
//		bson.M{
//			"address": address}).Count()
//
//	return &tempObj, inOneIntervalDelegations
//}
//func (d *DelegatorObj) SetInfo(log zap.Logger) {
//	d.Time = time.Now()
//	session := db.NewDBConn()
//	defer session.Close()
//	dbConn := session.DB(conf.NewConfig().DBName)
//	// 检查数据，如果验证人地址和账户地址都一致，则使用upsert。
//	_, err := dbConn.C("delegations").Upsert(bson.D{{"address", d.Address}, {"delegatoraddress", d.DelegatorAddress}}, &d)
//	//err := dbConn.C("delegations").Insert(&d)
//	if err != nil {
//		log.Error("Insert delegations error!", zap.String("error", err.Error()))
//	} else {
//		log.Info("Insert delegations success")
//	}
//}
//
//func (vdn *ValidatorDelegatorNums) SetInfo(log zap.Logger) {
//
//	session := db.NewDBConn()
//	defer session.Close()
//	dbConn := session.DB(conf.NewConfig().DBName)
//	_, err := dbConn.C("validatorDelegatorNums").Upsert(bson.D{{"validatoraddress", vdn.ValidatorAddress}}, &vdn)
//	if err != nil {
//		log.Error("Update or Insert validatorDelegatorNums error!", zap.String("error", err.Error()))
//	} else {
//		log.Info("Update or Insert validatorDelegatorNums success")
//	}
//}
//func (vdn *ValidatorDelegatorNums) GetInfo(address string) int {
//
//	session := db.NewDBConn()
//	defer session.Close()
//	dbConn := session.DB(conf.NewConfig().DBName)
//	_ = dbConn.C("validatorDelegatorNums").Find(bson.M{"validatoraddress": address}).One(&vdn)
//	return vdn.DelegatorNums
//}
//
//func (d *DelegatorObj) DeleteInfo(sign int) {
//	session := db.NewDBConn()
//	defer session.Close()
//	dbConn := session.DB(conf.NewConfig().DBName)
//	_,_ = dbConn.C("delegations").RemoveAll((bson.M{"sign": bson.M{"$ne":sign}}))
//}

//
//func (d *DelegatorObj) DeleteAllInfo() {
//	session := db.NewDBConn()
//	defer session.Close()
//	dbConn := session.DB(conf.NewConfig().DBName)
//	_ = dbConn.C("delegations").DropCollection()
//}
