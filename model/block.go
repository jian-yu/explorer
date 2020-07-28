package model

type BlockInfo struct {
	BlockMeta BlockMeta `json:"block_meta"`
	Block     Block     `json:"block"`
	IntHeight int       `json:"int_height"`
}

type BlockMeta struct {
	BlockID BlockID `json:"block_id"`
}
type Block struct {
	Header Header `json:"header"`
	//Data       Data       `json:"data"`
	//Evidence   Evidence   `json:"evidence"`
	//LastCommit LastCommit `json:"last_commit"`
}

type BlockID struct {
	Hash string `json:"hash"`
	//Parts Parts  `json:"parts"`
}
type Parts struct {
	Total string `json:"total"`
	Hash  string `json:"hash"`
}

type Header struct {
	//Version            Version `json:"version"`
	//ChainId            string  `json:"chain_id"`
	Height      string  `json:"height"`
	Time        string  `json:"time"` //p
	NumTxs      string  `json:"num_txs"`
	TotalTxs    string  `json:"total_txs"`
	LastBlockID BlockID `json:"last_block_id"`
	//LastCommitHash     string  `json:"last_commit_hash"`
	//DataHash           string  `json:"data_hash"`
	ValidatorsHash string `json:"validators_hash"`
	//NextValidatorsHash string  `json:"next_validators_hash"`
	//ConsensusHash      string  `json:"consensus_hash"`
	//AppHash            string  `json:"app_hash"`
	//LastResultHash     string  `json:"last_result_hash"`
	//EvidenceHash       string  `json:"evidence_hash"`
	ProposerAddress string `json:"proposer_address"`
}
type Version struct {
	Block string `json:"block"`
	App   string `json:"app"`
}

type Data struct {
	Txs []string `json:"txs"`
}

type Evidence struct {
	Evidence string `json:"evidence"`
}
type LastCommit struct {
	BlockID    BlockID          `json:"block_id"`
	Precommits []PreCommitsList `json:"precommits"`
}

type PreCommitsList struct {
	Type             int     `json:"type"`
	Height           string  `json:"height"`
	Round            string  `json:"round"`
	BlockID          BlockID `json:"block_id"`
	Timestamp        string  `json:"timestamp"`
	ValidatorAddress string  `json:"validator_address"`
	ValidatorIndex   string  `json:"validator_index"`
	Signature        string  `json:"signature"`
}
type BlocksHeights struct {
	IntHeight int `json:"int_height"`
}
