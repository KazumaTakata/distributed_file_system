package client

import (
	"fmt"
	"log"
	"net/rpc"

	"distributed_file_system/model"
)

func CreateDirMaster(dirname string) {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := &model.Args{Path: dirname}
	rep := &model.Res{}
	err = client.Call("Master.CreateDir", args, &rep)
	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Printf(rep.Message)
}

// func ChangeDirMaster(dirname string) {
// 	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
// 	if err != nil {
// 		log.Fatal("dialing:", err)
// 	}
// 	args := &model.Args{Path: dirname}
// 	rep := &model.Res{}
// 	err = client.Call("Master.ChangeDir", args, &rep)
// 	if err != nil {
// 		log.Fatal("arith error:", err)
// 	}

// 	fmt.Printf(rep.Message)
// }

func CreateFileMaster(filepath string) *model.ResFile {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := &model.Args{Path: filepath}
	rep := &model.ResFile{}
	err = client.Call("Master.CreateFile", args, &rep)
	if err != nil {
		log.Fatal("arith error:", err)
	}

	return rep
}

func CreateFileChunk(address string, filehash string) {

	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args2 := &model.FileObj{Filehash: filehash}
	rep2 := &model.Res{}
	err = client.Call("Chunk.CreateFile", args2, &rep2)
	if err != nil {
		log.Fatal("arith error:", err)
	}

}
