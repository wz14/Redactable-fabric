package ChamHash

import (
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
)

//fill payload with chamhash. return (filledPayloadBytes, hashValueOfPayload)
func FillPayload(PayloadBytes []byte)([]byte, []byte){
	payload := common.Payload{}
	proto.Unmarshal(PayloadBytes,&payload)

	if payload.Chamhash != nil{
		print("exist chamstruct in PaylaodBytes\n")
		return nil,nil
	}

	chash := BytesChamHashFromBytes(PayloadBytes)

	payload.Chamhash = chash

	valuehash := HashValueFromChamHashBytes(chash)

	filledPayloadBytes,_ := proto.Marshal(&payload)

	return filledPayloadBytes, valuehash
}

func GetHashOfPayloadStructure(filledPaylaodBytes []byte) []byte {
	payload := common.Payload{}
	proto.Unmarshal(filledPaylaodBytes,&payload)

	if payload.Chamhash == nil {
		print("chamstruct is not filled in PayloadBytes in gethash\n")
		return nil
	}

	return HashValueFromChamHashBytes(payload.Chamhash)
}

/*
	verify filledPayload With correct form and verify the signature
*/
func CheckFilledPayload(filledPayloadBytes []byte) bool {
	payload := common.Payload{}
	proto.Unmarshal(filledPayloadBytes,&payload)

	if payload.Chamhash == nil {
		print("chamstruct is not filled in PayloadBytes in checking\n")
		return false
	}

	chash := payload.Chamhash
	payload.Chamhash = nil
	orginBytes,_ := proto.Marshal(&payload)
	return ChamHashCheck(orginBytes,chash)
}
