package account

import (
	"explorer/controllers"
	"strings"

	"github.com/astaxie/beego"
	"github.com/shopspring/decimal"
)

type KindsRewardController struct {
	beego.Controller
	Base *controllers.BaseController
}

type KindsRewardErrMsg struct {
	Data error  `json:"data"`
	Msg  string `json:"msg"`
	Code string `json:"code"`
}
type KindsRewardMsg struct {
	Data Kinds  `json:"data"`
	Msg  string `json:"msg"`
	Code string `json:"code"`
}
type Kinds struct {
	// (amount percentage)
	Available   []decimal.Decimal `json:"available"`
	Delegated   []decimal.Decimal `json:"delegated"`
	Unbonding   []decimal.Decimal `json:"unbonding"`
	Reward      []decimal.Decimal `json:"reward"`
	Commission  []decimal.Decimal `json:"commission"`
	TotalAmount []decimal.Decimal `json:"total_amount"`
}

func (krc *KindsRewardController) URLMapping() {
	krc.Mapping("DelegatorAllKindsReward", krc.DelegatorAllKindsReward)
}

// @Title
// @Description
// @Success code 0
// @Failure code 1
// @router /delegatorAllKindsReward [get]
func (krc *KindsRewardController) DelegatorAllKindsReward() {
	krc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", krc.Ctx.Request.Header.Get("Origin"))
	address := krc.GetString("address", "")

	// get delegator's txs .计算各类奖励之和.(delegate, unbonding, Reward, CommissionReward,
	//Available=account's token sub (delegate+ unbonding+ Reward+ CommissionReward))
	//Reward http://172.38.8.89:1317/distribution/delegators/hsn190zh9q92xs8y7s4c0tpc784vass5lalt7f3n0h/rewards
	//delegate SUM http://172.38.8.89:1317/staking/delegators/hsn1502lgkad0tnc2szdww0whpxs30szz03lj6n06q/delegations
	//Unbonding Sum http://172.38.8.89:1317/staking/delegators/hsn1aqwurs5jfu5z0z3k99tljt9csausdqrcaewjwv/unbonding_delegations
	//Commission

	if address == "" || strings.Index(address, krc.Base.Bech32PrefixAccAddr) != 0 || strings.Index(address, krc.Base.Bech32PrefixValAddr) == 0 {
		var errMsg KindsRewardErrMsg
		errMsg.Data = nil
		errMsg.Code = "1"
		errMsg.Msg = "Delegator address is empty Or Error address!"
		krc.Data["json"] = errMsg
		krc.ServeJSON()
	}

	msg := krc.GetAllKindsAmount(address)
	// 检查该地址是否是验证人的账户，如果是普通账户它的commission 等于nil
	//var valToDelMapping model.ValidatorToDelegatorAddress
	validatorAddress, _ := krc.Base.Validator.CheckDelegatorAddress(address)
	//validatorAddress, _ := valToDelMapping.CheckDelegatorAddress(address)
	if validatorAddress == "" {
		msg.Data.Commission = nil
	}
	msg.Code = "0"
	msg.Msg = "OK"
	krc.Data["json"] = msg
	krc.ServeJSON()
}
func (krc *KindsRewardController) getAvailableAmount(address string) []decimal.Decimal {
	var availableAmount []decimal.Decimal
	_, strAmount := krc.Base.Account.GetInfo(address)
	//_, strAmount := account.GetInfo(address)
	decimalAmount, _ := decimal.NewFromString(strAmount)
	//, decimalPercentage
	availableAmount = append(availableAmount, decimalAmount)
	return availableAmount
}

// from LCD interface.
//The amount is so large that there is a negative number in the results.
func (krc *KindsRewardController) getReward(address string) []decimal.Decimal {
	//var delegateReward model.DelegateRewards
	var decimalAmount decimal.Decimal
	var rewards []decimal.Decimal
	amount := krc.Base.Account.GetDelegateReward(address)
	//amount := delegateReward.GetDelegateReward(address)
	decimalAmount, _ = decimal.NewFromString(amount)
	rewards = append(rewards, decimalAmount)
	return rewards
}
func (krc *KindsRewardController) getTotalDelegateAmount(address string) []decimal.Decimal {
	//var delegators model.Delegators
	var amount decimal.Decimal
	var delegate []decimal.Decimal
	infos := krc.Base.Account.GetDelegator(address)
	//infos := delegators.GetInfo(address)
	for _, item := range infos.Result {
		decimalAmount, _ := decimal.NewFromString(item.Balance.Amount)
		amount = amount.Add(decimalAmount)
	}
	delegate = append(delegate, amount)
	return delegate
}
func (krc *KindsRewardController) getTotalUnbondingAmount(address string) []decimal.Decimal {
	//var unbonding accountDetail.Unbonding
	var amount decimal.Decimal
	var unbond []decimal.Decimal
	infos := krc.Base.Account.GetUnbonding(address)
	//infos := unbonding.GetInfo(address)
	for _, item := range infos.Result {
		for _, entrie := range item.Entries {
			decimalAmount, _ := decimal.NewFromString(entrie.Balance)
			amount = amount.Add(decimalAmount)
		}
	}
	unbond = append(unbond, amount)
	return unbond
}
func (krc *KindsRewardController) getTotalCommissionAmount(address string) []decimal.Decimal {
	//var txs models.Txs
	var commission []decimal.Decimal
	var decimalCommissionAmount decimal.Decimal
	commissionTxs := krc.Base.Transaction.GetDelegatorCommissionTx(address)
	//commissionTxs := txs.GetDelegatorCommissionTx(address)
	if len(*commissionTxs) == 0 {
		decimalCommissionAmount, _ = decimal.NewFromString("0.0")
	} else {
		for _, item := range *commissionTxs {
			for index, delegator := range item.DelegatorAddress {
				if delegator == address {
					decimalWithDrawCommissionAmount := decimal.NewFromFloat(item.WithDrawCommissionAmout[index])
					decimalCommissionAmount = decimalCommissionAmount.Add(decimalWithDrawCommissionAmount)
				}
			}
		}
	}

	commission = append(commission, decimalCommissionAmount)
	return commission
}
func (krc *KindsRewardController) GetAllKindsAmount(address string) *KindsRewardMsg {
	var msg KindsRewardMsg
	msg.Data.Reward = krc.getReward(address)
	msg.Data.Available = krc.getAvailableAmount(address)
	msg.Data.Delegated = krc.getTotalDelegateAmount(address)
	msg.Data.Unbonding = krc.getTotalUnbondingAmount(address)
	msg.Data.Commission = krc.getTotalCommissionAmount(address)
	var totalAmount []decimal.Decimal
	var percentage decimal.Decimal
	decimalAmount := (msg.Data.Available[0]).Add(msg.Data.Commission[0]).Add(msg.Data.Unbonding[0]).Add(msg.Data.Delegated[0]).Add(msg.Data.Reward[0])
	percentage = msg.Data.Available[0].Div(decimalAmount)
	msg.Data.Available = append(msg.Data.Available, percentage)
	percentage = msg.Data.Commission[0].Div(decimalAmount)
	msg.Data.Commission = append(msg.Data.Commission, percentage)
	percentage = msg.Data.Unbonding[0].Div(decimalAmount)
	msg.Data.Unbonding = append(msg.Data.Unbonding, percentage)
	percentage = msg.Data.Delegated[0].Div(decimalAmount)
	msg.Data.Delegated = append(msg.Data.Delegated, percentage)
	percentage = msg.Data.Reward[0].Div(decimalAmount)
	msg.Data.Reward = append(msg.Data.Reward, percentage)
	msg.Data.TotalAmount = append(totalAmount, decimalAmount)
	return &msg
}
