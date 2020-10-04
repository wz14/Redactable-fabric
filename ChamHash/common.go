package ChamHash

import (
	"crypto/sha256"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"net"
)

type ChamHash struct {
	Hashvalue            []byte
	Randomvalue          []byte
	Etdcipher            []byte
}

const SECURITY_PARAMETER int = 1024
const CIPHER_HOST = "127.0.0.1:1234"

type ChamPublicKey struct {
	publicKey [SECURITY_PARAMETER]byte
}

type ChamPrivateKey struct {
	privateKey [SECURITY_PARAMETER]byte
}

// this function receive prpbytes without field in chamhash,
// and return (filled prpbytes, hashvalue in prpbytes for sign)
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

func BytesChamHashFromBytes(longbytes []byte) []byte {
	s := sha256.New()
	s.Write(longbytes)
	sha2value := s.Sum(nil)
	return MockBytesChamHash(sha2value)
	//return BytesChamHashFromSHA256(sha2value)
}

func MockBytesChamHash(sha256Hash []byte) []byte {
	chamhx := common.Chamhash{
		HashValue:            []byte("123"),
		RandomValue:          []byte("235"),
		Etdcipher:            []byte("324"),
	}
	b,err := proto.Marshal(&chamhx)
	if err != nil{
		print("chamhash from sha256 fail!")
	}
	return b

}

func BytesChamHashFromSHA256(sha256Hash []byte) []byte {
	chamhx := ChamHashHash(sha256Hash)
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
func ChamHashHash(message []byte) (common.Chamhash){
	conn, err := net.Dial("tcp", CIPHER_HOST)
	if err != nil{
		print("bad network")
	}
	defer conn.Close()
	conn.Write(message)
	buf := make([]byte, SECURITY_PARAMETER*10)
	count,err := conn.Read(buf)
	if err != nil{
		print("bad network")
	}
	if count<4*SECURITY_PARAMETER{
		print("bad communication")
	}
	c := common.Chamhash{
		HashValue:            buf[:SECURITY_PARAMETER],
		RandomValue:          buf[SECURITY_PARAMETER : SECURITY_PARAMETER*2],
		Etdcipher:            buf[SECURITY_PARAMETER*2:SECURITY_PARAMETER*4],
	}
	return c
}

func ChamHashCheck(message []byte, hash []byte) bool {
	return true
}

func ChamHashAdapt(privatekey ChamPrivateKey, m1 []byte, m2 []byte, hash common.Chamhash) (common.Chamhash){
	return common.Chamhash{}
}

func GetHashFromChamHash(hash common.Chamhash)([]byte){
	return []byte("")
}

//mock
