package ChamHash

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/flogging"
	"net"
)

var logger = flogging.MustGetLogger("ChamHash")

type ChamHash struct {
	Hashvalue   []byte
	Randomvalue []byte
	Etdcipher   []byte
}

const SecurityParameter int = 1024
const CipherHost = "192.168.0.2:1234"

type ChamPublicKey struct {
	publicKey [SecurityParameter]byte
}

type ChamPrivateKey struct {
	privateKey [SecurityParameter]byte
}

// this function receive prpbytes without field in chamhash,
// and return (filled prpbytes, hashvalue in prpbytes for sign)

func BytesChamHashFromBytes(longbytes []byte) []byte {
	s := sha256.New()
	s.Write(longbytes)
	sha2value := s.Sum(nil)
	//return MockBytesChamHash(sha2value)
	return BytesChamHashFromSHA256(sha2value)
}

func MockBytesChamHash(sha256Hash []byte) []byte {
	chamhx := common.Chamhash{
		HashValue:   []byte("123"),
		RandomValue: []byte("235"),
		Etdcipher:   []byte("324"),
	}
	b, err := proto.Marshal(&chamhx)
	if err != nil {
		logger.Infof("chamhash from sha256 fail!")
	}
	return b

}

func BytesChamHashFromSHA256(sha256Hash []byte) []byte {
	chamhx := ChamHashHash(sha256Hash)
	b, err := proto.Marshal(&chamhx)
	if err != nil {
		logger.Info("chamhash from sha256 fail!")
	}
	return b
}

func HashValueFromChamHashBytes(chamHashBytes []byte) []byte {
	if chamHashBytes == nil {
		return nil
	}
	chamhash := common.Chamhash{}
	err := proto.Unmarshal(chamHashBytes, &chamhash)
	if err != nil {
		logger.Info("hash value from cham hash fail!")
	}
	return chamhash.GetHashValue()
}

// this function change chamhash while message changes, but chamhash.Hashvalue is fixed.
func ChangeChamHash(m1 []byte, m2 []byte, OldChamHash []byte) []byte {
	return []byte("")
}

// return (private_key, public_key)
func ChamHashKeyGen() (ChamPrivateKey, ChamPublicKey) {
	return ChamPrivateKey{}, ChamPublicKey{}
}

// return (ChamHash(hashValue,randomValue,Etdcipher))
func ChamHashHash(message []byte) common.Chamhash {
	conn, err := net.Dial("tcp", CipherHost)
	if err != nil {
		logger.Infof("can't connect %s: %s", CipherHost, err)
		logger.Infof("hashed message: ", hex.EncodeToString(message))
	}
	logger.Infof("success to connect %s", CipherHost)
	defer conn.Close()
	conn.Write(message)
	buf := make([]byte, SecurityParameter*10)
	count, err := conn.Read(buf)
	if err != nil {
		logger.Info("fail to read message from ", CipherHost)
	}
	if count < 4*SecurityParameter {
		logger.Info("bad formed message from ", CipherHost)
	}
	c := common.Chamhash{
		HashValue:   buf[:SecurityParameter],
		RandomValue: buf[SecurityParameter : SecurityParameter*2],
		Etdcipher:   buf[SecurityParameter*2 : SecurityParameter*4],
	}
	return c
}

func ChamHashCheck(message []byte, hash []byte) bool {
	return true
}

func ChamHashAdapt(privatekey ChamPrivateKey, m1 []byte, m2 []byte, hash common.Chamhash) common.Chamhash {
	return common.Chamhash{}
}

func GetHashFromChamHash(hash common.Chamhash) []byte {
	return []byte("")
}

//mock
