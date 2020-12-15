package ChamHash

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
)

// while modify the payload, change the chamhash struct content.
// paylaod.ChamHash and newpayload.ChamHash should be empty.
func UpdatePaylaod(charmHashBytes []byte, payloadBytes []byte, newpayloadBytes []byte) ([]byte, error) {
	payload := common.Payload{}
	err := proto.Unmarshal(payloadBytes, &payload)

	if err != nil {
		return nil, err
	}

	if payload.Chamhash != nil {
		return nil, errors.New("payload is not well formed")
	}

	newpayload := common.Payload{}
	err = proto.Unmarshal(newpayloadBytes, &newpayload)

	if err != nil {
		return nil, err
	}

	if newpayload.Chamhash != nil {
		return nil, errors.New("nwe payload is not well formed")
	}

	newCharmHash := ChamHashAdapt(Bytes2Sha256Bytes(payloadBytes), Bytes2Sha256Bytes(newpayloadBytes), charmHashBytes)
	newpayload.Chamhash = newCharmHash
	UpdatedPayloadBytes, err := proto.Marshal(&newpayload)

	if err != nil {
		return nil, err
	}
	return UpdatedPayloadBytes, nil
}

// fill payload with chamhash. return (filledPayloadBytes, hashValueOfPayload)
func FillPayload(PayloadBytes []byte) ([]byte, []byte) {
	payload := common.Payload{}
	proto.Unmarshal(PayloadBytes, &payload)

	if payload.Chamhash != nil {
		print("exist chamstruct in PaylaodBytes\n")
		return nil, nil
	}

	chash := BytesChamHashFromBytes(PayloadBytes)

	payload.Chamhash = chash

	valuehash := HashValueFromChamHashBytes(chash)

	filledPayloadBytes, _ := proto.Marshal(&payload)

	return filledPayloadBytes, valuehash
}

func GetHashOfPayloadStructure(filledPaylaodBytes []byte) []byte {
	payload := common.Payload{}
	proto.Unmarshal(filledPaylaodBytes, &payload)

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
	proto.Unmarshal(filledPayloadBytes, &payload)

	if payload.Chamhash == nil {
		print("chamstruct is not filled in PayloadBytes in checking\n")
		return false
	}

	chash := payload.Chamhash
	payload.Chamhash = nil
	orginBytes, _ := proto.Marshal(&payload)
	return ChamHashCheck(Bytes2Sha256Bytes(orginBytes), chash)
}
