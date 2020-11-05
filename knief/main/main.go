package main

import knief "github.com/hyperledger/fabric/knief"

func main(){
	h := knief.GetLedgerHeight()
	println(h)
}