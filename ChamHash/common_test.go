package ChamHash

import (
	"bytes"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"testing"
)

func TestBytes2Sha256Bytes(t *testing.T) {
	m1 := []byte("123")
	m2 := []byte("456")
	m3 := []byte("123")
	sha1 := Bytes2Sha256Bytes(m1)
	sha2 := Bytes2Sha256Bytes(m2)
	sha3 := Bytes2Sha256Bytes(m3)
	if !bytes.Equal(sha1, sha3) {
		t.FailNow()
	}

	if bytes.Equal(sha1, sha2) {
		t.FailNow()
	}
}

func TestChamHashHash(t *testing.T) {
	m1 := []byte("123123")
	chbytes := ChamHashHash(m1)
	chamHash1 := common.ChamHash{}
	err := proto.Unmarshal(chbytes, &chamHash1)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if chamHash1.Hash == "" {
		t.Log("chamHash Hash string is null !")
		t.FailNow()
	}

	if chamHash1.HelperData == "" {
		t.Log("chamHash HelperData string is null !")
		t.FailNow()
	}

	chbytes = ChamHashHash(m1)
	chamHash2 := common.ChamHash{}
	err = proto.Unmarshal(chbytes, &chamHash2)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	if chamHash1.Hash != chamHash2.Hash || chamHash1.HelperData != chamHash2.HelperData {
		t.Log("chamHash is not equal with same message")
		t.FailNow()
	}
	//t.Logf("chamHash1.Hash: %s", chamHash1.Hash)
	//t.Logf("chamHash1.HelperData: %s", chamHash1.HelperData)
}

func BenchmarkChamHashHash(b *testing.B) {
	m1 := []byte("123123")
	chbytes := ChamHashHash(m1)
	chamHash1 := common.ChamHash{}
	err := proto.Unmarshal(chbytes, &chamHash1)
	if err != nil {
		b.Log(err.Error())
		b.FailNow()
	}
}

func TestChamHashCheck(t *testing.T) {
	m1 := []byte("345678")
	chbytes := ChamHashHash(m1)
	chamHash1 := common.ChamHash{}
	err := proto.Unmarshal(chbytes, &chamHash1)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	cs := ChamHashCheck(m1, chbytes)
	if !cs {
		t.Log("check pass fail")
		t.FailNow()
	}
}
