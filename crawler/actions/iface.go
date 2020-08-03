package actions

type Action interface {
	GetPublic()
	GetBlock()
	GetValidators()
	GetValidatorsSet()
	GetDelegations()
	GetDelegatorNums()
	GetTxs()
	GetTxs2()
	GetGenesis()
}
