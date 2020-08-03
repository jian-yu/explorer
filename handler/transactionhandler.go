package handler

import "explorer/model"

type TransactionHandler struct {
	base *BaseHandler
}

func NewTransactionHandler(base *BaseHandler) *TransactionHandler {
	return &TransactionHandler{base: base}
}

func (th *TransactionHandler) DelegatorTxs(address string, page int, size int) ([]*model.Txs, int) {
	return th.base.Transaction.GetDelegatorTxs(address, page, size)
}

func (th *TransactionHandler) Txs(before, after, limit int) ([]*model.Txs, int) {
	return th.base.Transaction.GetTxs(before, after, limit)
}

func (th *TransactionHandler) TxsByTypeAndTime(typo string, startTime int64, endTime int64, before, after, limit int) ([]*model.Txs, int) {
	return th.base.Transaction.GetTxsByTypeAndTime(typo, startTime, endTime, before, after, limit)
}

func (th *TransactionHandler) TxByHash(hash string) model.Txs {
	return th.base.Transaction.GetDetail(hash)
}
