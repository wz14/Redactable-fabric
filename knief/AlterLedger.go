package knief

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/ChamHash"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/pkg/errors"
)

func IsStringInBlock(b common.Block, s string) (bool, error) {
	bbytes, err := SerializeBlock(&b)
	if err != nil {
		return false, err
	}
	return bytes.Contains(bbytes, []byte(s)), nil
}

func IsStringInEnv(tx common.Envelope, s string) (bool, error) {
	txbytes, err := proto.Marshal(&tx)
	if err != nil {
		return false, err
	}
	return bytes.Contains(txbytes, []byte(s)), nil
}

func AlterAllStringInBlock(b common.Block, old string, new string) ([]byte, error) {
	j, err := IsStringInBlock(b, old)
	if err != nil {
		return nil, err
	} else if j == false {
		return nil, errors.New("old string not in block bytes")
	}
	bbytes, err := SerializeBlock(&b)
	if err != nil {
		return nil, err
	}
	newbbytes := bytes.ReplaceAll(bbytes, []byte(old), []byte(new))
	return newbbytes, nil
}

func AlterTxEnvelop(tx common.Envelope, old string, new string) ([]byte, error) {
	j, err := IsStringInEnv(tx, old)
	if err != nil {
		return nil, err
	} else if j == false {
		return nil, errors.New("old string not in tx bytes")
	}
	txbytes, err := proto.Marshal(&tx)
	if err != nil {
		return nil, err
	}
	newtxbytes := bytes.ReplaceAll(txbytes, []byte(old), []byte(new))
	return newtxbytes, nil
}

func AlterBlockWithTxBytes(b common.Block, txbytes []byte, loc int) ([]byte, error) {
	num := len(b.Data.Data)
	if loc >= num {
		return nil, errors.New("loc is bigger than the number of tx in block")
	}
	b.Data.Data[loc] = txbytes
	bbytes, err := SerializeBlock(&b)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to seriliaze block")
	}
	return bbytes, nil
}

// this is a function to alter the "change car9 to Dev" to the "change car9 to ABC"
func ChangeDevToABC(tx common.Envelope) ([]byte, error) {
	// find chaincodeproposal
	// unmarshal payload
	payload := common.Payload{}
	err := proto.Unmarshal(tx.Payload, &payload)
	if err != nil {
		return nil, err
	}

	// unmarshal transaction
	transaction := peer.Transaction{}
	err = proto.Unmarshal(payload.Data, &transaction)
	if err != nil {
		return nil, err
	}

	// unmarshal transaction actions
	action := transaction.Actions[0]
	cap := peer.ChaincodeActionPayload{}
	err = proto.Unmarshal(action.Payload, &cap)
	if err != nil {
		return nil, err
	}

	// unmarshal chaincodeproposalPyload
	cpp := peer.ChaincodeProposalPayload{}
	err = proto.Unmarshal(cap.ChaincodeProposalPayload, &cpp)
	if err != nil {
		return nil, err
	}

	//extract chaincodeSpec
	cis := peer.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(cpp.Input, &cis)
	if err != nil {
		return nil, err
	}

	functionArgs := cis.ChaincodeSpec.Input.Args
	for i, arg := range functionArgs {
		fmt.Printf("%d arg is %s.\n", i, string(arg))
	}

	// change function Args
	cis.ChaincodeSpec.Input.Args[2] = []byte("ABC")

	// show changed content
	functionArgs = cis.ChaincodeSpec.Input.Args
	for i, arg := range functionArgs {
		fmt.Printf("%d arg is %s.\n", i, string(arg))
	}
	//
	newCISBytes, err := proto.Marshal(&cis)
	if err != nil {
		return nil, err
	}

	//
	cpp.Input = newCISBytes
	newCPPBytes, err := proto.Marshal(&cpp)
	if err != nil {
		return nil, err
	}

	// change proposal hash 1, for proposal hash 2, there two proposal hash in proposal/txutils.go
	// one is clean the TransientMap filed, one is not clean. but actually it should be clean
	NewproposalHash, err := protoutil.GetProposalHash2(payload.Header, newCPPBytes)
	if err != nil {
		return nil, err
	}

	prp := peer.ProposalResponsePayload{}
	err = proto.Unmarshal(cap.Action.ProposalResponsePayload, &prp)
	if err != nil {
		return nil, err
	}
	logger.Infof("old prp bytes len = %d", len(cap.Action.ProposalResponsePayload))
	ss := []byte("{\"make\":\"Holden\",\"model\":\"Barina\",\"colour\":\"brown\",\"owner\":\"Dve\"}")
	if bytes.Contains(cap.Action.ProposalResponsePayload, ss) {
		logger.Infof("proposalResponsePayload====>%s", hex.EncodeToString(cap.Action.ProposalResponsePayload))
		logger.Infof("greate! the value of car9 is in prp")
	}

	chamHash := prp.ChamHashStruct

	prp.ChamHashStruct = nil
	EmptyPrpBytes, err := proto.Marshal(&prp)
	if err != nil {
		return nil, err
	}

	prp.ProposalHash = NewproposalHash
	// recalculate chamHashStruct in prp
	NewEmptyPrpBytes, err := proto.Marshal(&prp)
	if err != nil {
		return nil, err
	}

	//modifyPrpBytes, chamHash := ChamHash.FillPrpStructureWithChamHash(prpBytes)
	//here need a function in ChamHash convert(prpBytes, oldproposalHash, newproposalHash) to a new prpBytes
	//TODO: how to fight with cham hash structure.
	newchamHash, err := ChamHash.UpdateProposalResponsePayload(EmptyPrpBytes, NewEmptyPrpBytes, chamHash)
	if err != nil {
		return nil, err
	}

	prp.ChamHashStruct = newchamHash
	newPrpBytes, err := proto.Marshal(&prp)
	if err != nil {
		return nil, err
	}

	cap.ChaincodeProposalPayload = newCPPBytes
	// change prp bytes to new value of car9

	if bytes.Contains(newPrpBytes, ss) {
		logger.Infof("greate! the value of car9 is in new prp")
	}
	ss2 := []byte("{\"make\":\"Holden\",\"model\":\"Barina\",\"colour\":\"brown\",\"owner\":\"ABC\"}")
	finalPrpBytes := bytes.ReplaceAll(newPrpBytes, ss, ss2)
	if bytes.Contains(finalPrpBytes, ss2) {
		logger.Infof("greate! the value of car9 (ABC version) is in new prp")
	}

	logger.Infof("new prp bytes len = %d", len(finalPrpBytes))
	cap.Action.ProposalResponsePayload = finalPrpBytes

	NewCAPBytes, err := proto.Marshal(&cap)
	if err != nil {
		return nil, err
	}

	transaction.Actions[0].Payload = NewCAPBytes
	NewTransactionBytes, err := proto.Marshal(&transaction)
	if err != nil {
		return nil, err
	}

	oldChamHash := payload.Chamhash

	// empty payload.chamHash
	payload.Chamhash = nil
	PayloadBytes, err := proto.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	payload.Data = NewTransactionBytes
	NewPayloadBytes, err := proto.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	newPayloadChamHash, err := ChamHash.UpdatePaylaod(oldChamHash, PayloadBytes, NewPayloadBytes)
	if err != nil {
		return nil, err
	}
	payload.Chamhash = newPayloadChamHash

	NewPayloadBytes, err = proto.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	tx.Payload = NewPayloadBytes
	NewEnvelopBytes, err := proto.Marshal(&tx)
	if err != nil {
		return nil, err
	}

	return NewEnvelopBytes, nil
}
