package controllers

import (
	"encoding/json"
	"explorer/handler"
	"github.com/astaxie/beego"
	"strconv"
	"time"
)

type ValidatorsController struct {
	beego.Controller
	*handler.ValidatorHandler
}

type (
	Validator struct {
		AccountAddress     string          `json:"account_address" sql:",notnull, unique"`
		OperatorAddress    string          `json:"operator_address" sql:",notnull, unique"`
		ConsensusPubKey    json.RawMessage `json:"consensus_pubkey" sql:",notnull, unique"`
		ConsensusAddress   string          `json:"consensus_address" sql:",notnull, unique"`
		Jailed             bool            `json:"jailed"`
		Status             string          `json:"status"`
		Tokens             string          `json:"tokens"`
		Power              int64           `json:"power"`
		DelegatorShares    string          `json:"delegator_shares"`
		Description        Description     `json:"description"`
		BondHeight         int64           `json:"bond_height"`
		BondIntraTxCounter int64           `json:"bond_intra_tx_counter"`
		UnbondingHeight    int64           `json:"unbonding_height"`
		UnbondingTime      time.Time       `json:"unbonding_time"`
		Commission         Commission      `json:"commission"`
	}

	// Description wraps validator's description information
	Description struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}

	// Commission wraps validator's commission information
	Commission struct {
		Rate          string    `json:"rate"`
		MaxRate       string    `json:"max_rate"`
		MaxChangeRate string    `json:"max_change_rate"`
		UpdateTime    time.Time `json:"update_time"`
	}
)

func (vc *ValidatorsController) URLMapping() {
	vc.Mapping("Validators", vc.Validators)
}

type MSGS struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
type errMsg struct {
	Code string `json:"code"`
	Data error  `json:"data"`
	Msg  string `json:"msg"`
}

// @Title 获取Validators List
// @Description get validators
// @Success code 0
// @Failure code 1
// @router /validators [get]
func (vc *ValidatorsController) Validators() {
	vc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", vc.Ctx.Request.Header.Get("Origin"))
	typo := vc.GetString("type", "")

	validatorList := vc.ValidatorHandler.Validators(typo)

	var msgs MSGS
	msgs.Code = "0"
	msgs.Msg = "OK"
	if len(validatorList) == 0 {
		var errMsg errMsg
		errMsg.Data = nil
		errMsg.Code = "1"
		errMsg.Msg = "No Record!"
		vc.Data["json"] = errMsg
		vc.ServeJSON()
		return
	}
	msgs.Data = validatorList
	vc.Data["json"] = msgs
	vc.ServeJSON()
}

// @Title 获取Validators List
// @Description get validators
// @Success code 0
// @Failure code 1
// @router /stake/validators [get]
func (vc *ValidatorsController) StakeValidators() {
	vc.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", vc.Ctx.Request.Header.Get("Origin"))

	var validators []*Validator

	stakeVals := vc.ValidatorHandler.StakeValidators()

	for _, val := range stakeVals {
		unbondingHeight, _ := strconv.ParseInt(val.Validator.Result[0].UnbondingHeight, 10, 64)
		power, _ := strconv.ParseInt(val.Power, 10, 64)
		pubkey, _ := json.Marshal(val.Validator.Result[0].ConsensusPubkey)

		val := &Validator{
			AccountAddress:   val.Owner,
			OperatorAddress:  val.Validator.Result[0].OperatorAddress,
			ConsensusPubKey:  pubkey,
			ConsensusAddress: "",
			Jailed:           val.Validator.Result[0].Jailed,
			Status:           strconv.Itoa(val.Validator.Result[0].Status),
			Tokens:           val.Validator.Result[0].Tokens,
			Power:            power,
			DelegatorShares:  val.Validator.Result[0].DelegatorShares,
			Description: Description{
				Moniker:  val.Validator.Result[0].Description.Moniker,
				Identity: val.Validator.Result[0].Description.Identity,
				Website:  val.Validator.Result[0].Description.Website,
				Details:  val.Validator.Result[0].Description.Details,
			},
			BondHeight:         0,
			BondIntraTxCounter: 0,
			UnbondingHeight:    unbondingHeight,
			UnbondingTime:      val.Validator.Result[0].UnbondingTime,
			Commission: Commission{
				Rate:          val.Validator.Result[0].Commission.CommissionRates.Rate,
				MaxRate:       val.Validator.Result[0].Commission.CommissionRates.MaxRate,
				MaxChangeRate: val.Validator.Result[0].Commission.CommissionRates.MaxChangeRate,
				UpdateTime:    val.Validator.Result[0].Commission.UpdateTime,
			},
		}
		validators = append(validators, val)
	}
	vc.Data["json"] = validators
	vc.ServeJSON()
}
