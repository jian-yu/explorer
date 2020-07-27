package actions

import (
	"encoding/json"
	"explorer/model"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type validatorsSets struct {
	Result model.ValidatorSet `json:"result"`
}

func (a *action) GetValidatorsSet() {
	/*
		获取当前高度， 开始同步区块验证人信息直至当前高度。
		再次获取高度，若高度相同sleep，若不同再次同步
	*/
	// Defer is useless in loops, does not apply defer or organizes part of the code into a function.
	var sets validatorsSets
	url := a.LcdURL + "/validatorsets/"

	for {
		aimHeight, ValidatorSetHeight := a.GetHeight()
		//如果高度差大于110，仅同步最新的100个高度地址
		if aimHeight-ValidatorSetHeight >= 110 {
			ValidatorSetHeight = aimHeight - 101
		}
		if aimHeight > ValidatorSetHeight {
			for aimHeight > ValidatorSetHeight {
				ValidatorSetHeight = ValidatorSetHeight + 1

				rsp, err := a.R().Get(url + strconv.Itoa(ValidatorSetHeight))
				if err != nil {
					log.Err(err).Interface(`url`, url).Msg(`GetValidatorsSet`)
					continue
				}

				err = json.Unmarshal(rsp.Body(), &sets)
				if err != nil {
					log.Err(err).Interface(`rsp`, rsp).Interface(`url`, url).Msg(`GetValidatorsSet`)
					intHeight, _ := strconv.Atoi(sets.Result.BlockHeight)
					sets.Result.Height = intHeight
					sets.Result.Time = time.Now()
					a.Validator.SetValidatorSet(sets.Result)
				}
				//time.Sleep(time.Millisecond * 2)
			}
		}
		time.Sleep(time.Second * 4)
	}
}

func (a *action) GetHeight() (int, int) {
	// get public's Height,ValidatorsSet's Height
	conn := a.MgoOperator.GetDBConn()
	defer conn.Session.Close()

	var tempValidatorsSet model.ValidatorSet
	_ = conn.C("validatorsSet").Find(nil).Sort("-height").One(&tempValidatorsSet)
	ValidatorsSetHeight := tempValidatorsSet.Height
	var tempPublic model.Information
	_ = conn.C("public").Find(nil).Sort("-height").One(&tempPublic)
	PublicHeight := tempPublic.Height
	return PublicHeight, ValidatorsSetHeight
}
