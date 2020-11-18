package knief

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	ab "github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/common/ledger/blockledger"
	"github.com/hyperledger/fabric/common/ledger/blockledger/fileledger"
	"github.com/hyperledger/fabric/common/metrics/disabled"
	"github.com/pkg/errors"
)

// fetch block from dir

func ExtractEnvelopesFromBlock(block *common.Block) (*[]common.Envelope, uint64, error) {
	number := len(block.Data.Data)
	txs := make([]common.Envelope, number)
	for i := 0; i < number; i++ {
		err := proto.Unmarshal(block.Data.Data[i], &txs[i])
		if err != nil {
			return nil, 0, errors.WithMessage(err, "error in parsing the envelop in block")
		}
	}
	logger.Infof("Extract envelops from %dth block", block.Header.Number)
	return &txs, uint64(number), nil
}

func ExtractEnvelope(block *common.Block, index int) (*common.Envelope, error) {
	if block.Data == nil {
		return nil, errors.New("block data is nil")
	}

	envelopeCount := len(block.Data.Data)
	if index < 0 || index >= envelopeCount {
		return nil, errors.New("envelope index out of bounds")
	}
	marshaledEnvelope := block.Data.Data[index]
	envelope, err := GetEnvelopeFromBlock(marshaledEnvelope)
	if err != nil {
		err = errors.WithMessagef(err, "block data does not carry an envelope at index %d", index)
		return nil, err
	}
	logger.Infof("Extract %dth envelops from %dth block", index, block.Header.Number)
	return envelope, nil
}

func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
	// Block always begins with an envelope
	var err error
	env := &common.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Envelope")
	}

	return env, nil
}

// helper function
func ExtractHeaderFromBlock(block *common.Block) (*common.BlockHeader, error) {
	return block.Header, nil
}

func ExtractDataFromBlock(block *common.Block) (*common.BlockData, error) {
	return block.Data, nil
}

func ExtractMetaDataFromBlock(block *common.Block) (*common.BlockMetadata, error) {
	return block.Metadata, nil
}

func GetLedgerHeight(fac blockledger.Factory) uint64 {
	rw, err := fac.GetOrCreate(ChannelName)
	if err != nil {
		logger.Info("factory get channel manager fail")
	}
	h := rw.Height()
	logger.Info("getLedgerHeight:", h)
	return h
}

func GetBlockFromNumber(fac blockledger.Factory, num uint64) (*common.Block, error) {
	rw, err := fac.GetOrCreate(ChannelName)
	if err != nil {
		return nil, err
	}
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
	logger.Infof("get %dth block from ledger", num)
	return block, nil
}

func GetFactory() (blockledger.Factory, string, error) {
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
	blockbytes, err := proto.Marshal(block)
	if err != nil {
		return nil, err
	}
	return blockbytes, nil
}

func DeserializeBlock(blk []byte) (*common.Block, error) {
	b := common.Block{}
	err := proto.Unmarshal(blk, &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
