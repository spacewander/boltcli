package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	bolt "go.etcd.io/bbolt"
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
		b, _ := tx.CreateBucket([]byte("bucket"))
		b.Put([]byte("key"), []byte("value"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		b.Put([]byte("key"), []byte("value"))
		return nil
	})
	assert.Equal(suite.T(), "true", ExecCmdInCli("exists", "bucket"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("exists", "bucket", "key"))

	assert.Equal(suite.T(), "true", ExecCmdInCli("exists", "bucket", "subbucket"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("exists", "bucket", "subbucket", "subbucket"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("exists", "bucket", "subbucket", "non-exist", "key"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("exists", "bucket", "subbucket", "subbucket", "key"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("exists", "bucket", "subbucket", "subbucket", "non-exist"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("exists", "bucket2", "subbucket2"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("exists", "bucket", "subbucket2"))
}

func (suite *CmdSuite) TestGet() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'get' command", ExecCmdInCli("get", "bucket"))
	assert.Equal(suite.T(), `""`, ExecCmdInCli("get", "bucket", "key"))

	DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("bucket"))
		if err != nil {
			return err
		}
		b.Put([]byte("key"), []byte("value"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		return b.Put([]byte("key"), []byte("value"))
	})
	assert.Equal(suite.T(), `"value"`, ExecCmdInCli("get", "bucket", "key"))
	assert.Equal(suite.T(), `""`, ExecCmdInCli("get", "bucket", "non-exist", "key"))
	assert.Equal(suite.T(), `""`, ExecCmdInCli("get", "bucket", "subbucket", "key"))
	assert.Equal(suite.T(), `"value"`, ExecCmdInCli("get", "bucket", "subbucket", "subbucket", "key"))
}

func (suite *CmdSuite) TestSet() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'set' command",
		ExecCmdInCli("set", "bucket", "key"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("set", "bucket", "key", "value"))
	assert.Equal(suite.T(), `"value"`, ExecCmdInCli("get", "bucket", "key"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("set", "bucket", "subbucket", "key", "value"))
	assert.Equal(suite.T(), `"value"`, ExecCmdInCli("get", "bucket", "subbucket", "key"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("set", "bucket1", "bucket2", "key", "value"))
	assert.Equal(suite.T(), `"value"`, ExecCmdInCli("get", "bucket1", "bucket2", "key"))
}

func (suite *CmdSuite) TestDel() {
	assert.Equal(suite.T(), "ERR wrong number of arguments for 'del' command", ExecCmdInCli("del"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("del", "bucket"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("del", "bucket", "key"))

	DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("bucket"))
		if err != nil {
			return err
		}
		b.Put([]byte("key"), []byte("value"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		return b.Put([]byte("key"), []byte("value"))
	})
	assert.Equal(suite.T(), "true", ExecCmdInCli("del", "bucket", "key"))
	assert.Equal(suite.T(), `""`, ExecCmdInCli("get", "bucket", "key"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("del", "bucket", "key"))

	assert.Equal(suite.T(), "false", ExecCmdInCli("del", "bucket", "subbucket", "non-exist"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("del", "bucket", "subbucket", "non-exist-bucket", "key"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("del", "bucket", "subbucket", "subbucket", "key"))
	assert.Equal(suite.T(), `""`, ExecCmdInCli("get", "bucket", "subbucket", "subbucket", "key"))
	assert.Equal(suite.T(), "true", ExecCmdInCli("del", "bucket", "subbucket"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("exists", "bucket", "subbucket"))

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
		b.Put([]byte("key_1"), []byte("value"))
		b.Put([]byte("key_2"), []byte("value"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		return b.Put([]byte("key"), []byte("value"))
	})
	assert.Equal(suite.T(), "0", ExecCmdInCli("delglob", "bucket", "non-exist", "*"))
	assert.Equal(suite.T(), "1", ExecCmdInCli("delglob", "bucket", "subbucket", "subbucket", "*"))
	assert.Equal(suite.T(), `""`, ExecCmdInCli("get", "bucket", "subbucket", "subbucket", "key"))
	assert.Equal(suite.T(), "1", ExecCmdInCli("delglob", "bucket", "subbucket", "sub*"))
	assert.Equal(suite.T(), "false", ExecCmdInCli("exists", "bucket", "subbucket", "subbucket"))

	assert.Equal(suite.T(), "3", ExecCmdInCli("delglob", "bucket", "*"))
	assert.Equal(suite.T(), `""`, ExecCmdInCli("get", "bucket", "key_1"))

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
		b := tx.Bucket([]byte("key_0"))
		b.Put([]byte("subb1ket"), []byte("value"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		b.CreateBucket([]byte("subbucket"))
		return nil
	})
	assert.Equal(suite.T(), `1) "key_1"`, ExecCmdInCli("buckets", "key_1"))
	assert.Equal(suite.T(),
		" 1) \"key_0\"\n 2) \"key_1\"\n 3) \"key_2\"\n 4) \"key_3\"\n 5) \"key_4\"\n 6) \"key_5\"\n 7) \"key_6\"\n 8) \"key_7\"\n 9) \"key_8\"\n10) \"key_9\"",
		ExecCmdInCli("buckets", "key_*"))
	assert.Equal(suite.T(), ``, ExecCmdInCli("buckets", "key_0", "non-exist-bucket", "*"))
	assert.Equal(suite.T(), `1) "subbucket"`, ExecCmdInCli("buckets", "key_0", "subb*ket"))
	assert.Equal(suite.T(), `1) "subbucket"`, ExecCmdInCli("buckets", "key_0", "subbucket", "*bucket"))
	assert.Equal(suite.T(), ``, ExecCmdInCli("buckets", "subb", "*bucket"))
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
		b, _ = b.CreateBucket([]byte("subbucket"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		b.Put([]byte("key"), []byte("value"))
		b.CreateBucket([]byte("keyButBucketActually"))
		return nil
	})
	assert.Equal(suite.T(), ``, ExecCmdInCli("keys", "bucket", "non-exist*"))
	assert.Equal(suite.T(), `1) "key_1"`, ExecCmdInCli("keys", "bucket", "key_1"))
	assert.Equal(suite.T(),
		" 1) \"key_0\"\n 2) \"key_1\"\n 3) \"key_2\"\n 4) \"key_3\"\n 5) \"key_4\"\n 6) \"key_5\"\n 7) \"key_6\"\n 8) \"key_7\"\n 9) \"key_8\"\n10) \"key_9\"",
		ExecCmdInCli("keys", "bucket", "key_*"))
	assert.Equal(suite.T(), ``, ExecCmdInCli("keys", "bucket", "subbucket", "k*"))
	assert.Equal(suite.T(), ``, ExecCmdInCli("keys", "bucket", "subbucket", "non-exist", "k*"))
	assert.Equal(suite.T(), `1) "key"`, ExecCmdInCli("keys", "bucket", "subbucket", "subbucket", "k*"))
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
		b, _ = b.CreateBucket([]byte("subbucket"))
		b, _ = b.CreateBucket([]byte("subbucket"))
		b.Put([]byte("key"), []byte("value"))
		return nil
	})
	assert.Equal(suite.T(), ``, ExecCmdInCli("keyvalues", "bucket", "non-exist*"))
	assert.Equal(suite.T(), `key_1) "value_1"`, ExecCmdInCli("keyvalues", "bucket", "key_1"))
	assert.Equal(suite.T(),
		"key_0) \"value_0\"\nkey_1) \"value_1\"\nkey_2) \"value_2\"\nkey_3) \"value_3\"\nkey_4) \"value_4\"\nkey_5) \"value_5\"\nkey_6) \"value_6\"\nkey_7) \"value_7\"\nkey_8) \"value_8\"\nkey_9) \"value_9\"",
		ExecCmdInCli("keyvalues", "bucket", "key_*"))
	assert.Equal(suite.T(), ``, ExecCmdInCli("keyvalues", "bucket", "subbucket", "k*"))
	assert.Equal(suite.T(), ``, ExecCmdInCli("keyvalues", "bucket", "subbucket", "non-exist", "k*"))
	assert.Equal(suite.T(), `key) "value"`, ExecCmdInCli("keyvalues", "bucket", "subbucket", "subbucket", "k*"))
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
	assert.NotEqual(suite.T(), `""`, ExecCmdInCli("stats"))

	info, _ := stats()
	freeAlloc, _ := info.(map[string]interface{})["FreeAlloc"].(int64)
	assert.True(suite.T(), freeAlloc > 0)
	txStatusWrite, _ := info.(map[string]interface{})["TxStats"].(map[string]interface{})["Write"].(int64)
	assert.True(suite.T(), txStatusWrite > 0)
}
