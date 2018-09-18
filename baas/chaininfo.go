package baas

import (
	"chain-sync/models"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func GetChainInfo() (*models.ChainInfo, error) {
	chainInfo := &models.ChainInfo{}
	sdk, err := fabsdk.New(config.FromFile("./sdk.yaml"))
	if err != nil {
		return chainInfo, errors.New(fmt.Sprintf("Failed to create new SDK: %s", err))
	}
	defer sdk.Close()
	channelProvider := sdk.ChannelContext(channelID,
		fabsdk.WithUser(user),
		fabsdk.WithOrg(org))
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		return chainInfo, errors.New(fmt.Sprintf(fmt.Sprintf("create ledger client fail: %s", err)))
	}
	chainInfoResponse, err := ledgerClient.QueryInfo()
	if err != nil {
		return chainInfo, errors.New(fmt.Sprintf("queryinfo fail: %s", err))
	}
	chainInfo.Height = chainInfoResponse.BCI.GetHeight()
	chainInfo.CurrentBlockHash = chainInfoResponse.BCI.GetCurrentBlockHash()
	chainInfo.PreviousBlockHash = chainInfoResponse.BCI.GetPreviousBlockHash()
	return chainInfo, nil
}
