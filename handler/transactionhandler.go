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
