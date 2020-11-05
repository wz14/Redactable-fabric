package main

import (
	"github.com/hyperledger/fabric/knief"
)

func main(){
	fac, _, _ := knief.GetFactory()
	_ = knief.GetLedgerHeight(fac)
	var num uint64 = 6
	block, _ := knief.GetBlockFromNumber(fac, num)
	blockBytes, _:= knief.SerializeBlock(block)
	l := len(blockBytes)
	println(l)
}