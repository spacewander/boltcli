package main

import (
	"strings"
)

var CmdHelp = map[string][2]string{
	"del": [2]string{
		"[bucket ...] bucket/key",
		strings.Join([]string{
			"Deletes the key in the specified bucket, and returns true.",
			"If no key is given, deletes the bucket, and return true.",
			"If the bucket/key does not exists, returns false.",
		}, "\n"),
	},
	"delglob": [2]string{
		"[bucket ...] bucket/key-pattern",
		strings.Join([]string{
			"Deletes the buckets/keys matching the given glob pattern, and returns the number of items deleted.",
			"If bucket does not exist, returns 0",
		}, "\n"),
	},
	"exists": [2]string{
		"[bucket ...] bucket/key",
		strings.Join([]string{
			"Checks if a given bucket/key exists.",
			"Returns true if it does, otherwise false.",
		}, "\n"),
	},
	"get": [2]string{
		"[bucket ...] bucket key",
		strings.Join([]string{
			"Returns the value of the given key in the specified bucket.",
			"Returns an empty string if the bucket or key does not exist.",
		}, "\n"),
	},
	"help": [2]string{
		"command",
		strings.Join([]string{
			"Shows the help output for the given command.",
		}, "\n"),
	},
	"set": [2]string{
		"[bucket ...] bucket key value",
		strings.Join([]string{
			"Sets the value of the given key in the specified bucket, and returns true.",
			"If the bucket does not exist it will be created.",
		}, "\n"),
	},
	"buckets": [2]string{
		"[bucket ...] bucket-pattern",
		strings.Join([]string{
			"Lists all buckets matching the given glob pattern.",
		}, "\n"),
	},
	"keys": [2]string{
		"[bucket ...] bucket key-pattern",
		strings.Join([]string{
			"Lists all keys in the specified bucket matching the given glob pattern.",
		}, "\n"),
	},
	"keyvalues": [2]string{
		"[bucket ...] bucket key-pattern",
		strings.Join([]string{
			"Lists all keys and their associated values in the specified bucket matching the given glob pattern.",
		}, "\n"),
	},
	"stats": [2]string{
		"",
		strings.Join([]string{
			"Returns the result of bolt.DB.Stats()",
		}, "\n"),
	},
}
