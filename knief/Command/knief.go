package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/tools/protolator"
	"github.com/hyperledger/fabric/knief"
	"os"
)

func main() {
	fac, _, _ := knief.GetFactory()
	defer fac.Close()

	fmt.Printf("this ledger (%s) height:%s", knief.ChannelName, knief.GetLedgerHeight(fac))
	block, _ := knief.GetBlockFromNumber(fac, 7)
	tx, _ := knief.ExtractEnvelope(block, 0)
	txbytes, _ := proto.Marshal(tx)
	println("tx bytes == ", len(txbytes))
	filename := knief.DataLoc + "/chains/mychannel/" + "blockfile_000000"
	newbytes, _ := knief.ChangeDevToABC(*tx)
	err := saveTx2File(txbytes, "./oldtx.json")
	if err != nil {
		println(err.Error())
	}
	err = saveTx2File(newbytes, "./newtx.json")
	if err != nil {
		println(err.Error())
	}
	err = knief.AlterBlockstorageForce(filename, txbytes, newbytes)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

func saveTx2File(txbytes []byte, filename string) error {
	tx := common.Envelope{}
	err := proto.Unmarshal(txbytes, &tx)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	err = protolator.DeepMarshalJSON(f, &tx)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
