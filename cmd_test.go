package main

import (
	"io/ioutil"
	"os"
	"testing"

	bolt "github.com/coreos/bbolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CmdSuite struct {
	suite.Suite
	dbPath string
}

func (suite *CmdSuite) SetupTest() {
	tmpfile, _ := ioutil.TempFile("", "boltcli")
	suite.dbPath = tmpfile.Name()
	initDB(suite.dbPath)
}

func TestCmdTestSuite(t *testing.T) {
	suite.Run(t, new(CmdSuite))
}

func (suite *CmdSuite) TearDownTest() {
	DB.Close()
	os.Remove(suite.dbPath)
}

func (suite *CmdSuite) TestExecCmdInCli() {
	assert.Equal(suite.T(), "ERR unknown command 'non-exist cmd'", ExecCmdInCli("non-exist cmd"))
	// Ignore the case of command name
	assert.Equal(suite.T(), "0", ExecCmdInCli("EXISTS", "bucket"))
}

func (suite *CmdSuite) TestExists() {
	assert.Equal(suite.T(), "0", ExecCmdInCli("exists", "bucket"))
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'exists' command", ExecCmdInCli("exists"))

	DB.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("bucket"))
		return nil
	})
	assert.Equal(suite.T(), "1", ExecCmdInCli("exists", "bucket"))
}

func (suite *CmdSuite) TestGet() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'get' command", ExecCmdInCli("get", "bucket"))
	assert.Equal(suite.T(), "ERR specific bucket 'bucket' does not exist", ExecCmdInCli("get", "bucket", "key"))

	DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("bucket"))
		if err != nil {
			return err
		}
		return b.Put([]byte("key"), []byte("value"))
	})
	assert.Equal(suite.T(), "value", ExecCmdInCli("get", "bucket", "key"))
	assert.Equal(suite.T(), "", ExecCmdInCli("get", "bucket", "non-exist"))
}

func (suite *CmdSuite) TestSet() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'set' command",
		ExecCmdInCli("set", "bucket", "key"))
	assert.Equal(suite.T(), "OK", ExecCmdInCli("set", "bucket", "key", "value"))
	assert.Equal(suite.T(), "value", ExecCmdInCli("get", "bucket", "key"))
}

func (suite *CmdSuite) TestDel() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'del' command", ExecCmdInCli("del"))
	assert.Equal(suite.T(), "OK", ExecCmdInCli("del", "bucket"))
	assert.Equal(suite.T(), "OK", ExecCmdInCli("del", "bucket", "key"))

	DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("bucket"))
		if err != nil {
			return err
		}
		return b.Put([]byte("key"), []byte("value"))
	})
	assert.Equal(suite.T(), "OK", ExecCmdInCli("del", "bucket", "key"))
	assert.Equal(suite.T(), "", ExecCmdInCli("get", "bucket", "key"))
	assert.Equal(suite.T(), "OK", ExecCmdInCli("del", "bucket", "key"))
	assert.Equal(suite.T(), "OK", ExecCmdInCli("del", "bucket"))
	assert.Equal(suite.T(), "0", ExecCmdInCli("exists", "bucket"))
}
