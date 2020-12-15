package ChamHash

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// update Proposal response payload, return Changed prpBytes
func UpdateProposalResponsePayload(prpbytes []byte, newprpbytes []byte, chamhashBytes []byte) ([]byte, error) {
	Prp := peer.ProposalResponsePayload{}
	err := proto.Unmarshal(prpbytes, &Prp)
	if err != nil {
		return nil, err
	}

	if Prp.ChamHashStruct != nil {
		return nil, errors.New("chamHash is not empty in prpbytes")
	}

	newPrp := peer.ProposalResponsePayload{}
	err = proto.Unmarshal(newprpbytes, &newPrp)
	if err != nil {
		return nil, err
	}

	if newPrp.ChamHashStruct != nil {
		return nil, errors.New("chamHash is not empty in NewprpBytes")
	}

	//ChamHashAdapt
	updateChamHashBytes := ChamHashAdapt(Bytes2Sha256Bytes(prpbytes), Bytes2Sha256Bytes(newprpbytes), chamhashBytes)
	newPrp.ChamHashStruct = updateChamHashBytes

	updatePrpBytes, err := proto.Marshal(&newPrp)
	if err != nil {
		return nil, err
	}
	return updatePrpBytes, nil
}

// equalPayload is to judge p1 and p2 with chamHashStruct in
func EqualProposalResponsePayload(p1 []byte, p2 []byte) (bool, error) {
	// first verify(p1.chamhashstruct, p1 , pk)
	// second verify(p2.chamhashstruct, p2 , pk)
	return true, nil
}

func FillPrpStructureWithChamHash(prpbytes []byte) ([]byte, []byte) {
	Prp := peer.ProposalResponsePayload{}
	err := proto.Unmarshal(prpbytes, &Prp)
	if err != nil {
		print("bad Prp bytes")
	}
	// check chamhash filled or not
	if Prp.ChamHashStruct != nil {
		return prpbytes, HashValueFromChamHashBytes(Prp.ChamHashStruct)
	}
	chambytes := BytesChamHashFromBytes(prpbytes)
	Prp.ChamHashStruct = chambytes
	realprpBytes, err := proto.Marshal(&Prp)
	if err != nil {
		print("bad convert")
	}
	return realprpBytes, HashValueFromChamHashBytes(chambytes)
}

// this function get Hash from PrpStructure
func GetHashOfPrpStructure(prpbytes []byte) []byte {
	Prp := peer.ProposalResponsePayload{}
	err := proto.Unmarshal(prpbytes, &Prp)
	if err != nil {
		print("bad Prp bytes")
	}
	// check chamhash filled or not
	if Prp.ChamHashStruct == nil {
		print("bad prpbytes")
		return nil
	}
	return HashValueFromChamHashBytes(Prp.ChamHashStruct)
}

func CheckChamHashOfPrpStructure(prpbytes []byte) bool {
	Prp := peer.ProposalResponsePayload{}
	err := proto.Unmarshal(prpbytes, &Prp)
	if err != nil {
		print("bad Prp bytes")
	}
	// check chamhash filled or not
	if Prp.ChamHashStruct == nil {
		print("bad prpbytes")
		return false
	}
	chash := Prp.ChamHashStruct
	Prp.ChamHashStruct = nil
	orginBytes, _ := proto.Marshal(&Prp)
	return ChamHashCheck(Bytes2Sha256Bytes(orginBytes), chash)
}
