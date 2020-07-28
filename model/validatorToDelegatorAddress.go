package model

type ValidatorToDelegatorAddress struct {
	ValidatorAddress string `json:"validator_address"`
	DelegatorAddress string `json:"delegator_address"`
}

type GenesisFile struct {
	Result struct {
		Genesis struct {
			AppState struct {
				Genutil struct {
					Gentxs []struct {
						Value struct {
							Msg []struct {
								Value struct {
									DelegatorAddress string `json:"delegator_address"`
									ValidatorAddress string `json:"validator_address"`
								} `json:"value"`
							} `json:"msg"`
						} `json:"value"`
					} `json:"gentxs"`
				} `json:"genutil"`
			} `json:"app_state"`
		} `json:"genesis"`
	} `json:"result"`
}

type CreatValidator struct {
	Txs []struct {
		Tx struct {
			Value struct {
				Msg []struct {
					Value struct {
						DelegatorAddress string `json:"delegator_address"`
						ValidatorAddress string `json:"validator_address"`
					} `json:"value"`
				} `json:"msg"`
			} `json:"value"`
		} `json:"tx"`
	} `json:"txs"`
}
