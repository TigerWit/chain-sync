package baas

import (
	"chain-sync/models"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/util"
	cb "github.com/hyperledger/fabric/protos/common"
	"time"
)

func GetBlock(blocknum uint64) (*models.Block, error) {
	sdk, err := fabsdk.New(config.FromFile("./sdk.yaml"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to create new SDK: %s", err))
	}
	defer sdk.Close()
	channelProvider := sdk.ChannelContext(channelID,
		fabsdk.WithUser(user),
		fabsdk.WithOrg(org))
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(fmt.Sprintf("create ledger client fail: %s", err)))
	}
	ledgerBlock, err := ledgerClient.QueryBlock(blocknum)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(fmt.Sprintf("query block fail: %s", err)))
	}
	block := &models.Block{}
	block.BlockNum = ledgerBlock.Header.Number
	// block.DataHash = hex.EncodeToString(ledgerBlock.Header.DataHash)
	block.PreHash = hex.EncodeToString(ledgerBlock.Header.PreviousHash)
	block.BlockHash = hex.EncodeToString(util.ComputeSHA256(tobytes(ledgerBlock.Header)))
	block.TxCount = len(ledgerBlock.Data.Data)
	// block时间取block中第一笔交易的创建时间
	firstTxEnvBytes := ledgerBlock.GetData().GetData()[0]
	firstTxEnv := &cb.Envelope{}
	if err := proto.Unmarshal(firstTxEnvBytes, firstTxEnv); err != nil {
		return block, errors.New(fmt.Sprintf("error reconstructing txenvelope(%s)", err))
	}
	payload := &cb.Payload{}
	if err := proto.Unmarshal(firstTxEnv.Payload, payload); err != nil {
		return block, errors.New(fmt.Sprintf("error reconstructing payload(%s)", err))
	}

	chhd := &cb.ChannelHeader{}
	if err := proto.Unmarshal(payload.Header.ChannelHeader, chhd); err != nil {
		return block, errors.New(fmt.Sprintf("error reconstructing channelheader(%s)", err))
	}
	createdt := time.Unix(chhd.GetTimestamp().Seconds, int64(chhd.GetTimestamp().GetNanos()))
	block.Createdt = &createdt
	return block, nil
}

func GetBlockNew(blocknum uint64) (*models.Block, error) {
        channelProvider := SDKInstance.ChannelContext(channelID,
                fabsdk.WithUser(user),
                fabsdk.WithOrg(org))
        ledgerClient, err := ledger.New(channelProvider)
        if err != nil {
                return nil, errors.New(fmt.Sprintf(fmt.Sprintf("create ledger client fail: %s", err)))
        }
        ledgerBlock, err := ledgerClient.QueryBlock(blocknum)
        if err != nil {
                return nil, errors.New(fmt.Sprintf(fmt.Sprintf("query block fail: %s", err)))
        }
        block := &models.Block{}
        block.BlockNum = ledgerBlock.Header.Number
        // block.DataHash = hex.EncodeToString(ledgerBlock.Header.DataHash)
        block.PreHash = hex.EncodeToString(ledgerBlock.Header.PreviousHash)
        block.BlockHash = hex.EncodeToString(util.ComputeSHA256(tobytes(ledgerBlock.Header)))
        block.TxCount = len(ledgerBlock.Data.Data)
        // block时间取block中第一笔交易的创建时间
        firstTxEnvBytes := ledgerBlock.GetData().GetData()[0]
        firstTxEnv := &cb.Envelope{}
        if err := proto.Unmarshal(firstTxEnvBytes, firstTxEnv); err != nil {
                return block, errors.New(fmt.Sprintf("error reconstructing txenvelope(%s)", err))
        }
        payload := &cb.Payload{}
        if err := proto.Unmarshal(firstTxEnv.Payload, payload); err != nil {
                return block, errors.New(fmt.Sprintf("error reconstructing payload(%s)", err))
        }

        chhd := &cb.ChannelHeader{}
        if err := proto.Unmarshal(payload.Header.ChannelHeader, chhd); err != nil {
                return block, errors.New(fmt.Sprintf("error reconstructing channelheader(%s)", err))
        }
        createdt := time.Unix(chhd.GetTimestamp().Seconds, int64(chhd.GetTimestamp().GetNanos()))
        block.Createdt = &createdt
        return block, nil
}

func GetBlockEventClient() (*event.Client, error) {
	channelProvider := SDKInstance.ChannelContext(channelID,
		fabsdk.WithUser(user),
		fabsdk.WithOrg(org))
	return event.New(channelProvider, event.WithBlockEvents())
}

func GetBlocks(blocknums ...uint64) ([]*models.Block, error) {
	blocks := []*models.Block{}
	sdk, err := fabsdk.New(config.FromFile("./sdk.yaml"))
	if err != nil {
		return blocks, errors.New(fmt.Sprintf("Failed to create new SDK: %s", err))
	}
	defer sdk.Close()
	channelProvider := sdk.ChannelContext(channelID,
		fabsdk.WithUser(user),
		fabsdk.WithOrg(org))
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		return blocks, errors.New(fmt.Sprintf(fmt.Sprintf("create ledger client fail: %s", err)))
	}
	for _, blocknum := range blocknums {
		ledgerBlock, err := ledgerClient.QueryBlock(blocknum)
		if err != nil {
			return blocks, errors.New(fmt.Sprintf(fmt.Sprintf("query block fail: %s", err)))
		}
		block := &models.Block{}
		block.BlockNum = ledgerBlock.Header.Number
		// block.DataHash = hex.EncodeToString(ledgerBlock.Header.DataHash)
		block.PreHash = hex.EncodeToString(ledgerBlock.Header.PreviousHash)
		block.BlockHash = hex.EncodeToString(util.ComputeSHA256(tobytes(ledgerBlock.Header)))
		block.TxCount = len(ledgerBlock.Data.Data)
		// block.Createdt
		blocks = append(blocks, block)
	}
	return blocks, nil
}
