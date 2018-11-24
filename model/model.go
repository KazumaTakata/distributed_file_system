package model

type Args struct {
	Path string
}

type ChunkArg struct {
	ID      string
	Address string
}

type Res struct {
	Success bool
	Message string
}

type ResFile struct {
	Filehash string
	Address  string
}

type FileObj struct {
	Filehash string
	Body     []byte
}
