package ChamHash

import (
	"crypto/sha256"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/flogging"
	"net"
)

var logger = flogging.MustGetLogger("ChamHash")

const SecurityParameter int = 1024
const CipherHost = "192.168.0.2:1234"

// this function receive prpbytes without field in chamhash,
// and return (filled prpbytes, hashvalue in prpbytes for sign)

func Bytes2Sha256Bytes(message []byte) []byte {
	s := sha256.New()
	s.Write(message)
	return s.Sum(nil)
}

func BytesChamHashFromBytes(longbytes []byte) []byte {
	sha2value := Bytes2Sha256Bytes(longbytes)
	//return MockBytesChamHash(sha2value)
	return BytesChamHashFromSHA256(sha2value)
}

func BytesChamHashFromSHA256(sha256Hash []byte) []byte {
	//return ChamHashHash(sha256Hash)
	return ChamHashHash(sha256Hash)
}

func HashValueFromChamHashBytes(chamHashBytes []byte) []byte {
	if chamHashBytes == nil {
		return nil
	}
	chamhash := common.ChamHash{}
	err := proto.Unmarshal(chamHashBytes, &chamhash)
	if err != nil {
		logger.Errorf("hash value from cham hash fail!")
	}
	return []byte(chamhash.GetHash())
}

func MockChamHashHash(message []byte) []byte {
	ch := common.ChamHash{
		Hash:       "LLLLLLLLLLLL",
		HelperData: "LLLLLLLLLLLL",
	}
	chbytes, _ := proto.Marshal(&ch)
	return chbytes
}

// return ((hashValue,randomValue,Etdcipher))
func ChamHashHash(message []byte) []byte {
	conn, err := net.Dial("tcp", CipherHost)
	if err != nil {
		logger.Errorf("can't connect %s: %s", CipherHost, err.Error())
	} else {
		logger.Infof("chamHashHash success to connect %s", CipherHost)
	}
	defer conn.Close()

	m1 := common.Message{Mes: message}
	m1bytes, err := proto.Marshal(&m1)

	if err != nil {
		logger.Error("marshal message fail in chamHash message")
	}

	tf := common.Transfer{
		Ttype: common.TransferType_Hash,
		Tdata: m1bytes,
	}

	tfbytes, err := proto.Marshal(&tf)
	if err != nil {
		logger.Error("marshal transfer fail in chamHash message")
	}

	_, err = conn.Write(tfbytes)
	if err != nil {
		logger.Errorf("connection to %s fail to write message", CipherHost)
	}

	buffer := make([]byte, SecurityParameter*10)
	var count int
	for {
		count, err = conn.Read(buffer)
		if err == nil {
			break
		}
	}

	recvtf := common.Transfer{}

	err = proto.Unmarshal(buffer[:count], &recvtf)
	if err != nil {
		logger.Errorf("fail to unmarshal receive message from %s", CipherHost)
	}

	if recvtf.Ttype != common.TransferType_HashRecv {
		logger.Error("parse wrong package from %: should be HashRecv but get %s", CipherHost, common.TransferType_name[int32(recvtf.Ttype)])
	}

	return recvtf.Tdata
}

// message should be sha256 bytes
func ChamHashCheck(message []byte, hash []byte) bool {
	//hash -> sha256
	//s := sha256.New()
	//s.Write(message)
	//m := s.Sum(nil)

	conn, err := net.Dial("tcp", CipherHost)
	if err != nil {
		logger.Errorf("can't connect %s: %s", CipherHost, err.Error())
		return false
	} else {
		logger.Infof("ChamHashCheck success to connect %s", CipherHost)
	}
	defer conn.Close()

	m1 := common.Message{Mes: message}
	chFromBytes := common.ChamHash{}
	err = proto.Unmarshal(hash, &chFromBytes)
	if err != nil {
		logger.Error("hash couldn't be marshal to chamHash")
		return false
	}

	ck := common.CheckSet{M: &m1, Ch: &chFromBytes}
	ckbytes, err := proto.Marshal(&ck)
	if err != nil {
		logger.Error("marshal Checkset fail")
		return false
	}

	tf := common.Transfer{
		Ttype: common.TransferType_Check,
		Tdata: ckbytes,
	}

	tfbytes, err := proto.Marshal(&tf)
	if err != nil {
		logger.Error("marshal transfer fail in chamHash message")
	}

	_, err = conn.Write(tfbytes)
	if err != nil {
		logger.Errorf("connection to %s fail to write message", CipherHost)
	}

	// multi_thread? Timeout ?
	buffer := make([]byte, SecurityParameter*10)

	var count int
	for {
		count, err = conn.Read(buffer)
		if err == nil {
			break
		}
	}

	recvtf := common.Transfer{}

	err = proto.Unmarshal(buffer[:count], &recvtf)
	if err != nil {
		logger.Errorf("fail to unmarshal receive message from %s", CipherHost)
	}

	if recvtf.Ttype != common.TransferType_CheckRecv {
		logger.Error("parse wrong package from %: should be HashRecv but get %s", CipherHost, common.TransferType_name[int32(recvtf.Ttype)])
	}

	cs := common.CheckState{}
	_ = proto.Unmarshal(recvtf.Tdata, &cs)

	return cs.Check
}

/*
	collision : H(m1).hash = H(m2).hash
	check : check(m1, hash1) -> true; check(m2, hash2) -> true
	input: m1,m2,h1
	return: h2
	need : m1, m2 should be sha256 bytes
*/
func ChamHashAdapt(m1 []byte, m2 []byte, hash []byte) []byte {

	conn, err := net.Dial("tcp", CipherHost)
	if err != nil {
		logger.Errorf("can't connect %s: %s", CipherHost, err.Error())
		return nil
	} else {
		logger.Infof("ChamHashAdapt success to connect %s", CipherHost)
	}
	defer conn.Close()

	M1 := common.Message{Mes: m1}
	M2 := common.Message{Mes: m2}

	chFromBytes := common.ChamHash{}
	err = proto.Unmarshal(hash, &chFromBytes)
	if err != nil {
		logger.Error("hash couldn't be marshal to chamHash")
		return nil
	}

	ck := common.AdaptSet{M1: &M1, M2: &M2, Ch: &chFromBytes}
	ckbytes, err := proto.Marshal(&ck)
	if err != nil {
		logger.Error("marshal Checkset fail")
		return nil
	}

	tf := common.Transfer{
		Ttype: common.TransferType_Adapt,
		Tdata: ckbytes,
	}

	tfbytes, err := proto.Marshal(&tf)
	if err != nil {
		logger.Error("marshal transfer fail in chamHash message")
	}

	_, err = conn.Write(tfbytes)
	if err != nil {
		logger.Errorf("connection to %s fail to write message", CipherHost)
	}

	// multi_thread? Timeout ?
	buffer := make([]byte, SecurityParameter*10)
	var count int
	for {
		count, err = conn.Read(buffer)
		if err != nil {
			logger.Errorf("fail to read message from %s", CipherHost)
			logger.Error(err.Error())
		} else {
			break
		}
	}

	recvtf := common.Transfer{}

	err = proto.Unmarshal(buffer[:count], &recvtf)
	if err != nil {
		logger.Errorf("fail to unmarshal receive message from %s", CipherHost)
	}

	if recvtf.Ttype != common.TransferType_AdaptRecv {
		logger.Error("parse wrong package from %: should be HashRecv but get %s", CipherHost, common.TransferType_name[int32(recvtf.Ttype)])
	}

	return recvtf.Tdata
}
