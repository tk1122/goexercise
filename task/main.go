package main

import (
	"github.com/mitchellh/go-homedir"
	"goercises/task/cmd"
	"goercises/task/db"
	"log"
	"path/filepath"
)

const dbName = "task.db"

func main() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("cannot get home dir: %v\n", err)
	}
	dbPath := filepath.Join(home, dbName)
	err = db.Init(dbPath)
	if err != nil {
		log.Fatalf("cannot connect to db %v: %v\n", dbName, err)
	}
	err = cmd.RootCmd.Execute()
	if err != nil {
		log.Fatalf("cannot execute root command: %v\n", err)
	}
}
