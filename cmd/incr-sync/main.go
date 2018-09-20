package main

import (
	"chain-sync/baas"
	"flag"
	"fmt"
	"os"
)

const (
	GONUM      = 10
	HANDLERNUM = 10
)

var (
	from      uint64
	height    uint64
	curheight *uint64
	//控制goroutine数量
	chhandler  = make(chan bool, GONUM)
	chblocknum = make(chan uint64, HANDLERNUM)
	//历史区块同步超越标识
	passflag = make(chan bool, 1)
	//历史区块同步执行标识
	execflag = true
)

func sync(n uint64) {
	block, err := baas.GetBlock(n)
	if err == nil {
		fmt.Println(fmt.Sprintf("insert block %d", n))
		block.Insert()
	}
	txs, err := baas.GetTxsByBlockNum(n)
	if err == nil {
		for _, tx := range txs {
			fmt.Println(fmt.Sprintf("insert tx %s", tx.TxID))
			tx.Insert()
		}
	}
}

func push2chblocknum(blocknum uint64) {
	if pass := <-passflag; pass {
		chblocknum <- blocknum
		passflag <- true
		return
	}
	height = blocknum
	if execflag {
		execflag = false
		for i := from; i <= *curheight; i++ {
			chblocknum <- i
		}
		passflag <- true
	}
}

func main() {
	chaininfo, _ := baas.GetChainInfo()
	height = chaininfo.GetHeight()
	curheight = &height
	passflag <- false
	eventClient, _ := baas.GetBlockEventClient()
	reg, blockEventCh, _ := eventClient.RegisterBlockEvent()
	defer eventClient.Unregister(reg)

	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "block sync")
		flag.PrintDefaults()
	}
	flag.Uint64Var(&from, "from", 0, "The beginning block num.")
	flag.Parse()

	for i := 0; i < GONUM; i++ {
		chhandler <- true
	}

	// go func() {
	// 	for i := from; i <= height; i++ {
	// 		chblocknum <- i
	// 	}
	// }()

	go func() {
		for {
			select {
			case blockEvent := <-blockEventCh:
				push2chblocknum(blockEvent.Block.Header.Number)
			}
		}
	}()

	for {
		select {
		case num := <-chblocknum:
			select {
			case <-chhandler:
				go func() {
					sync(num)
					chhandler <- true
				}()
			}
			// case blockEvent := <-blockEventCh:
			// fmt.Printf("received block event: add block (%d)\n", blockEvent.Block.Header.Number)
			// go sync(blockEvent.Block.Header.Number)

		}
	}
}
