package handler

import "explorer/model"

type BlockHandler struct {
	base *BaseHandler
}

func NewBlockHandler(base *BaseHandler) *BlockHandler {
	return &BlockHandler{base: base}
}

func (b *BlockHandler) GetBlocks(before, after, limit int) []*model.BlockInfo {
	return b.base.Block.GetBlocks(before, after, limit)
}

func (b *BlockHandler) GetLatestBlock() model.BlockInfo {
	return b.base.Block.GetLatestBlock()
}
