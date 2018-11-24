package chunk

import (
	"distributed_file_system/model"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Chunk struct {
	ServerID string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func StartServer() {
	fmt.Printf("chunk server start...")

	chunk := new(Chunk)

	file := "./data/metadata"

	_, err := os.Stat(file)
	if err == nil {
		log.Printf("file %s exists", file)
		readGob(chunk)
	} else if os.IsNotExist(err) {
		log.Printf("file %s not exists", file)
	} else {
		log.Printf("file %s stat error: %v", file, err)
	}

	// chunk.joinMessage(address)

	rpc.Register(chunk)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1235")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}

func (m *Chunk) joinMessage(address string) error {
	defer writeGob(m)

	serverID := m.ServerID

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := &model.ChunkArg{ID: serverID, Address: address}
	rep := &model.Res{}
	err = client.Call("Master.AddNewChunkServer", args, &rep)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	return nil
}

func (m *Chunk) CreateFile(args *model.FileObj, res *model.Res) error {

	filehash := args.Filehash
	// filebody := args.Body

	err := ioutil.WriteFile("data/"+filehash, []byte{}, 0644)
	check(err)

	return nil
}

func writeGob(object interface{}) error {
	filePath := "./data/metadata"
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(object interface{}) error {
	filePath := "./data/metadata"
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
