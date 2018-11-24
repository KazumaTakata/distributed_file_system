package main

import (
	"distributed_file_system/client/lib"
	"os"
)

func main() {

	command := os.Args[1]

	ExecuteCommand(command, os.Args[2:])

}

func ExecuteCommand(command string, others []string) {
	if command == "mkdir" {
		dirname := others[0]
		client.CreateDirMaster(dirname)
	} else if command == "touch" {
		filepath := others[0]
		req := client.CreateFileMaster(filepath)
		client.CreateFileChunk(req.Address, req.Filehash)
	}
}
