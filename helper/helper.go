package helper

import (
	"log"

	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type Config struct {
	DefaultIdentity   string `json:"identity"`
	DefaultNamespace  string `json:"namespace"`
	ContractFolder    string `json:"contractDir"`
	ContractManifest  string `json:"contractManifest"`
	SetupDataSaveFile string `json:"setupDataSaveFile"`
	TestnetFolder     string `json:"testnetFolder"`
}

type SetupData struct {
	Folder     string
	Identity   alsdk.StreamID
	Namespace  string
	KeyHandler alsdk.KeyHandler
	Conn       alsdk.Connection
}

func HandleError(e error, note string) {
	log.Println(note)
	log.Fatalln(e)
}

func HandleALError(e error, resp alsdk.Response, note string) {
	if len(resp.Summary.Errors) > 0 {
		for i, e := range resp.Summary.Errors {
			log.Printf("Error %d: %s\n\n", i, e)
		}
	}

	HandleError(e, note)
}
