package knief

import (
	"fmt"
	"testing"
)

func TestChangeCarOwner(t *testing.T) {
	fac, _, _ := GetFactory()
	block, _ := GetBlockFromNumber(fac, 7)
	txs, num, _ := ExtractEnvelopesFromBlock(block)
	if num != 1 {
		fmt.Printf("get block envelop fail \n")
	}
	txbytes, _ := ChangeDevToABC((*txs)[0])
	block.Data.Data[0] = txbytes
	fac.Close()
}
