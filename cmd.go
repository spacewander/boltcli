package main

import (
	"fmt"
	"strconv"
	"strings"

	bolt "github.com/coreos/bbolt"
)

func del(args []string) (res interface{}, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "del")
	}
	err = DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(args[0]))
		if b == nil {
			return nil
		}
		if argsLen > 1 {
			return b.Delete([]byte(args[1]))
		}
		return tx.DeleteBucket([]byte(args[0]))
	})
	if err != nil {
		return nil, err
	}
	return "OK", nil
}

func exists(args []string) (res interface{}, err error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "exists")
	}
	var b *bolt.Bucket
	err = DB.View(func(tx *bolt.Tx) error {
		b = tx.Bucket([]byte(args[0]))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if b == nil {
		return false, nil
	}
	return true, nil
}

func get(args []string) (res interface{}, err error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "get")
	}
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(args[0]))
		if b == nil {
			return fmt.Errorf("specific bucket '%s' does not exist", args[0])
		}
		res = b.Get([]byte(args[1]))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func set(args []string) (res interface{}, err error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "set")
	}
	err = DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(args[0]))
		if b == nil {
			b, err = tx.CreateBucket([]byte(args[0]))
			if err != nil {
				return err
			}
		}
		return b.Put([]byte(args[1]), []byte(args[2]))
	})
	if err != nil {
		return nil, err
	}
	return "OK", nil
}

type cmd func([]string) (interface{}, error)

var cmdMap = map[string]cmd{
	"del":    del,
	"exists": exists,
	"get":    get,
	"set":    set,
}

// ExecCmdInCli run given cmd with args, return formatted string according to cmd result.
func ExecCmdInCli(cmd string, args ...string) string {
	f, ok := cmdMap[strings.ToLower(cmd)]
	if !ok {
		return fmt.Sprintf("ERR unknown command '%s'", cmd)
	}
	// keep the case unchanged so that we could distinguish
	// uppercase key from lowercase key.
	res, err := f(args)
	if err != nil {
		return fmt.Sprintf("ERR %v", err)
	}
	switch res := res.(type) {
	case bool:
		return strconv.FormatBool(res)
	case []byte:
		return string(res)
	case string:
		return res
	default:
		panic(fmt.Sprintf(
			"The type of result returns from command '%s' with args %v is unsupported",
			cmd, args))
	}
}
