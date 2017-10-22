package main

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"

	bolt "github.com/coreos/bbolt"
	"github.com/gobwas/glob"
)

func del(args ...string) (res interface{}, err error) {
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
	return true, nil
}

func delGlob(args ...string) (res interface{}, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "delglob")
	}
	count := 0
	// Only one glob pattern is suppored
	pattern, err := glob.Compile(args[len(args)-1])
	if err != nil {
		return nil, err
	}
	err = DB.Update(func(tx *bolt.Tx) error {
		if argsLen == 1 {
			c := tx.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if pattern.Match(string(k)) {
					err = tx.DeleteBucket(k)
					if err != nil {
						return err
					}
					count++
				}
			}
		} else {
			b := tx.Bucket([]byte(args[0]))
			if b != nil {
				c := b.Cursor()
				for k, _ := c.First(); k != nil; k, _ = c.Next() {
					if pattern.Match(string(k)) {
						err = b.Delete(k)
						if err == bolt.ErrIncompatibleValue {
							err = b.DeleteBucket(k)
						}
						if err != nil {
							return err
						}
						count++
					}
				}
			}
		}
		return nil
	})
	return count, nil
}

func exists(args ...string) (res interface{}, err error) {
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

func get(args ...string) (res interface{}, err error) {
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

func set(args ...string) (res interface{}, err error) {
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
	return true, nil
}

func buckets(args ...string) (res interface{}, err error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "buckets")
	}
	pattern, err := glob.Compile(args[0])
	if err != nil {
		return
	}
	res = []string{}
	err = DB.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(bname []byte, b *bolt.Bucket) error {
			name := string(bname)
			if pattern.Match(name) {
				res = append(res.([]string), name)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func keys(args ...string) (res interface{}, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "keys")
	}
	pattern, err := glob.Compile(args[1])
	if err != nil {
		return
	}
	res = []string{}
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(args[0]))
		if b == nil {
			return nil
		}
		b.ForEach(func(k, v []byte) error {
			key := string(k)
			if pattern.Match(key) {
				res = append(res.([]string), key)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func keyvalues(args ...string) (res interface{}, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "keyvalues")
	}
	pattern, err := glob.Compile(args[1])
	if err != nil {
		return
	}
	res = map[string]interface{}{}
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(args[0]))
		if b == nil {
			return nil
		}
		b.ForEach(func(k, v []byte) error {
			key := string(k)
			if pattern.Match(key) {
				res.(map[string]interface{})[key] = string(v)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func stats(_ ...string) (res interface{}, err error) {
	info := map[string]interface{}{}
	stat := DB.Stats()
	val := reflect.ValueOf(stat)
	for i := 0; i < val.NumField(); i++ {
		valField := val.Field(i)
		typeField := val.Type().Field(i)
		if typeField.Type.Name() == "int" {
			info[typeField.Name] = strconv.FormatInt(valField.Int(), 10)
		}
	}
	val = reflect.ValueOf(stat.TxStats)
	info["TxStats"] = map[string]interface{}{}
	for i := 0; i < val.NumField(); i++ {
		valField := val.Field(i)
		typeField := val.Type().Field(i)
		if typeField.Type.Name() == "int" || typeField.Type.Name() == "Duration" {
			info["TxStats"].(map[string]interface{})[typeField.Name] = strconv.FormatInt(valField.Int(), 10)
		}
	}
	return info, nil
}

type cmd func(...string) (interface{}, error)

var cmdMap = map[string]cmd{
	"del":       del,
	"delglob":   delGlob,
	"exists":    exists,
	"get":       get,
	"set":       set,
	"buckets":   buckets,
	"keys":      keys,
	"keyvalues": keyvalues,
	"stats":     stats,
}

// Format ["o1", "o2"] to string
// 1) "o1"\n
// 2) "o2"
func formatListToStr(list []string) string {
	paddingNum := strconv.Itoa(int(math.Log10(float64(len(list)))) + 1)
	padded := make([]string, len(list))
	for i, data := range list {
		padded[i] = fmt.Sprintf("%"+paddingNum+`d) "%s"`, i+1, data)
	}
	return strings.Join(padded, "\n")
}

// Format {"a": "10", "b": "20", "c": {"c1": 30}} to string
// a) "10"\n
// b) "20"\n
// c)\n
//     c1) "30"
func formatMapToStr(collection map[string]interface{}, prefix string) string {
	formatted := make([]string, len(collection))
	keys := make([]string, len(collection))
	i := 0
	for k := range collection {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for i, k := range keys {
		switch v := collection[k].(type) {
		case string:
			formatted[i] = fmt.Sprintf(`%s%s) "%s"`, prefix, k, v)
		case map[string]interface{}:
			nestedMap := formatMapToStr(v, prefix+"    ")
			formatted[i] = fmt.Sprintf("%s%s)\n%s", prefix, k, nestedMap)
		}
		i++
	}
	return strings.Join(formatted, "\n")
}

// ExecCmdInCli run given cmd with args, return formatted string according to cmd result.
func ExecCmdInCli(cmd string, args ...string) string {
	f, ok := cmdMap[strings.ToLower(cmd)]
	if !ok {
		return fmt.Sprintf("ERR unknown command '%s'", cmd)
	}
	// keep the case unchanged so that we could distinguish
	// uppercase key from lowercase key.
	res, err := f(args...)
	if err != nil {
		return fmt.Sprintf("ERR %v", err)
	}
	switch res := res.(type) {
	case bool:
		return strconv.FormatBool(res)
	case []byte:
		return string(res)
	case []string:
		return formatListToStr(res)
	case map[string]interface{}:
		return formatMapToStr(res, "")
	case int:
		return strconv.Itoa(res)
	default:
		panic(fmt.Sprintf(
			"The type of result returns from command '%s' with args %v is unsupported",
			cmd, args))
	}
}
