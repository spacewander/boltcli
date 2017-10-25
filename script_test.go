package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalLuaScript(t *testing.T) {
	tmpfile, _ := ioutil.TempFile("", "boltcli")
	dbPath := tmpfile.Name()
	initDB(dbPath)

	err := StartScript("test.lua")
	assert.Nil(t, err)

	DB.Close()
	os.Remove(dbPath)
}
