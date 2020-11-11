package knief

import (
	"bytes"
	"errors"
	"github.com/hyperledger/fabric-protos-go/common"
)

func IsStringInBlock(b common.Block, s string) (bool,error) {
	bbytes, err := SerializeBlock(&b)
	if err != nil{
		return false, err
	}
	return bytes.Contains(bbytes, []byte(s)),nil
}

func AlterAllStringInBlock(b common.Block, old string, new string) ([]byte, error){
	j, err :=IsStringInBlock(b, old)
	if err != nil{
		return nil, err
	}else if j == false{
		return nil, errors.New("old string not in block bytes")
	}
	bbytes, err := SerializeBlock(&b)
	if err != nil{
		return nil, err
	}
	newbbytes := bytes.ReplaceAll(bbytes, []byte(old), []byte(new))
	return newbbytes, nil
}
