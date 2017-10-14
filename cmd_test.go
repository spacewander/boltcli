package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
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
	assert.Equal(suite.T(), ExecCmdInCli("non-exist cmd"), "ERR unknown command 'non-exist cmd'")
}

func (suite *CmdSuite) TestExists() {
	assert.Equal(suite.T(), ExecCmdInCli("exists", "test"), "0")
	assert.Equal(suite.T(), ExecCmdInCli("exists"), "ERR wrong number of arguments for 'exists' command")

	DB.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("test"))
		return nil
	})
	assert.Equal(suite.T(), ExecCmdInCli("exists", "test"), "1")
}
