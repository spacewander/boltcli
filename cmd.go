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
	found := false
	err = DB.Update(func(tx *bolt.Tx) error {
		if argsLen == 1 {
			err = tx.DeleteBucket([]byte(args[0]))
			if err == nil {
				found = true
			} else if err == bolt.ErrBucketNotFound {
				return nil
			}
			return err
		}
		b := tx.Bucket([]byte(args[0]))
		if b == nil {
			return nil
		}
		for i := 1; i < argsLen-1; i++ {
			b = b.Bucket([]byte(args[i]))
			if b == nil {
				return nil
			}
		}
		key := []byte(args[argsLen-1])
		err = b.DeleteBucket(key)
		if err == nil {
			found = true
		} else if err == bolt.ErrBucketNotFound || err == bolt.ErrIncompatibleValue {
			value := b.Get(key)
			if value == nil {
				return nil
			}
			found = true
			return b.Delete(key)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return found, nil
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
			if b == nil {
				return nil
			}
			for i := 1; i < argsLen-1; i++ {
				b = b.Bucket([]byte(args[i]))
				if b == nil {
					return nil
				}
			}
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
		return nil
	})
	return count, nil
}

func exists(args ...string) (res interface{}, err error) {
	argsLen := len(args)
	if argsLen < 1 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "exists")
	}
	var b *bolt.Bucket
	found := false
	err = DB.View(func(tx *bolt.Tx) error {
		b = tx.Bucket([]byte(args[0]))
		if b == nil {
			return nil
		}
		if argsLen == 1 {
			found = true
			return nil
		}
		for i := 1; i < argsLen-1; i++ {
			b = b.Bucket([]byte(args[i]))
			if b == nil {
				return nil
			}
		}
		lastWord := []byte(args[argsLen-1])
		if b.Bucket(lastWord) == nil && b.Get(lastWord) == nil {
			return nil
		}
		found = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	return found, nil
}

func get(args ...string) (res interface{}, err error) {
	argsLen := len(args)
	if argsLen < 2 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "get")
	}
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(args[0]))
		if b == nil {
			return nil
		}
		i := 1
		for ; i < argsLen-1; i++ {
			b = b.Bucket([]byte(args[i]))
			if b == nil {
				return nil
			}
		}
		res = b.Get([]byte(args[i]))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		res = ""
	}
	return
}

func set(args ...string) (res interface{}, err error) {
	argsLen := len(args)
	if argsLen < 3 {
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
		for i := 1; i < argsLen-2; i++ {
			subb := b.Bucket([]byte(args[i]))
			if subb == nil {
				subb, err = b.CreateBucket([]byte(args[i]))
				if err != nil {
					return err
				}
			}
			b = subb
		}
		return b.Put([]byte(args[argsLen-2]), []byte(args[argsLen-1]))
	})
	if err != nil {
		return nil, err
	}
	return true, nil
}

func buckets(args ...string) (res interface{}, err error) {
	argsLen := len(args)
	if argsLen < 1 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "buckets")
	}
	pattern, err := glob.Compile(args[argsLen-1])
	if err != nil {
		return
	}
	res = []string{}
	err = DB.View(func(tx *bolt.Tx) error {
		if argsLen > 1 {
			b := tx.Bucket([]byte(args[0]))
			if b == nil {
				return nil
			}
			for i := 1; i < argsLen-1; i++ {
				b = b.Bucket([]byte(args[i]))
				if b == nil {
					return nil
				}
			}
			b.ForEach(func(k, v []byte) error {
				name := string(k)
				if pattern.Match(name) && b.Bucket(k) != nil {
					res = append(res.([]string), name)
				}
				return nil
			})
		} else {
			tx.ForEach(func(bname []byte, b *bolt.Bucket) error {
				name := string(bname)
				if pattern.Match(name) {
					res = append(res.([]string), name)
				}
				return nil
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func keys(args ...string) (res interface{}, err error) {
	argsLen := len(args)
	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "keys")
	}
	pattern, err := glob.Compile(args[argsLen-1])
	if err != nil {
		return
	}
	res = []string{}
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(args[0]))
		if b == nil {
			return nil
		}
		for i := 1; i < argsLen-1; i++ {
			b = b.Bucket([]byte(args[i]))
			if b == nil {
				return nil
			}
		}
		b.ForEach(func(k, v []byte) error {
			key := string(k)
			if pattern.Match(key) && b.Bucket(k) == nil {
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
	argsLen := len(args)
	if argsLen < 2 {
		return nil, fmt.Errorf("wrong number of arguments for '%s' command", "keyvalues")
	}
	pattern, err := glob.Compile(args[argsLen-1])
	if err != nil {
		return
	}
	res = map[string]interface{}{}
	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(args[0]))
		if b == nil {
			return nil
		}
		for i := 1; i < argsLen-1; i++ {
			b = b.Bucket([]byte(args[i]))
			if b == nil {
				return nil
			}
		}
		b.ForEach(func(k, v []byte) error {
			key := string(k)
			if pattern.Match(key) && b.Bucket(k) == nil {
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
			info[typeField.Name] = valField.Int()
		}
	}
	val = reflect.ValueOf(stat.TxStats)
	info["TxStats"] = map[string]interface{}{}
	for i := 0; i < val.NumField(); i++ {
		valField := val.Field(i)
		typeField := val.Type().Field(i)
		if typeField.Type.Name() == "int" || typeField.Type.Name() == "Duration" {
			info["TxStats"].(map[string]interface{})[typeField.Name] = valField.Int()
		}
	}
	return info, nil
}

type cmd func(...string) (interface{}, error)

var CmdMap = map[string]cmd{
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
		case int64:
			formatted[i] = fmt.Sprintf(`%s%s) %v`, prefix, k, v)
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
	f, ok := CmdMap[strings.ToLower(cmd)]
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
		return fmt.Sprintf("\"%s\"", string(res))
	case string:
		return fmt.Sprintf("\"%s\"", res)
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
