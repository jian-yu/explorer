package web

import (
	"explorer/handler"
	"explorer/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type webHandler struct {
	validatorHandler  *handler.ValidatorHandler
	accountHandler    *handler.AccountHandler
	delegationHandler *handler.DelegationsHandler
	tokenHandler      *handler.TokensHandler
}

func NewWebHandler(
	engine *gin.Engine,
	val *handler.ValidatorHandler,
	acc *handler.AccountHandler,
	del *handler.DelegationsHandler,
	token *handler.TokensHandler,
) *webHandler {
	wh := &webHandler{
		val,
		acc,
		del,
		token,
	}
	router := engine.Group("/api/v1")
	router.GET("/asset", wh.handleAssetName)
	router.GET("/assets", wh.handleAssets)
	router.GET("/asset-holders", wh.handleAssetHolder)

	return wh
}

func (wh *webHandler) handleAssetName(ctx *gin.Context) {
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

func (wh *webHandler) handleAssets(ctx *gin.Context) {
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

func (wh *webHandler) handleAssetHolder(ctx *gin.Context) {
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
