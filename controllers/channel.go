package controllers

import (
	"chain-sync/models"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type RCData struct {
	Status    *Status           `json:"status"`
	ChainInfo *models.ChainInfo `json:"chain_info"`
}

func (m *MainController) GetChannelInfo() {
	rData := &RCData{}
	status := &Status{}
	chainInfo := &models.ChainInfo{}
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
	chainInfoResponse, err := ledgerClient.QueryInfo()
	if err != nil {
		status.Code = 500
		status.Msg = fmt.Sprintf("queryinfo fail: %s", err)
		rData.Status = status
		m.Data["json"] = rData
		m.ServeJSON()
		return
	}
	status.Success()
	rData.Status = status
	chainInfo.Height = chainInfoResponse.BCI.GetHeight()
	chainInfo.CurrentBlockHash = chainInfoResponse.BCI.GetCurrentBlockHash()
	chainInfo.PreviousBlockHash = chainInfoResponse.BCI.GetPreviousBlockHash()
	chainInfo.TxCount, _ = models.GetTxCount()
	rData.ChainInfo = chainInfo
	m.Data["json"] = rData
	m.ServeJSON()
	return
}
