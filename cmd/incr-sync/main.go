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
	from uint64
	//控制goroutine数量
	chhandler  = make(chan bool, GONUM)
	chblocknum = make(chan uint64, HANDLERNUM)
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

func main() {
	chaininfo, _ := baas.GetChainInfo()
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

	go func() {
		for i := from; i <= chaininfo.GetHeight(); i++ {
			chblocknum <- i
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
		case blockEvent := <-blockEventCh:
			fmt.Printf("received block event: add block (%d)\n", blockEvent.Block.Header.Number)
			go sync(blockEvent.Block.Header.Number)
		}
	}
}
