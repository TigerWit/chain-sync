package models

import (
	"time"
)

type Block struct {
	BlockNum  uint64     `xorm:"blocknum" json:"block_num"`
	PreHash   string     `xorm:"prehash" json:"pre_hash"`
	BlockHash string     `xorm:"blockhash" json:"block_hash"`
	TxCount   int        `xorm:"txcount" json:"tx_count"`
	Createdt  *time.Time `xorm:"createdt" json:"created"`
}

func (b *Block) Insert() (int64, error) {
	return engine.Insert(b)
}

func GetLastBlocks(num int) ([]*Block, error) {
	blocks := []*Block{}
	err := engine.Limit(num).Desc("blocknum").Find(&blocks)
	for _, block := range blocks {
		Createdtf, _ := time.Parse(TIME_LAYOUT, block.Createdt.UTC().Format(TIME_LAYOUT))
		block.Createdt = &Createdtf
	}
	return blocks, err
}
