package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

func exists(args []string) (res interface{}, err error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "exists")
	}
	var b *bolt.Bucket
	DB.View(func(tx *bolt.Tx) error {
		b = tx.Bucket([]byte(args[0]))
		return nil
	})
	if b == nil {
		return 0, nil
	}
	return 1, nil
}

type cmd func([]string) (interface{}, error)

var cmdMap = map[string]cmd{
	"exists": exists,
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
	case int:
		return strconv.Itoa(res)
	default:
		panic(fmt.Sprintf(
			"The type of result returns from command '%s' with args %v is unsupported",
			cmd, args))
	}
}
