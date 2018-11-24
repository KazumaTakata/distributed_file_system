package main

import (
	"distributed_file_system/client"
	"distributed_file_system/master"
	"os"
)

func main() {

	kind := os.Args[1]
	if kind == "master" {
		master.StartServer()
	} else if kind == "client" {
		command := os.Args[2]
		if command == "mkdir" {
			client.CreateDirMaster()
		}

	}

}
