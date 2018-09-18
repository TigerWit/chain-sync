package baas

import (
	"encoding/asn1"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"math"
)

type asn1Header struct {
	Number       int64
	PreviousHash []byte
	DataHash     []byte
}

func tobytes(cb *common.BlockHeader) []byte {
	asn1Header := asn1Header{
		PreviousHash: cb.PreviousHash,
		DataHash:     cb.DataHash,
	}
	if cb.Number > uint64(math.MaxInt64) {
		panic(fmt.Errorf("Golang does not currently support encoding uint64 to asn1"))
	} else {
		asn1Header.Number = int64(cb.Number)
	}
	result, err := asn1.Marshal(asn1Header)
	if err != nil {
		panic(err)
	}
	return result
}
