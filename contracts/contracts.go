package contracts

import (
	"fmt"

	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/helper"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type ContractHandler struct {
	setup     *helper.SetupData
	config    *helper.Config
	manifest  contractManifest
	store     []files.ContractStore
	contracts []Contract
}

type Contract struct {
	name    string
	data    string
	version string
}

type contractManifest struct {
	Contracts []contractMetadata `json:"contracts"`
}

type contractMetadata struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
	Exclude bool   `json:"exclude"`
	Hash    string `json:"hash"`
}

func SetupContractHandler(config *helper.Config, setup *helper.SetupData) ContractHandler {
	ch := ContractHandler{
		setup:    setup,
		config:   config,
		manifest: contractManifest{},
		store:    []files.ContractStore{},
	}

	fmt.Println("Loading contracts...")
	ch.loadContracts()

	return ch
}

func (ch *ContractHandler) UpdateContracts() {
	updater := ch.getContractUpdater()
	updater.Update()
}

func (ch *ContractHandler) GetContractData() []files.ContractStore {
	return ch.store
}

func (ch *ContractHandler) labelContract(contract Contract, contractId string) {
	fmt.Printf("\nLabeling contract %s..\n", contract.name)

	input := alsdk.DataWrapper{
		"namespace": ch.setup.Namespace,
		"contract":  contractId,
		"link":      contract.name,
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  ch.setup.Identity,
		Contract:  "contract",
		Namespace: "default",
		Entry:     "link",
		Input:     input,
		Key:       ch.setup.KeyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		helper.HandleError(err, fmt.Sprintf("Error building contract link transaction for contract %s", contract.name))
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(*tx, ch.setup.Conn)
	if err != nil {
		helper.HandleALError(err, resp, fmt.Sprintf("Error running contract lin transaction for contract %s", contract.name))
	}

	fmt.Printf("Link created for contract %s.\n", contract.name)
}
