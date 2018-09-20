package models

import (
	"encoding/hex"
)

type ChainInfo struct {
	Height            uint64 `json:"height"`
	TxCount           uint64 `json:"tx_count"`
	CurrentBlockHash  []byte `json:"current_block_hash"`
	PreviousBlockHash []byte `json:"previous_block_hash"`
}

func (c *ChainInfo) GetHeight() uint64 {
	if c != nil {
		return c.Height
	}
	return 0
}

func (c *ChainInfo) GetTxCount() uint64 {
	if c != nil {
		return c.TxCount
	}
	return 0
}

func (c *ChainInfo) GetCurrentBlockHash() string {
	if c != nil {
		return hex.EncodeToString(c.CurrentBlockHash)
	}
	return ""
}

func (c *ChainInfo) GetPreviousBlockHash() string {
	if c != nil {
		return hex.EncodeToString(c.PreviousBlockHash)
	}
	return ""
}
