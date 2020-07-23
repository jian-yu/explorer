package common

import (
	"explorer/db"
	"explorer/model"
	"gopkg.in/mgo.v2/bson"
)

type transaction struct {
	db.MgoOperator
}

func NewTransaction(m db.MgoOperator) Transaction {
	return &transaction{
		MgoOperator: m,
	}
}

func (t *transaction) SetInfo(tx model.Txs) {
	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	c := conn.C("Txs")
	err := c.Insert(&tx)
	if err != nil {
		logger.Err(err).Interface(`Txs`, tx).Msg(`SetInfo`)
	}
}

func (t *transaction) GetInfo(head int, page int, size int) ([]model.Txs, int) {
	var TxsSet = make([]model.Txs, size)
	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	if page <= 0 {
		// default page
		page = 0
	}
	if size == 0 {
		size = 5
	}

	if head == 0 {
		var txs model.Txs
		_ = conn.C("Txs").Find(nil).Sort("-height").One(&txs)
		head = txs.Height
	}
	err := conn.C("Txs").Find(bson.M{"height": bson.M{
		"$lte": head}}).Sort("-height").Limit(size).Skip(page * size).All(&TxsSet)
	if err != nil {
		logger.Err(err).Interface(`params`,[]interface{}{head,page,size}).Msg(`GetInfo`)
	}
	totalTxsCount, _ := conn.C("Txs").Find(nil).Count()
	return TxsSet, totalTxsCount
}

func (t *transaction) GetDetail(txHash string) model.Txs {
	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	var temptxs model.Txs
	_ = conn.C("Txs").Find(bson.M{"txhash": txHash}).One(&temptxs)
	return temptxs
}

func (t *transaction) CheckHash(txHash string) int {
	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	var temptxs model.Txs
	_ = conn.C("Txs").Find(bson.M{"txhash": txHash}).One(&temptxs)
	if temptxs.TxHash == "" {
		return 0
	}
	return 1
}

func (t *transaction) GetValidatorsTransactions(address string) {
	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	var temptxs []model.Txs
	//type undelegate,delegate ...
	_ = conn.C("Txs").Find(bson.M{"txhash": address}).One(&temptxs)
}

func (t *transaction) GetPowerEventInfo(address string, page, size int) (*[]model.Txs, int) {
	var txsSet []model.Txs
	var query []bson.M
	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	if page < 0 {
		// default page
		page = 0
	}
	if size == 0 {
		size = 5
	}
	q1 := bson.M{"type": "unbonding"}
	q2 := bson.M{"type": "delegate"}
	q3 := bson.M{"type": "redelegate"}
	query = append(query, q1, q2, q3)
	count, _ := conn.C("Txs").Find(bson.M{"$or": query, "validatoraddress": bson.M{"$elemMatch": bson.M{"$eq": address}}, "result": true}).Count()
	_ = conn.C("Txs").Find(bson.M{"$or": query, "validatoraddress": bson.M{"$elemMatch": bson.M{"$eq": address}}, "result": true}).Sort("-height").Limit(size).Skip(page * size).All(&txsSet)
	return &txsSet, count
}

func (t *transaction) GetDelegatorTxs(address string, page, size int) (*[]model.Txs, int) {
	var txsSet []model.Txs
	if page < 0 {
		// default page
		page = 0
	}

	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	var query []bson.M
	q1 := bson.M{"delegatoraddress": bson.M{"$elemMatch": bson.M{"$eq": address}}}
	//from_address
	//out_puts_address
	//voter_address
	q2 := bson.M{"fromaddress": bson.M{"$elemMatch": bson.M{"$eq": address}}}
	q3 := bson.M{"outputsaddress": bson.M{"$elemMatch": bson.M{"$eq": address}}}
	q5 := bson.M{"inputsaddress": bson.M{"$elemMatch": bson.M{"$eq": address}}}
	q4 := bson.M{"voteraddress": bson.M{"$elemMatch": bson.M{"$eq": address}}}

	query = append(query, q1, q2, q3, q4, q5)
	count, _ := conn.C("Txs").Find(bson.M{"$or": query}).Count()
	_ = conn.C("Txs").Find(bson.M{"$or": query}).Sort("-height").Limit(size).Skip(page * size).All(&txsSet)

	return &txsSet, count
}
func (t *transaction) GetDelegatorCommissionTx(address string) *[]model.Txs {
	var txsSet []model.Txs

	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("Txs").Find(bson.M{"type": "commission", "delegatoraddress": bson.M{"$elemMatch": bson.M{"$eq": address}}}).All(&txsSet)
	return &txsSet
}
func (t *transaction) GetDelegatorRewardTx(address string) *[]model.Txs {
	var txsSet []model.Txs
	var query []bson.M

	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	q1 := bson.M{"type": "commission"}
	q2 := bson.M{"type": "reward"}
	query = append(query, q1, q2)
	_ = conn.C("Txs").Find(bson.M{"$or": query, "delegatoraddress": bson.M{"$elemMatch": bson.M{"$eq": address}}}).All(&txsSet)
	return &txsSet
}
func (t *transaction) GetSpecifiedHeight(head int, page int, size int) ([]model.Txs, int) {
	var TxsSet = make([]model.Txs, size)

	if page <= 0 {
		// default page
		page = 0
	}
	if size == 0 {
		size = 5
	}

	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("Txs").Find(bson.M{"height": head}).Sort("-height").Limit(size).Skip(page * size).All(&TxsSet)
	totalTxsCount, _ := conn.C("Txs").Find(bson.M{"height": head}).Count()
	return TxsSet, totalTxsCount
}

func (t *transaction) GetTxHeight(tx model.Txs) int {
	conn := t.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	_ = conn.C("Txs").Find(nil).Sort("-height").One(&tx)
	return tx.Height
}
