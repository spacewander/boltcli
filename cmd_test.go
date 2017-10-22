package main

import (
	"io/ioutil"
	"os"
	"strconv"
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
	assert.Equal(suite.T(), "false", ExecCmdInCli("EXISTS", "bucket"))
}

func (suite *CmdSuite) TestExists() {
	assert.Equal(suite.T(), "false", ExecCmdInCli("exists", "bucket"))
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'exists' command", ExecCmdInCli("exists"))

	DB.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("bucket"))
		return nil
	})
	assert.Equal(suite.T(), "true", ExecCmdInCli("exists", "bucket"))
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
	assert.Equal(suite.T(), "true", ExecCmdInCli("set", "bucket", "key", "value"))
	assert.Equal(suite.T(), "value", ExecCmdInCli("get", "bucket", "key"))
}

func (suite *CmdSuite) TestDel() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'del' command", ExecCmdInCli("del"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("del", "bucket"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("del", "bucket", "key"))

	DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("bucket"))
		if err != nil {
			return err
		}
		return b.Put([]byte("key"), []byte("value"))
	})
	assert.Equal(suite.T(), "true", ExecCmdInCli("del", "bucket", "key"))
	assert.Equal(suite.T(), "", ExecCmdInCli("get", "bucket", "key"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("del", "bucket", "key"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("del", "bucket"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("exists", "bucket"))
}

func (suite *CmdSuite) TestDelGlob() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'delglob' command", ExecCmdInCli("delglob"))
	assert.Equal(suite.T(), "0", ExecCmdInCli("delglob", "bucket*"))
	assert.Equal(suite.T(), "0", ExecCmdInCli("delglob", "bucket", "key*"))

	DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("bucket"))
		if err != nil {
			return err
		}
		b.CreateBucket([]byte("keyButBucketActually"))
		b.Put([]byte("key_1"), []byte("value"))
		return b.Put([]byte("key_2"), []byte("value"))
	})
	assert.Equal(suite.T(), "3", ExecCmdInCli("delglob", "bucket", "key*"))
	assert.Equal(suite.T(), "", ExecCmdInCli("get", "bucket", "key_1"))
	assert.Equal(suite.T(), "1", ExecCmdInCli("delglob", "bucket*"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("exists", "bucket"))
}

func (suite *CmdSuite) TestBuckets() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'buckets' command", ExecCmdInCli("buckets"))
	assert.Equal(suite.T(), "", ExecCmdInCli("buckets", "non-exist"))

	DB.Update(func(tx *bolt.Tx) error {
		for i := 0; i < 10; i++ {
			suffix := strconv.Itoa(i)
			tx.CreateBucket([]byte("key_" + suffix))
		}
		return nil
	})
	assert.Equal(suite.T(), `1) "key_1"`, ExecCmdInCli("buckets", "key_1"))
	assert.Equal(suite.T(),
		" 1) \"key_0\"\n 2) \"key_1\"\n 3) \"key_2\"\n 4) \"key_3\"\n 5) \"key_4\"\n 6) \"key_5\"\n 7) \"key_6\"\n 8) \"key_7\"\n 9) \"key_8\"\n10) \"key_9\"",
		ExecCmdInCli("buckets", "key_*"))
}

func (suite *CmdSuite) TestKeys() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'keys' command", ExecCmdInCli("keys", "bucket"))
	assert.Equal(suite.T(), "", ExecCmdInCli("keys", "non-exist", "key_*"))

	DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("bucket"))
		if err != nil {
			return err
		}
		for i := 0; i < 10; i++ {
			suffix := strconv.Itoa(i)
			b.Put([]byte("key_"+suffix), []byte("value_"+suffix))
		}
		return nil
	})
	assert.Equal(suite.T(), ``, ExecCmdInCli("keys", "bucket", "non-exist*"))
	assert.Equal(suite.T(), `1) "key_1"`, ExecCmdInCli("keys", "bucket", "key_1"))
	assert.Equal(suite.T(),
		" 1) \"key_0\"\n 2) \"key_1\"\n 3) \"key_2\"\n 4) \"key_3\"\n 5) \"key_4\"\n 6) \"key_5\"\n 7) \"key_6\"\n 8) \"key_7\"\n 9) \"key_8\"\n10) \"key_9\"",
		ExecCmdInCli("keys", "bucket", "key_*"))
}

func (suite *CmdSuite) TestKeyValues() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'keyvalues' command",
		ExecCmdInCli("keyvalues", "bucket"))
	assert.Equal(suite.T(), "", ExecCmdInCli("keyvalues", "non-exist", "key_*"))

	DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("bucket"))
		if err != nil {
			return err
		}
		for i := 0; i < 10; i++ {
			suffix := strconv.Itoa(i)
			b.Put([]byte("key_"+suffix), []byte("value_"+suffix))
		}
		return nil
	})
	assert.Equal(suite.T(), ``, ExecCmdInCli("keyvalues", "bucket", "non-exist*"))
	assert.Equal(suite.T(), `key_1) "value_1"`, ExecCmdInCli("keyvalues", "bucket", "key_1"))
	assert.Equal(suite.T(),
		"key_0) \"value_0\"\nkey_1) \"value_1\"\nkey_2) \"value_2\"\nkey_3) \"value_3\"\nkey_4) \"value_4\"\nkey_5) \"value_5\"\nkey_6) \"value_6\"\nkey_7) \"value_7\"\nkey_8) \"value_8\"\nkey_9) \"value_9\"",
		ExecCmdInCli("keyvalues", "bucket", "key_*"))
}

func (suite *CmdSuite) TestStats() {
	// warm up the stats
	DB.Update(func(tx *bolt.Tx) error {
		for i := 0; i < 10; i++ {
			suffix := strconv.Itoa(i)
			tx.CreateBucket([]byte(suffix))
		}
		return nil
	})
	assert.NotEqual(suite.T(), "", ExecCmdInCli("stats"))

	info, _ := stats()
	freeAlloc, _ := strconv.Atoi(info.(map[string]interface{})["FreeAlloc"].(string))
	assert.True(suite.T(), freeAlloc > 0)
	txStatusWrite, _ := strconv.Atoi(
		info.(map[string]interface{})["TxStats"].(map[string]interface{})["Write"].(string))
	assert.True(suite.T(), txStatusWrite > 0)
}
