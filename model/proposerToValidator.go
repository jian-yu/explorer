package model

type ValidatorAddressAndKey struct {
	OperatorAddress string `json:"operator_address"`
	ConsensusPubkey string `json:"consensus_pubkey"`
	ProposerHash    string `json:"proposer_hash"`
}
