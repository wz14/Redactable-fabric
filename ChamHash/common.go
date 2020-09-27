package ChamHash

import (
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
)

type ChamHash struct {
	Hashvalue            []byte
	Randomvalue          []byte
	Etdcipher            []byte
}

const SECURITY_PARAMETER int = 1024

type ChamPublicKey struct {
	publicKey [SECURITY_PARAMETER]byte
}

type ChamPrivateKey struct {
	privateKey [SECURITY_PARAMETER]byte
}

func InitPublicKey() ChamPublicKey{
	return ChamPublicKey{}
}

func BytesChamHashFromSHA256(sha256Hash []byte) []byte {
	champk := InitPublicKey()
	chamhx := ChamHashHash(champk,sha256Hash)
	b,err := proto.Marshal(&chamhx)
	if err != nil{
		print("chamhash from sha256 fail!")
	}
	return b
}

func HashValueFromChamHashBytes(chamHashBytes []byte) []byte {
	chamhash := common.Chamhash{}
	err := proto.Unmarshal(chamHashBytes,&chamhash)
	if err != nil {
		print("hash value from cham hash fail!")
	}
	return chamhash.GetHashValue()
}

// this function change chamhash while message changes, but chamhash.Hashvalue is fixed.
func ChangeChamHash(m1 []byte, m2 []byte, OldChamHash []byte)([]byte){
	return []byte("")
}

// return (private_key, public_key)
func ChamHashKeyGen()(ChamPrivateKey, ChamPublicKey) {
	return ChamPrivateKey{},ChamPublicKey{}
}

// return (ChamHash(hashValue,randomValue,Etdcipher))
func ChamHashHash(publickey ChamPublicKey, message []byte) (common.Chamhash){
	return common.Chamhash{}
}

func ChamHashCheck(publickey ChamPublicKey, message []byte, hash common.Chamhash) bool {
	return true
}

func ChamHashAdapt(privatekey ChamPrivateKey, m1 []byte, m2 []byte, hash common.Chamhash) (common.Chamhash){
	return common.Chamhash{}
}

func GetHashFromChamHash(hash common.Chamhash)([]byte){
	return []byte("")
}

/*
func main() {
	ChamHashFromProposalHash([]byte("Hello"))
}
*/
