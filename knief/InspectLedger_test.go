package knief

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/ledger/blockledger"
	"testing"
)

func TestInspectLedger(t *testing.T) {
	fac, _, _ := GetFactory()
	testInsepctLedger(fac)
	testAlterLedger(fac)
}

func testInsepctLedger(fac blockledger.Factory) {
	high := GetLedgerHeight(fac)
	println("ledger hight = ", high)
	var num uint64 = 7
	block, _ := GetBlockFromNumber(fac, num)
	blockBytes, _ := SerializeBlock(block)
	l := len(blockBytes)
	println("block length =", l)
	datah := block.Header.PreviousHash
	x := common.BlockchainInfo{CurrentBlockHash: datah}
	bytesx, _ := json.Marshal(&x)
	println("block data hash =", string(bytesx))
	println("block data hash base64 =", base64.StdEncoding.EncodeToString(datah))
	txs, n, _ := ExtractEnvelopesFromBlock(block)
	fmt.Printf("there are %d transaction in number %d block\n", n, num)
	for i := 0; i < int(n); i++ {
		sig := (*txs)[i].Signature
		println(hex.EncodeToString(sig))
	}
}

func testAlterLedger(fac blockledger.Factory) {
	high := GetLedgerHeight(fac)
	Toyota := "Toyota"
	for i := 0; i < int(high); i++ {
		block, _ := GetBlockFromNumber(fac, uint64(i))
		if tf, _ := IsStringInBlock(*block, Toyota); tf {
			fmt.Printf("%s bytes is in block %dth \n", Toyota, i)
		}
	}
}
