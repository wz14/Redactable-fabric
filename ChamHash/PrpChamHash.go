package ChamHash

import (
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/gogo/protobuf/proto"
)

func FillPrpStructureWithChamHash(prpbytes []byte) ([]byte, []byte){
	Prp := peer.ProposalResponsePayload{}
	err := proto.Unmarshal(prpbytes,&Prp)
	if err != nil{
		print("bad Prp bytes")
	}
	// check chamhash filled or not
	if Prp.ChamHashStruct != nil{
		return prpbytes,HashValueFromChamHashBytes(Prp.ChamHashStruct)
	}
	chambytes := BytesChamHashFromBytes(prpbytes)
	Prp.ChamHashStruct = chambytes
	realprpBytes,err := proto.Marshal(&Prp)
	if err!=nil{
		print("bad convert")
	}
	return realprpBytes,HashValueFromChamHashBytes(chambytes)
}

// this function get Hash from PrpStructure
func GetHashOfPrpStructure(prpbytes []byte) []byte {
	Prp := peer.ProposalResponsePayload{}
	err := proto.Unmarshal(prpbytes,&Prp)
	if err != nil{
		print("bad Prp bytes")
	}
	// check chamhash filled or not
	if Prp.ChamHashStruct == nil{
		print("bad prpbytes")
		return nil
	}
	return HashValueFromChamHashBytes(Prp.ChamHashStruct)
}

func CheckChamHashOfPrpStructure(prpbytes []byte) bool {
	Prp := peer.ProposalResponsePayload{}
	err := proto.Unmarshal(prpbytes,&Prp)
	if err != nil{
		print("bad Prp bytes")
	}
	// check chamhash filled or not
	if Prp.ChamHashStruct == nil{
		print("bad prpbytes")
		return false
	}
	chash := Prp.ChamHashStruct
	Prp.ChamHashStruct = nil
	orginBytes,_ := proto.Marshal(&Prp)
	return ChamHashCheck(orginBytes,chash)
}
