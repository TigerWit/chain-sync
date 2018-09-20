package main

import (
	"chain-sync/baas"
	"flag"
	"fmt"
	"os"
)

const (
	GONUM = 10
)

var (
	from uint64
	to   uint64
	//控制goroutine数量
	complete = make(chan bool, GONUM)
	overflag = make(chan bool)
)

func main() {
	chaininfo, _ := baas.GetChainInfo()

	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "block sync")
		flag.PrintDefaults()
	}
	flag.Uint64Var(&from, "from", 0, "The beginning block num.")
	flag.Uint64Var(&to, "to", chaininfo.GetHeight(), "The end block num.")
	flag.Parse()

	for i := 0; i < GONUM; i++ {
		complete <- true
	}

	for i := from; i <= to; i++ {
		select {
		case <-complete:
			go func(n uint64) {
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
				complete <- true
				if n == to {
					overflag <- true
				}
			}(i)
		}
	}
	_ = <-overflag
}
