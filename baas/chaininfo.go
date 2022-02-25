package baas

import (
	"chain-sync/models"
	"errors"
	"fmt"
)

func GetChainInfo() (*models.ChainInfo, error) {
	chainInfo := &models.ChainInfo{}
	//sdk, err := fabsdk.New(config.FromFile("./sdk.yaml"))
	//if err != nil {
	//	return chainInfo, errors.New(fmt.Sprintf("Failed to create new SDK: %s", err))
	//}
	//defer sdk.Close()
	//channelProvider := sdk.ChannelContext(channelID,
	//	fabsdk.WithUser(user),
	//	fabsdk.WithOrg(org))
	//ledgerClient, err := ledger.New(channelProvider)
	//if err != nil {
	//	return chainInfo, errors.New(fmt.Sprintf(fmt.Sprintf("create ledger client fail: %s", err)))
	//}
	//chainInfoResponse, err := ledgerClient.QueryInfo()
	//if err != nil {
	//	return chainInfo, errors.New(fmt.Sprintf("queryinfo fail: %s", err))
	//}
	//chainInfo.Height = chainInfoResponse.BCI.GetHeight()
	//chainInfo.CurrentBlockHash = chainInfoResponse.BCI.GetCurrentBlockHash()
	//chainInfo.PreviousBlockHash = chainInfoResponse.BCI.GetPreviousBlockHash()
	blocks, err := models.GetLastBlocks(1)
	if err != nil {
		return chainInfo, errors.New(fmt.Sprintf("Failed to query explorer db: %s", err))
	}
	chainInfo.Height = blocks[0].BlockNum
	chainInfo.CurrentBlockHash = []byte(blocks[0].BlockHash)
	chainInfo.PreviousBlockHash = []byte(blocks[0].PreHash)
	return chainInfo, nil
}
