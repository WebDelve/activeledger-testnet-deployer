package contracts

import (
	"encoding/json"
	"fmt"

	"dynamicledger.com/testnet-deployer/helper"
	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type ContractUpdater struct {
	manifest          structs.ContractManifest
	contractData      []structs.Contract
	contractsToUpdate []structs.Contract
	setup             *structs.SetupData
	config            *structs.Config
	transactions      []contractUpdateTx
}

type contractUpdateTx struct {
	tx           alsdk.Transaction
	contractName string
}

func (ch *ContractHandler) getContractUpdater() ContractUpdater {

	u := ContractUpdater{
		manifest:          ch.Manifest,
		contractData:      ch.Contracts,
		contractsToUpdate: []structs.Contract{},
		setup:             ch.Setup,
		config:            ch.Config,
		transactions:      []contractUpdateTx{},
	}

	return u

}

func (cu *ContractUpdater) GetChangedContracts() []structs.Contract {
	return cu.contractsToUpdate
}

func (cu *ContractUpdater) Update() {
	fmt.Println("Updating changed contracts...")
	cu.findChangedContracts()

	cu.updateChangedContracts()

	fmt.Println("Contracts updated")
}

func (cu *ContractUpdater) findChangedContracts() {
	for _, c := range cu.contractData {
		if cu.contractChanged(c) {
			cu.contractsToUpdate = append(cu.contractsToUpdate, c)
		}
	}
}

func (cu *ContractUpdater) contractChanged(contract structs.Contract) bool {
	hash := getContractHash(contract)

	for _, c := range cu.manifest.Contracts {
		if c.Name == contract.Name {
			return c.Hash != hash
		}
	}

	return false
}

func requestVersionUpdate(currVersion string) string {
	fmt.Printf(
		"Current version is set to \"%s\", enter new version (leave blank to use existing): ",
		currVersion,
	)

	var newVersion string
	fmt.Scanln(&newVersion)

	var blank string
	if newVersion == blank {
		fmt.Printf("\nVersion unchanged, will use %s\n", currVersion)
		return currVersion
	}

	return newVersion
}

func (cu *ContractUpdater) updateChangedContracts() {
	contracts := cu.contractsToUpdate

	// Build the update transactions
	for i, contract := range contracts {

		newVersionNum := requestVersionUpdate(contract.Version)
		if newVersionNum != contract.Version {
			contract.Version = newVersionNum
			cu.contractsToUpdate[i].Version = contract.Version
		}

		cu.buildContractUpdateTx(contract)
	}

	// Run them all
	cu.runTransactions()
}

func (cu *ContractUpdater) buildContractUpdateTx(contract structs.Contract) {
	input := alsdk.DataWrapper{
		"version":   contract.Version,
		"namespace": cu.setup.Namespace,
		"name":      contract.Name,
		"contract":  contract.Data,
	}

	contractId := contract.Id

	txOpts := alsdk.TransactionOpts{
		StreamID:       cu.setup.Identity,
		OutputStreamID: alsdk.StreamID(contractId),
		Contract:       "contract",
		Namespace:      "default",
		Entry:          "update",
		Input:          input,
		Output:         alsdk.DataWrapper{},
		Key:            cu.setup.KeyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		helper.HandleError(err, fmt.Sprintf("Error building contract update transaction for contract %s", contract.Name))
	}

	tx := txHan.GetTransaction()

	txData := contractUpdateTx{
		tx:           tx,
		contractName: contract.Name,
	}

	cu.transactions = append(cu.transactions, txData)
}

func (cu *ContractUpdater) runTransactions() {
	for _, t := range cu.transactions {
		resp, err := alsdk.Send(t.tx, cu.setup.Conn)
		if err != nil {
			helper.HandleALError(err, resp, fmt.Sprintf("Error running update transaction %s", t.contractName))
		}
	}
}
