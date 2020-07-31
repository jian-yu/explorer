package model

type (
	Asset struct {
		CreateTime      interface{} `json:"createTime"`
		UpdateTime      interface{} `json:"updateTime"`
		ID              int         `json:"id"`
		Asset           string      `json:"asset"`
		MappedAsset     string      `json:"mappedAsset"`
		Name            string      `json:"name"`
		AssetImg        string      `json:"assetImg"`
		Supply          float64     `json:"supply"`
		Price           float64     `json:"price"`
		QuoteUnit       string      `json:"quoteUnit"`
		ChangeRange     float64     `json:"changeRange"`
		Owner           string      `json:"owner"`
		Mintable        int         `json:"mintable"`
		Visible         interface{} `json:"visible"`
		Description     string      `json:"description"`
		AssetCreateTime interface{} `json:"assetCreateTime"`
		Transactions    int         `json:"transactions"`
		Holders         int         `json:"holders"`
		OfficialSiteURL string      `json:"officialSiteUrl"`
		ContactEmail    string      `json:"contactEmail"`
		MediaList       []struct {
			MediaName string `json:"mediaName"`
			MediaURL  string `json:"mediaUrl"`
			MediaImg  string `json:"mediaImg"`
		} `json:"mediaList"`
	}

	AssetInfo struct {
		TotalNum      int      `json:"totalNum"`
		AssetInfoList []*Asset `json:"assetInfoList"`
	}

	AssetHolders struct {
		TotalNum       int            `json:"totalNum"`
		AddressHolders []*AddressHolders `json:"addressHolders"`
	}

	AddressHolders struct {
		Address    string      `json:"address"`
		Quantity   float64     `json:"quantity"`
		Percentage float64     `json:"percentage"`
		Tag        interface{} `json:"tag"`
	}
)
