package knief

import (
	"github.com/gogo/protobuf/proto"
	"testing"
)

func TestAlterBlockstorage(t *testing.T) {
	fac, _, _ := GetFactory()
	block, _ := GetBlockFromNumber(fac, 7)
	tx, _ := ExtractEnvelope(block, 0)
	txbytes, _ := proto.Marshal(tx)
	filename := DataLoc + "/chains/mychannel/" + "blockfile_000000"
	err := AlterBlockstorageForce(filename, txbytes, txbytes)
	if err != nil {
		logger.Errorf("%s", err.Error())
	}
	fac.Close()
}

func TestAlterBlockstorageForce(t *testing.T) {
	fac, _, _ := GetFactory()
	defer fac.Close()
	block, _ := GetBlockFromNumber(fac, 7)
	tx, _ := ExtractEnvelope(block, 0)
	txbytes, _ := proto.Marshal(tx)
	filename := DataLoc + "/chains/mychannel/" + "blockfile_000000"
	newbytes, _ := ChangeDevToABC(*tx)
	err := AlterBlockstorageForce(filename, txbytes, newbytes)
	if err != nil {
		logger.Errorf("%s", err.Error())
	}
	fac.Close()

}
