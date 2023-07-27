package structs

import (
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type SetupData struct {
	Folder     string
	Identity   alsdk.StreamID
	Namespace  string
	KeyHandler alsdk.KeyHandler
	Conn       alsdk.Connection
	Contracts  []ContractStore
}
