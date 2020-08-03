package controllers

import (
	"explorer/handler"
	"github.com/astaxie/beego"
	"github.com/shopspring/decimal"
)

type TokensController struct {
	beego.Controller
	*handler.TokensHandler
}

type Token struct {
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	OriginalSymbol string `json:"original_symbol"`
	TotalSupply    string `json:"total_supply"`
	Owner          string `json:"owner"`
	Mintable       bool   `json:"mintable"`
}

// @Title
// @Description
// @Success code 0
// @Failure code 1
// @router /tokens [get]
func (t *TokensController) Tokens() {
	t.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", t.Ctx.Request.Header.Get("Origin"))
	page, _ := t.GetInt("offset", 0)
	if page <= 0 {
		page = 1
	}
	rows, _ := t.GetInt("limit", 1)
	if rows <= 0 {
		rows = 1
	}

	var tokens []*Token

	validatorTokens := t.TokensHandler.GetValidatorsToken()
	total := len(validatorTokens)
	for _, val := range validatorTokens {
		total := decimal.NewFromFloat(val.Extra.TotalToken)
		token := &Token{
			Name:           val.Validator.AKA,
			Symbol:         val.Extra.Identity,
			OriginalSymbol: val.Extra.Identity,
			TotalSupply:    total.String(),
			Owner:          val.Owner,
			Mintable:       mintable(val.Validator.Status),
		}
		tokens = append(tokens, token)
	}

	if rows >= total {
		t.Data["json"] = tokens
	} else {
		offset := (page - 1) * rows
		if offset+rows >= total {
			t.Data["json"] = tokens[total-rows : total]
		} else {
			t.Data["json"] = tokens[offset : offset+rows]
		}
	}

	t.ServeJSON()
}

func mintable(status int) bool {
	if status == 2 {
		return true
	}
	return false
}
