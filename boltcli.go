package main

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
)

var (
	// DB is the global boltdb instance which will be inited in the beginning.
	DB *bolt.DB
)

func initDB(dbPath string) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatalf("Could not open %s: %v", dbPath, err)
	}
	DB = db
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("database filename is required.")
	}
	initDB(os.Args[1])
	DB.Close()
}
