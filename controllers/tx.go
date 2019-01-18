package controllers

import (
	"chain-sync/baas"
	"chain-sync/models"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	cb "github.com/hyperledger/fabric/protos/common"
	"encoding/json"
)

type BlockInfo struct {
	models.Block
	TxIds []string `json:"tx_ids"`
}

type RTData struct {
	Status    *Status    `json:"status"`
	BlockInfo *BlockInfo `json:"block_info"`
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
	var ids []string
	for _, txEnvBytes := range ledgerBlock.GetData().GetData() {
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
		ids = append(ids, chhd.GetTxId())
	}
	blockInfo.TxIds = ids
	block, _ := baas.GetBlock(blocknum)
	blockInfo.Block = *block
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

type OLTData struct {
	Status     *Status     `json:"status"`
	OriginData *OriginData `json:"origin_data"`
}

type OriginData struct {
	UserID     string `json:"user_id"`
	Ticket     string `json:"ticket"`
	Sympol     string `json:"symbol"`
	CMD        string `json:"cmd"`
	Volume     string `json:"volume"`
	OpenTime   int64  `json:"open_time"`
	OpenPrice  string `json:"open_price"`
	CloseTime  int64  `json:"close_time"`
	ClosePrice string `json:"close_price"`
}

func (m *MainController) GetOriginByHash() {
	rData := &OLTData{}
	status := &Status{}
	hash := m.GetString("hash")
	mode := m.GetString("mode")
	OriginDataStr, err := models.GetOriginByHash(mode, hash)
	if err != nil {
		status.Code = 500
		status.Msg = fmt.Sprintf("Get origin fail: %s", err)
		rData.Status = status
		m.Data["json"] = rData
		m.ServeJSON()
		return
	}
	status.Success()
	rData.Status = status
	originData := &OriginData{}
	err = json.Unmarshal([]byte(OriginDataStr),originData)
	if err != nil{
		status.Code = 500
		status.Msg = fmt.Sprintf("Json Unmarshal origin fail: %s", err)
		rData.Status = status
		m.Data["json"] = rData
		m.ServeJSON()
		return
	}
	rData.OriginData = originData
	m.Data["json"] = rData
	m.ServeJSON()
}
