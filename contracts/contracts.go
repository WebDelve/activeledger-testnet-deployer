package contracts

import (
	"fmt"

	"dynamicledger.com/testnet-deployer/helper"
	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type ContractHandler struct {
	Setup     *structs.SetupData
	Config    *structs.Config
	Manifest  structs.ContractManifest
	Store     []structs.ContractStore
	Contracts []structs.Contract
}

func SetupContractHandler(config *structs.Config, setup *structs.SetupData) ContractHandler {
	ch := ContractHandler{
		Setup:    setup,
		Config:   config,
		Manifest: structs.ContractManifest{},
		Store:    []structs.ContractStore{},
	}

	fmt.Println("Loading contracts...")
	ch.loadContracts()

	return ch
}

func (ch *ContractHandler) UpdateContracts() {
	updater := ch.getContractUpdater()
	updater.Update()

	changedContracts := updater.GetChangedContracts()

	ch.mergeInChangedContracts(changedContracts)
	ch.updateContractHashes()
}

func (ch *ContractHandler) mergeInChangedContracts(changedContracts []structs.Contract) {
	for _, changed := range changedContracts {
		for i, contract := range ch.Contracts {
			if contract.Id == changed.Id {
				ch.Contracts[i] = changed
				break
			}
		}
	}
}

func (ch *ContractHandler) GetContractData() []structs.ContractStore {
	return ch.Store
}

func (ch *ContractHandler) labelContract(contract structs.Contract, contractId string) {
	fmt.Printf("\nLabeling contract %s..\n", contract.Name)

	input := alsdk.DataWrapper{
		"namespace": ch.Setup.Namespace,
		"contract":  contractId,
		"link":      contract.Name,
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  ch.Setup.Identity,
		Contract:  "contract",
		Namespace: "default",
		Entry:     "link",
		Input:     input,
		Key:       ch.Setup.KeyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		helper.HandleError(err, fmt.Sprintf("Error building contract link transaction for contract %s", contract.Name))
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(*tx, ch.Setup.Conn)
	if err != nil {
		helper.HandleALError(err, resp, fmt.Sprintf("Error running contract lin transaction for contract %s", contract.Name))
	}

	fmt.Printf("Link created for contract %s.\n", contract.Name)
}
