package knief

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	ab "github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/ledger/blockledger"
	"github.com/hyperledger/fabric/common/ledger/blockledger/fileledger"
	"github.com/hyperledger/fabric/common/metrics/disabled"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("fetchBlock")
var DataLoc = "./production/orderer"
var ChannelName = "mychannel"
// fetch block from dir

func GetLedgerHeight(fac blockledger.Factory) uint64{
	rw,_ := fac.GetOrCreate(ChannelName)
	h := rw.Height()
	logger.Info("getLedgerHeight:",h)
	return h
}

func GetBlockFromNumber(fac blockledger.Factory, num uint64) (*common.Block, error){
	rw,_ := fac.GetOrCreate(ChannelName)
	h := rw.Height()
	if h < num {
		return nil, errors.New("getBlock Num is higher than height of ledger")
	}
	// construct seekPosition
	seekPosition := &ab.SeekPosition{
		Type: &ab.SeekPosition_Specified{
			Specified: &ab.SeekSpecified{
				Number: num,
			},
		},
	}
	it, _ := rw.Iterator(seekPosition)
	block, _ := it.Next()
	return block,nil
}

func GetFactory() (blockledger.Factory, string, error){
	metricProv := disabled.Provider{}
	ld := DataLoc
	logger.Info("Ledger dir:", ld)
	lf, err := fileledger.New(ld, &metricProv)
	if err != nil {
		return nil, "", errors.WithMessage(err, "Error in opening ledger factory")
	}
	return lf, ld, nil
}

func SerializeBlock(block *common.Block) ([]byte, error) {
	blockbytes, _ := proto.Marshal(block)
	return blockbytes, nil
}

func DeserializeBlock(blk []byte) (*common.Block){
	b := common.Block{}
	proto.Unmarshal(blk, &b)
	return &b
}
