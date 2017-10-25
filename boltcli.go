package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	bolt "github.com/coreos/bbolt"
)

var (
	// DB is the global boltdb instance which will be inited in the beginning.
	DB *bolt.DB
	// DbPath is the path of given db file
	DbPath string

	shouldPrintVersion = flag.Bool("version", false, "Output version and exit.")
	version            = "1.0.0"

	scriptPath = flag.String("e", "", "Eval the Lua script in given path")
)

func initDB(dbPath string) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatalf("Could not open %s: %v", dbPath, err)
	}
	DB = db
	DbPath = dbPath
}

func printVersion() {
	fmt.Printf("boltcli %s\n", version)
}

func main() {
	flag.Parse()
	if *shouldPrintVersion {
		printVersion()
		os.Exit(0)
	}
	if flag.NArg() < 1 {
		log.Fatalf("database filename is required.")
	}
	initDB(flag.Arg(0))
	defer DB.Close()
	if *scriptPath != "" {
		err := StartScript(*scriptPath)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		StartCli()
	}
}
