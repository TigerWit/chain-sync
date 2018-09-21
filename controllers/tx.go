package controllers

import (
	"chain-sync/models"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	cb "github.com/hyperledger/fabric/protos/common"
	"time"
)

type BlockInfo struct {
	models.Block
	TxIds []string `json:"tx_ids"`
}

type RTData struct {
	Status    *Status    `json:"status"`
	BlockInfo *BlockInfo `block_info`
}

func (m *MainController) GetTxsByBlocknum() {
	blocknum, _ := m.GetUint64("blocknum")
	rData := &RTData{}
	status := &Status{}
	sdk, err := fabsdk.New(config.FromFile("./sdk.yaml"))
	if err != nil {
		status.Code = 500
		status.Msg = fmt.Sprintf("Failed to create new SDK: %s", err)
		rData.Status = status
		m.Data["json"] = rData
		m.ServeJSON()
		return
	}
	defer sdk.Close()
	channelProvider := sdk.ChannelContext(channelID,
		fabsdk.WithUser(user),
		fabsdk.WithOrg(org))
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		status.Code = 500
		status.Msg = fmt.Sprintf("create ledger client fail: %s", err)
		rData.Status = status
		m.Data["json"] = rData
		m.ServeJSON()
		return
	}

	ledgerBlock, err := ledgerClient.QueryBlock(blocknum)
	if err != nil {
		status.Code = 500
		status.Msg = fmt.Sprintf("query block fail: %s", err)
		rData.Status = status
		m.Data["json"] = rData
		m.ServeJSON()
		return
	}
	blockInfo := &BlockInfo{}
	blockInfo.BlockNum = ledgerBlock.Header.Number
	blockInfo.PreHash = hex.EncodeToString(ledgerBlock.Header.PreviousHash)
	blockInfo.BlockHash = hex.EncodeToString(util.ComputeSHA256(tobytes(ledgerBlock.Header)))
	blockInfo.TxCount = len(ledgerBlock.Data.Data)

	var ids []string
	for index, txEnvBytes := range ledgerBlock.GetData().GetData() {
		txEnv := &cb.Envelope{}
		if err := proto.Unmarshal(txEnvBytes, txEnv); err != nil {
			status.Code = 500
			status.Msg = fmt.Sprintf("error reconstructing txenvelope(%s)", err)
			rData.Status = status
			m.Data["json"] = rData
			m.ServeJSON()
			return
		}

		payload := &cb.Payload{}
		if err := proto.Unmarshal(txEnv.Payload, payload); err != nil {
			status.Code = 500
			status.Msg = fmt.Sprintf("error reconstructing payload(%s)", err)
			rData.Status = status
			m.Data["json"] = rData
			m.ServeJSON()
			return
		}

		chhd := &cb.ChannelHeader{}
		if err := proto.Unmarshal(payload.Header.ChannelHeader, chhd); err != nil {
			status.Code = 500
			status.Msg = fmt.Sprintf("error reconstructing channelheader(%s)", err)
			rData.Status = status
			m.Data["json"] = rData
			m.ServeJSON()
			return
		}
		if index == 0 {
			createdt := time.Unix(chhd.GetTimestamp().Seconds, int64(chhd.GetTimestamp().GetNanos()))
			blockInfo.Createdt = &createdt
		}
		ids = append(ids, chhd.GetTxId())
	}
	blockInfo.TxIds = ids
	status.Success()
	rData.Status = status
	rData.BlockInfo = blockInfo
	m.Data["json"] = rData
	m.ServeJSON()
}

type RLTData struct {
	Status *Status               `json:"status"`
	Txs    []*models.Transaction `json:"transactions"`
}

func (m *MainController) GetLastTxs() {
	rData := &RLTData{}
	status := &Status{}
	num, _ := m.GetInt("num")
	txs, err := models.GetLastTxs(num)
	if err != nil {
		status.Code = 500
		status.Msg = fmt.Sprintf("Query txs fail: %s", err)
		rData.Status = status
		m.Data["json"] = rData
		m.ServeJSON()
		return
	}
	status.Success()
	rData.Status = status
	rData.Txs = txs
	m.Data["json"] = rData
	m.ServeJSON()
}
