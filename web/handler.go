package web

import (
	"encoding/json"
	"explorer/handler"
	"explorer/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HTTPHandler struct {
	validatorHandler   *handler.ValidatorHandler
	accountHandler     *handler.AccountHandler
	delegationHandler  *handler.DelegationsHandler
	tokensHandler      *handler.TokensHandler
	transactionHandler *handler.TransactionHandler
	blockHandler       *handler.BlockHandler
}

func NewHTTPHandler(
	engine *gin.Engine,
	val *handler.ValidatorHandler,
	acc *handler.AccountHandler,
	del *handler.DelegationsHandler,
	token *handler.TokensHandler,
	trans *handler.TransactionHandler,
	block *handler.BlockHandler,
) *HTTPHandler {
	wh := &HTTPHandler{
		val,
		acc,
		del,
		token,
		trans,
		block,
	}
	router := engine.Group("/api/v1")
	router.GET("/asset", wh.handleAssetName)
	router.GET("/assets", wh.handleAssets)
	router.GET("/asset-holders", wh.handleAssetHolder)
	router.GET("/accounts/:address", wh.handleAccount)
	router.GET("/account/txs", wh.handleAccountTxs)
	router.GET("/stake/validators", wh.handleValidators)
	router.GET("/stake/validators/:address", wh.handleValidator)
	router.GET("/tokens", wh.handleTokens)
	router.GET("/txs", wh.handleTxs)
	router.GET("/tx", wh.handleTxByHash)
	router.GET("/blocks", wh.handleBlocks)
	router.GET("/blocks/latest", wh.handleLastBlock)
	return wh
}

func (wh *HTTPHandler) handleAssetName(ctx *gin.Context) {
	name := ctx.Query("asset")
	if name == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewWebError(-2400, "not found asset param"))
		return
	}

	val := wh.validatorHandler.ValidatorByAsset(name)
	if val == nil {
		ctx.JSON(http.StatusOK, &model.Asset{})
		return
	}

	info := wh.validatorHandler.PublicInfo()
	valDetail := wh.validatorHandler.ValidatorDetail(val.ValidatorAddress)
	price, _ := strconv.ParseFloat(info.Price, 64)
	asset := &model.Asset{
		CreateTime:      nil,
		UpdateTime:      nil,
		ID:              0,
		Asset:           val.AKA,
		MappedAsset:     valDetail.Identity,
		Name:            val.AKA,
		AssetImg:        val.Avater,
		Supply:          val.VotingPower.Amount,
		Price:           price,
		QuoteUnit:       "hst",
		ChangeRange:     0,
		Owner:           wh.validatorHandler.ValidatorAccount(val.ValidatorAddress),
		Mintable:        val.Status,
		Visible:         nil,
		Description:     valDetail.Details,
		AssetCreateTime: nil,
		Transactions:    0,
		Holders:         0,
		OfficialSiteURL: valDetail.WebSite,
		ContactEmail:    "",
		MediaList:       nil,
	}

	ctx.JSON(http.StatusOK, asset)
}

func (wh *HTTPHandler) handleAssets(ctx *gin.Context) {
	var p struct {
		Page int `form:"page"`
		Rows int `form:"rows"`
	}

	if err := ctx.BindQuery(&p); err != nil {
		log.Err(err).Msg(`handleAssets`)
		return
	}

	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Rows <= 0 {
		p.Rows = 1
	}

	var ret = &model.AssetInfo{}
	var tmpAssets []*model.Asset
	var total = 0

	info := wh.validatorHandler.PublicInfo()

	validators := wh.validatorHandler.Validators("active")
	if len(validators) == 0 {
		ctx.JSON(http.StatusOK, ret)
		return
	}

	for i, val := range validators {
		valDetail := wh.validatorHandler.ValidatorDetail(val.ValidatorAddress)
		price, _ := strconv.ParseFloat(info.Price, 64)
		asset := &model.Asset{
			CreateTime:      nil,
			UpdateTime:      nil,
			ID:              i,
			Asset:           val.AKA,
			MappedAsset:     valDetail.Identity,
			Name:            val.AKA,
			AssetImg:        val.Avater,
			Supply:          val.VotingPower.Amount,
			Price:           price,
			QuoteUnit:       "hst",
			ChangeRange:     0,
			Owner:           wh.validatorHandler.ValidatorAccount(val.ValidatorAddress),
			Mintable:        val.Status,
			Visible:         nil,
			Description:     valDetail.Details,
			AssetCreateTime: nil,
			Transactions:    0,
			Holders:         0,
			OfficialSiteURL: valDetail.WebSite,
			ContactEmail:    "",
			MediaList:       nil,
		}
		tmpAssets = append(tmpAssets, asset)
		total++
	}

	if p.Rows >= total {
		ret.AssetInfoList = tmpAssets
	} else {
		offset := (p.Page - 1) * p.Rows
		if offset+p.Rows >= total {
			ret.AssetInfoList = tmpAssets[total-p.Rows : total]
		} else {
			ret.AssetInfoList = tmpAssets[offset : offset+p.Rows]
		}
	}
	ret.TotalNum = total

	ctx.JSON(http.StatusOK, ret)
}

func (wh *HTTPHandler) handleAssetHolder(ctx *gin.Context) {
	var p struct {
		Page      int    `form:"page"`
		Rows      int    `form:"rows"`
		AssetName string `form:"asset"`
	}

	if err := ctx.BindQuery(&p); err != nil {
		log.Err(err).Msg(`handleAssetHolder`)
		return
	}

	if p.Page <= 0 {
		p.Page = 0
	}
	if p.Rows <= 0 {
		p.Rows = 1
	}

	val := wh.validatorHandler.ValidatorByAsset(p.AssetName)
	if val == nil {
		ctx.JSON(http.StatusOK, &model.AddressHolders{})
		return
	}

	delegations := wh.validatorHandler.Delegations(val.ValidatorAddress, p.Page, p.Rows)

	var assetHolders model.AssetHolders
	var addrHolders []*model.AddressHolders
	for _, item := range delegations.Delegations {
		addrHolder := &model.AddressHolders{
			Address:    item.Address,
			Quantity:   item.Amount,
			Percentage: item.AmountPercentage,
		}
		addrHolders = append(addrHolders, addrHolder)
	}
	assetHolders.TotalNum = delegations.TotalDelegations
	assetHolders.AddressHolders = addrHolders

	ctx.JSON(http.StatusOK, assetHolders)
}

func (wh *HTTPHandler) handleAccount(ctx *gin.Context) {
	var p struct {
		Address string `uri:"address" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&p); err != nil {
		log.Err(err).Msg(`handleAccount`)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewWebError(-24001, "not found account address"))
		return
	}

	info := wh.accountHandler.Account(p.Address)

	an, _ := strconv.ParseInt(info.Detail.Result.Value.AccountNumber, 10, 64)
	seq, _ := strconv.ParseInt(info.Detail.Result.Value.Sequence, 10, 64)

	free := info.Tokens.Available[0].String()
	locked := info.Tokens.Unbonding[0].String()
	frozen := info.Tokens.Delegated[0].String()

	account := &Account{
		Address:       info.BaseInfo.Address,
		PublicKey:     nil,
		AccountNumber: an,
		Sequence:      seq,
		Flags:         0,
		Balances: []AccountToken{
			{Symbol: "hst", Free: free, Locked: locked, Frozen: frozen},
		},
	}

	ctx.JSON(http.StatusOK, account)
}

func (wh *HTTPHandler) handleValidators(ctx *gin.Context) {
	var validators []*Validator

	stakeVals := wh.validatorHandler.StakeValidators()

	for _, val := range stakeVals {
		unbondingHeight, _ := strconv.ParseInt(val.Validator.Result[0].UnbondingHeight, 10, 64)
		power, _ := strconv.ParseInt(val.Power, 10, 64)
		pubKey, _ := json.Marshal(val.Validator.Result[0].ConsensusPubkey)

		validator := &Validator{
			AccountAddress:   val.Owner,
			OperatorAddress:  val.Validator.Result[0].OperatorAddress,
			ConsensusPubKey:  pubKey,
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
		validators = append(validators, validator)
	}

	ctx.JSON(http.StatusOK, validators)
}

func (wh *HTTPHandler) handleTokens(ctx *gin.Context) {
	var p struct {
		Offset int `form:"offset"`
		Limit  int `form:"limit"`
	}

	if err := ctx.BindQuery(&p); err != nil {
		log.Err(err).Msg(`handleTokens`)
		return
	}

	if p.Offset <= 0 {
		p.Offset = 1
	}
	if p.Limit <= 0 {
		p.Limit = 1
	}

	var tokens []*Token
	var result interface{}

	var mintable = func(status int) bool {
		if status == 2 {
			return true
		}
		return false
	}

	validatorTokens := wh.tokensHandler.GetValidatorsToken()
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

	if p.Limit >= total {
		result = tokens
	} else {
		offset := (p.Offset - 1) * p.Limit
		if offset+p.Limit >= total {
			result = tokens[total-p.Limit : total]
		} else {
			result = tokens[offset : offset+p.Limit]
		}
	}

	ctx.JSON(http.StatusOK, result)
}

func (wh *HTTPHandler) handleValidator(ctx *gin.Context) {
	var p struct {
		Address string `uri:"address"`
	}

	if err := ctx.BindUri(&p); err != nil {
		log.Err(err).Msg("handleValidator")
		return
	}

	val := wh.validatorHandler.ValidatorByAddress(p.Address)
	if val == nil {
		ctx.JSON(http.StatusOK, &Validator{})
		return
	}

	unbondingHeight, _ := strconv.ParseInt(val.Validator.Result[0].UnbondingHeight, 10, 64)
	power, _ := strconv.ParseInt(val.Power, 10, 64)
	pubKey, _ := json.Marshal(val.Validator.Result[0].ConsensusPubkey)
	validator := &Validator{
		AccountAddress:   val.Owner,
		OperatorAddress:  val.Validator.Result[0].OperatorAddress,
		ConsensusPubKey:  pubKey,
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

	ctx.JSON(http.StatusOK, validator)
}

func (wh *HTTPHandler) handleTxs(ctx *gin.Context) {
	var p Param
	var txs []*model.Txs
	var total = 0

	if err := ctx.BindQuery(&p); err != nil {
		log.Err(err).Msg(`handleTxs`)
		return
	}

	if p.Limit > 100 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewWebError(-24002, "limit must equal to or smaller than 100"))
		return
	}

	if p.Asset != "" {
		wh.handleAssetName(ctx)
		return
	}

	if p.Type != "" {
		txs, total = wh.transactionHandler.TxsByTypeAndTime(p.Type, p.StartTime, p.EndTime, p.Before, p.After, p.Limit)
	} else {
		txs, total = wh.transactionHandler.Txs(p.Before, p.After, p.Limit)
	}

	var txDatas []*TxData

	for i, tx := range txs {
		var msgs []Message
		var rawMsgs []map[string]interface{}
		_ = json.Unmarshal(tx.Message, &rawMsgs)
		for _, item := range rawMsgs {
			value, _ := json.Marshal(item["value"])
			msg := Message{
				Type:  item["type"].(string),
				Value: value,
			}
			msgs = append(msgs, msg)
		}
		var signs []Signature
		var rawSigns []map[string]interface{}
		_ = json.Unmarshal(tx.Sign, &rawSigns)
		for _, item := range rawSigns {

			sign := Signature{
				Pubkey:        "",
				Address:       "",
				Sequence:      "",
				Signature:     item["signature"].(string),
				AccountNumber: "",
			}
			signs = append(signs, sign)
		}

		txData := &TxData{
			ID:         int32(p.After + i),
			Height:     int64(tx.Height),
			Result:     tx.Result,
			TxHash:     tx.TxHash,
			Messages:   msgs,
			Signatures: signs,
			Memo:       tx.Memo,
			Code:       0,
			Timestamp:  time.Unix(tx.TxTime, 0),
		}
		txDatas = append(txDatas, txData)
	}

	type ret struct {
		Data  interface{} `json:"data"`
		Total int         `json:"total"`
	}

	ctx.JSON(http.StatusOK, &ret{
		Data:  txDatas,
		Total: total,
	})
}

func (wh *HTTPHandler) handleTxByHash(ctx *gin.Context) {
	var p struct {
		Hash string `form:"hash"`
	}

	if err := ctx.BindQuery(&p); err != nil {
		log.Err(err).Msg(`handleTxByHash`)
		return
	}

	tx := wh.transactionHandler.TxByHash(p.Hash)

	var msgs []Message
	var rawMsgs []map[string]interface{}
	_ = json.Unmarshal(tx.Message, &rawMsgs)
	for _, item := range rawMsgs {
		value, _ := json.Marshal(item["value"])
		msg := Message{
			Type:  item["type"].(string),
			Value: value,
		}
		msgs = append(msgs, msg)
	}
	var signs []Signature
	var rawSigns []map[string]interface{}
	_ = json.Unmarshal(tx.Sign, &rawSigns)
	for _, item := range rawSigns {

		sign := Signature{
			Pubkey:        "",
			Address:       "",
			Sequence:      "",
			Signature:     item["signature"].(string),
			AccountNumber: "",
		}
		signs = append(signs, sign)
	}

	txData := &TxData{
		ID:         0,
		Height:     int64(tx.Height),
		Result:     tx.Result,
		TxHash:     tx.TxHash,
		Messages:   msgs,
		Signatures: signs,
		Memo:       tx.Memo,
		Code:       0,
		Timestamp:  time.Unix(tx.TxTime, 0),
	}

	ctx.JSON(http.StatusOK, txData)
}

func (wh *HTTPHandler) handleAccountTxs(ctx *gin.Context) {
	var p struct {
		Address string `form:"address"`
		Page    int    `form:"page"`
		Rows    int    `form:"rows"`
	}

	if err := ctx.BindQuery(&p); err != nil {
		log.Err(err).Msg(`handleAccountTxs`)
		return
	}

	txs, count := wh.transactionHandler.DelegatorTxs(p.Address, p.Page, p.Rows)

	var accTxs []*AccountTx
	for _, tx := range txs {
		fromAddrs := strings.Join(tx.FromAddress, ",")
		toAddrs := strings.Join(tx.ToAddress, ",")

		accTx := &AccountTx{
			TxHash:        tx.TxHash,
			BlockHeight:   int64(tx.Height),
			TxType:        tx.Type,
			TimeStamp:     tx.TxTime,
			FromAddr:      fromAddrs,
			ToAddr:        toAddrs,
			Value:         getAmount(tx.Amount),
			TxAsset:       "",
			TxQuoteAsset:  "",
			TxFee:         tx.Fee,
			TxAge:         0,
			OrderID:       "",
			Data:          tx.Data,
			Code:          0,
			Log:           tx.RawLog,
			ConfirmBlocks: int64(tx.Height),
			Memo:          tx.Memo,
			Source:        0,
			HasChildren:   0,
		}
		accTxs = append(accTxs, accTx)
	}

	accTxArr := &AccountTxs{
		TxNums:  count,
		TxArray: accTxs,
	}

	ctx.JSON(http.StatusOK, accTxArr)
}

func (wh *HTTPHandler) handleBlocks(ctx *gin.Context) {
	var p Param

	if err := ctx.BindQuery(&p); err != nil {
		log.Err(err).Msg(`handleBlocks`)
		return
	}

	blocks := wh.blockHandler.GetBlocks(p.Before, p.After, p.Limit)

	var blockList []*BlockData
	for _, block := range blocks {
		numTxs, _ := strconv.ParseInt(block.BlockMeta.Header.NumTxs, 10, 64)
		TotalTxs, _ := strconv.ParseInt(block.BlockMeta.Header.TotalTxs, 10, 64)
		tt, _ := time.ParseInLocation(time.RFC3339Nano, block.BlockMeta.Header.Time, time.Local)
		//var txs []Txs

		blockData := &BlockData{
			Height:        int64(block.IntHeight),
			Proposer:      block.BlockMeta.Header.ProposerAddress,
			Moniker:       "",
			BlockHash:     block.BlockMeta.BlockID.Hash,
			ParentHash:    block.BlockMeta.Header.LastBlockID.Hash,
			NumPrecommits: int64(len(block.Block.LastCommit.Precommits)),
			NumTxs:        numTxs,
			TotalTxs:      TotalTxs,
			Txs:           nil,
			Timestamp:     tt,
		}
		blockList = append(blockList, blockData)
	}
	ctx.JSON(http.StatusOK, blockList)
}

func (wh *HTTPHandler) handleLastBlock(ctx *gin.Context) {
	block := wh.blockHandler.GetLatestBlock()
	numTxs, _ := strconv.ParseInt(block.BlockMeta.Header.NumTxs, 10, 64)
	TotalTxs, _ := strconv.ParseInt(block.BlockMeta.Header.TotalTxs, 10, 64)
	tt, _ := time.ParseInLocation(time.RFC3339Nano, block.BlockMeta.Header.Time, time.Local)

	blockData := &BlockData{
		Height:        int64(block.IntHeight),
		Proposer:      block.BlockMeta.Header.ProposerAddress,
		Moniker:       "",
		BlockHash:     block.BlockMeta.BlockID.Hash,
		ParentHash:    block.BlockMeta.Header.LastBlockID.Hash,
		NumPrecommits: int64(len(block.Block.LastCommit.Precommits)),
		NumTxs:        numTxs,
		TotalTxs:      TotalTxs,
		Txs:           nil,
		Timestamp:     tt,
	}

	ctx.JSON(http.StatusOK, blockData)
}

func getAmount(amounts []float64) float64 {
	var totalAmount float64
	if len(amounts) <= 0 {
		return 0.0
	}
	for i := 0; i < len(amounts); i++ {
		totalAmount = totalAmount + amounts[i]
	}

	return totalAmount
}
