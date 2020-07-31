package asset

import (
	"explorer/handler"
	"explorer/model"
	"github.com/astaxie/beego"
	"strconv"
)

type Controller struct {
	beego.Controller
	*handler.ValidatorHandler
}

// @Title Get
// @Description get asset info
// @Success code 0
// @Failure code 1
// @router /assets [get]
func (a *Controller) Asset() {
	a.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", a.Ctx.Request.Header.Get("Origin"))
	assetName := a.GetString("asset", "")
	if assetName != "" {
		a.AssetByName()
		return
	}
	a.AssetPage()
}

// @Title Get
// @Description get asset info
// @Success code 0
// @Failure code 1
// @router /asset-holders [get]
func (a *Controller) AssetHolder() {
	a.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", a.Ctx.Request.Header.Get("Origin"))
	assetName := a.GetString("asset", "")
	if assetName == "" {
		a.Data["json"] = &model.Asset{}
		a.ServeJSON()
		return
	}
	page, _ := a.GetInt("page", 0)
	if page <= 1 {
		page = 0
	}
	rows, _ := a.GetInt("rows", 1)
	if rows <= 0 {
		rows = 1
	}

	val := a.ValidatorHandler.ValidatorByAsset(assetName)
	if val == nil {
		a.Data["json"] = &model.Asset{}
		a.ServeJSON()
		return
	}

	delegations := a.ValidatorHandler.Delegations(val.ValidatorAddress, page, rows)

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

	a.Data["json"] = assetHolders
	a.ServeJSON()
}

func (a *Controller) AssetPage() {
	page, _ := a.GetInt("page", 1)
	if page <= 0 {
		page = 1
	}
	rows, _ := a.GetInt("rows", 1)
	if rows <= 0 {
		rows = 1
	}

	var ret = &model.AssetInfo{}
	var tmpAssets []*model.Asset
	var total = 0

	info := a.ValidatorHandler.PublicInfo()

	validators := a.ValidatorHandler.Validators("active")
	if len(validators) == 0 {
		a.Data["json"] = ret
		a.ServeJSON()
		return
	}

	for i, val := range validators {
		valDetail := a.ValidatorHandler.ValidatorDetail(val.ValidatorAddress)
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
			Owner:           a.ValidatorHandler.ValidatorAccount(val.ValidatorAddress),
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

	if rows >= total {
		ret.AssetInfoList = tmpAssets
	} else {
		offset := (page - 1) * rows
		if offset+rows >= total {
			ret.AssetInfoList = tmpAssets[total-rows : total]
		} else {
			ret.AssetInfoList = tmpAssets[offset : offset+rows]
		}
	}
	ret.TotalNum = total

	a.Data["json"] = ret
	a.ServeJSON()
}

func (a *Controller) AssetByName() {
	assetName := a.GetString("asset", "")
	if assetName == "" {
		a.Data["json"] = &model.Asset{}
		a.ServeJSON()
		return
	}

	val := a.ValidatorHandler.ValidatorByAsset(assetName)
	if val == nil {
		a.Data["json"] = &model.Asset{}
		a.ServeJSON()
		return
	}

	info := a.ValidatorHandler.PublicInfo()
	valDetail := a.ValidatorHandler.ValidatorDetail(val.ValidatorAddress)
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
		Owner:           a.ValidatorHandler.ValidatorAccount(val.ValidatorAddress),
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

	a.Data["json"] = asset
	a.ServeJSON()
}
