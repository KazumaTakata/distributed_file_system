package master

import (
	"distributed_file_system/model"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type File struct {
	Name     string
	ChunkID  string
	Hashname string
}

type Directory struct {
	Name       string
	Directorys *map[string]*Directory
	Files      *map[string]*File
}

type Filetree struct {
	Root *Directory
}

type MetaData struct {
	FileTree  Filetree
	ChunkHash *map[string]string
}

type Master struct {
	Metadata MetaData
}

func StartServer() {
	fmt.Printf("master server start...")

	master := new(Master)

	file := "./data/metadata"
	master.Metadata = MetaData{ChunkHash: &map[string]string{}, FileTree: Filetree{Root: &Directory{Directorys: &map[string]*Directory{}, Name: "root", Files: &map[string]*File{}}}}

	_, err := os.Stat(file)
	if err == nil {
		log.Printf("file %s exists", file)
		readGob(master)
	} else if os.IsNotExist(err) {
		log.Printf("file %s not exists", file)

	} else {
		log.Printf("file %s stat error: %v", file, err)
	}

	// dummy chunk server
	(*master.Metadata.ChunkHash)["test_chunk"] = "127.0.0.1:1235"

	rpc.Register(master)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}

func (m *Master) getDir(path []string) *Directory {

	dir := m.Metadata.FileTree.Root

	for _, p := range path {
		dir = (*dir.Directorys)[string(p)]
		if dir == nil {
			break
		}
	}
	return dir
}

func (m *Master) CreateDir(args *model.Args, res *model.Res) error {
	defer writeGob(m)
	// dirs := m.Metadata.FileTree.Cwd
	path := args.Path
	pathlist := strings.Split(path, "/")
	dirpath := pathlist[1 : len(pathlist)-1]
	dirname := pathlist[len(pathlist)-1]

	curdir := m.getDir(dirpath)

	if curdir == nil {
		res.Success = false
		res.Message = fmt.Sprintf("Directory %s not exists.", strings.Join(pathlist[0:len(pathlist)-1], "/"))
		return nil
	}

	if _, ok := (*curdir.Directorys)[dirname]; ok {
		res.Success = false
		res.Message = fmt.Sprintf("Directory %s already exists.", dirname)
		return nil
	}
	(*curdir.Directorys)[dirname] = &Directory{Directorys: &map[string]*Directory{}, Name: dirname, Files: &map[string]*File{}}

	res.Success = true
	res.Message = fmt.Sprintf("Directory %s created.", dirname)

	return nil
}

// func (m *Master) ChangeDir(args *model.Args, res *model.Res) error {
// 	defer writeGob(m)
// 	dirs := m.Metadata.FileTree.Cwd
// 	path := args.Path

// 	if val, ok := (*dirs.Directorys)[path]; ok {
// 		m.Metadata.FileTree.Cwd = val
// 		res.Success = true
// 		res.Message = fmt.Sprintf("Change cwd to %s.", path)
// 		return nil
// 	}

// 	res.Success = false
// 	res.Message = fmt.Sprintf("Directory %s not exists.", path)
// 	return nil
// }

func (m *Master) CreateFile(args *model.Args, res *model.ResFile) error {
	path := args.Path
	pathlist := strings.Split(path, "/")
	dirpath := pathlist[1 : len(pathlist)-1]
	filename := pathlist[len(pathlist)-1]

	dirs := m.getDir(dirpath)

	chunkhash := m.Metadata.ChunkHash

	chunkIDs := make([]string, 0, len(*chunkhash))
	for k := range *chunkhash {
		chunkIDs = append(chunkIDs, k)
	}

	u1 := uuid.Must(uuid.NewV4())
	hashname := u1.String()

	chunkID := chunkIDs[0]

	file := File{Name: filename, ChunkID: chunkID, Hashname: hashname}

	(*dirs.Files)[filename] = &file

	res.Filehash = u1.String()
	res.Address = (*chunkhash)[chunkID]

	return nil
}

func (m *Master) AddNewChunkServer(args *model.ChunkArg, res *model.Res) error {

	chunkID := args.ID
	chunkAdd := args.Address

	(*m.Metadata.ChunkHash)[chunkID] = chunkAdd

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
