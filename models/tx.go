package models

import (
	"time"
)

type Transaction struct {
	BlockNum uint64     `xorm:"blocknum" json:"block_num"`
	TxID     string     `xorm:"txid" json:"tx_id"`
	Createdt *time.Time `xorm:"createdt" json:"created"`
}

func (t *Transaction) Insert() (int64, error) {
	return engine.Insert(t)
}

// func Inserts(entities []interface{}) (int64, error) {
// 	return engine.Insert(entities...)
// }

func GetLastTxs(num int) ([]*Transaction, error) {
	txs := []*Transaction{}
	err := engine.Limit(num).Desc("createdt").Find(&txs)
	for _, tx := range txs {
		Createdtf, _ := time.Parse(TIME_LAYOUT, tx.Createdt.UTC().Format(TIME_LAYOUT))
		tx.Createdt = &Createdtf
	}
	return txs, err
}

func GetTxCount() (uint64, error) {
	tx := new(Transaction)
	count, err := engine.Count(tx)
	return uint64(count), err
}
