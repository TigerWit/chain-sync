package baas

import (
	"chain-sync/models"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	// "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	cb "github.com/hyperledger/fabric/protos/common"
	"time"
)

func GetTxsByBlockNum(blocknum uint64) ([]*models.Transaction, error) {
	blocknum = blocknum - uint64(originBlockNum)
	txs := []*models.Transaction{}
	channelProvider := SDKInstance.ChannelContext(channelID,
		fabsdk.WithUser(user),
		fabsdk.WithOrg(org))
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		return txs, errors.New(fmt.Sprintf("create ledger client fail: %s", err))
	}
	ledgerBlock, err := ledgerClient.QueryBlock(blocknum)
	if err != nil {
		return txs, errors.New(fmt.Sprintf("query block fail: %s", err))
	}

	for _, txEnvBytes := range ledgerBlock.GetData().GetData() {
		tx := &models.Transaction{}
		tx.BlockNum = blocknum
		txEnv := &cb.Envelope{}
		if err := proto.Unmarshal(txEnvBytes, txEnv); err != nil {
			return txs, errors.New(fmt.Sprintf("error reconstructing txenvelope(%s)", err))
		}

		payload := &cb.Payload{}
		if err := proto.Unmarshal(txEnv.Payload, payload); err != nil {
			return txs, errors.New(fmt.Sprintf("error reconstructing payload(%s)", err))
		}

		chhd := &cb.ChannelHeader{}
		if err := proto.Unmarshal(payload.Header.ChannelHeader, chhd); err != nil {
			return txs, errors.New(fmt.Sprintf("error reconstructing channelheader(%s)", err))
		}
		tx.TxID = chhd.GetTxId()
		createdt := time.Unix(chhd.GetTimestamp().Seconds, int64(chhd.GetTimestamp().GetNanos()))
		tx.Createdt = &createdt
		txs = append(txs, tx)
	}
	return txs, nil
}
