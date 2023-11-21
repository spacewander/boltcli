![screenshot](./screenshot.png)

[![Travis](https://travis-ci.org/spacewander/boltcli.svg?branch=master)](https://travis-ci.org/spacewander/boltcli)
[![GoReportCard](http://goreportcard.com/badge/spacewander/boltcli)](http://goreportcard.com/report/spacewander/boltcli)
[![codecov.io](https://codecov.io/github/spacewander/boltcli/coverage.svg?branch=master)](https://codecov.io/github/spacewander/boltcli?branch=master)
[![license](https://img.shields.io/badge/License-GPLv3-green.svg)](https://github.com/spacewander/boltcli/blob/master/LICENSE)

## Feature

* Support CRUD commands on the bucket/key in repl-like command line. You can consider it as `redis-cli` for boltdb.
* You can eval Lua script with given database. It makes maintaining boltdb more easily.

## Usage

`boltcli  [-e script] /path/to/db`

## Commands

Documentation for commands is available with the built-in help command:
```
/tmp/test.db> help
Commands: buckets, del, delglob, exists, get, help, keys, keyvalues, set, stats
/tmp/test.db> help help
Command: help command

Shows the help output for the given command.

```

## Lua support

You could run a lua script on specific database like this: `boltcli -e your.lua db_path`.
`boltcli` provides a couple of API within the global variable `bolt`. For example:
```lua
bolt.get("bucket", "key") -- return the value of `key` as a lua string
-- is equal to > get bucket key in the command line
```

See [test.lua](./test.lua) as a concrete example.
