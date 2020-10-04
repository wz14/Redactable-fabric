package main

import (
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/ChamHash"
)

func main() {
	h := ChamHash.BytesChamHashFromBytes([]byte("12354666"))
	chamhash := common.Chamhash{}
	_ = proto.Unmarshal(h, &chamhash)
}

