package knief

import (
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

func GetLedgerHeight() (uint64){
	fac,_,_ := getFactory()
	rw,_ := fac.GetOrCreate(ChannelName)
	h := rw.Height()
	return h
}

func getFactory() (blockledger.Factory, string, error){
	metricProv := disabled.Provider{}
	ld := DataLoc
	logger.Info("Ledger dir:", ld)
	lf, err := fileledger.New(ld, &metricProv)
	if err != nil {
		return nil, "", errors.WithMessage(err, "Error in opening ledger factory")
	}
	return lf, ld, nil
}
