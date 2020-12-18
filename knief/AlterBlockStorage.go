package knief

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
)

func AlterBlockstorageForce(FileName string, oldtx []byte, newtx []byte) error {
	if len(oldtx) != len(newtx) {
		logger.Errorf("the length of oldtx bytes (%d) is not equal with the length of new tx (%d)", len(oldtx), len(newtx))
		return errors.New("the length of oldtx bytes is not equal with the length of new tx")
	}
	file, err := os.OpenFile(FileName, os.O_RDONLY, 0660)
	if err != nil {
		logger.Errorf("fail to read %s", FileName)
		return errors.New("fail to read all bytes")
	}
	content, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		logger.Errorf("fail to read all bytes from %s", FileName)
		return errors.New("fail to read all bytes")
	}

	file, err = os.OpenFile(FileName, os.O_WRONLY|os.O_TRUNC, 0660)
	defer file.Close()
	if bytes.Contains(content, oldtx) {
		logger.Infof("the old tx bytes is in that blockfile:%s ", FileName)
		// substr from blockfile (oldtx --> newtx)
		newcontent := bytes.ReplaceAll(content, oldtx, newtx)
		file.Write(newcontent)
		return nil
	} else {
		logger.Infof("the old tx bytes isn't in that blockfile$s ", FileName)
		return errors.New("no such tx bytes in blockfile")
	}
}

/*
func AlterBlockstorage(DirLoc string, ChannelID string, bbytes []byte, height int) error {
	p, err := fsblkstorage.NewProvider(
		fsblkstorage.NewConf(DirLoc, -1),
		&blkstorage.IndexConfig{
			AttrsToIndex: []blkstorage.IndexableAttr{blkstorage.IndexableAttrBlockNum}},
		&disabled.Provider{})

	fsblkstorage.NewConf(DirLoc, -1)
	indexConfig := blkstorage.IndexConfig{AttrsToIndex: []blkstorage.IndexableAttr{blkstorage.IndexableAttrBlockNum}}
	levelConfig := &leveldbhelper.Conf{
		DBPath:                getIndexDir(),
		ExpectedFormatVersion: dataFormatVersion(&indexConfig),
	}

	if err != nil {
		logger.Infof("get %d channel readwriter fail", ChannelID)
		return err
	}

	bs, err := p.OpenBlockStore(ChannelName)
	if err != nil {
		logger.Errorf("open block store fail with channel %s", ChannelName)
		return err
	}

	bs.GetBlockchainInfo()

	return nil
}

func GetFsBlkStorage(DirLoc string, ChannelID string,){
	conf := fsblkstorage.NewConf(DirLoc, -1)
	indexConfig := blkstorage.IndexConfig{AttrsToIndex: []blkstorage.IndexableAttr{blkstorage.IndexableAttrBlockNum}}
	dbConf := &leveldbhelper.Conf{
		DBPath:                conf.GetIndexDir(),
		ExpectedFormatVersion: dataFormatVersion(&indexConfig),
	}
	fileMgr := fsblkstorage.NewBlockfileMgr(ChannelID, conf, indexConfig, )

}


func dataFormatVersion(indexConfig *blkstorage.IndexConfig) string {
	// in version 2.0 we merged three indexable into one `IndexableAttrTxID`
	if indexConfig.Contains(blkstorage.IndexableAttrTxID) {
		return dataformat.Version20
	}
	return dataformat.Version1x
}

func getIndexDir() string {
	return filepath.Join(DataLoc, "index")
}
*/
