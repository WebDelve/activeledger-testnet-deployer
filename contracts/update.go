package contracts

import (
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
	tx           *alsdk.Transaction
	contractName string
}

func (ch *ContractHandler) getContractUpdater() ContractUpdater {

	u := ContractUpdater{
		manifest:          ch.Manifest,
		contractData:      ch.Contracts,
		contractsToUpdate: []structs.Contract{},
		setup:             ch.Setup,
		config:            ch.Config,
	}

	return u

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

func (cu *ContractUpdater) updateChangedContracts() {
	contracts := cu.contractsToUpdate

	// Build the update transactions
	for _, contract := range contracts {
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

	contractId := alsdk.StreamID(contract.Id)
	output := alsdk.DataWrapper{}
	output[contractId] = "{}"

	txOpts := alsdk.TransactionOpts{
		StreamID:  cu.setup.Identity,
		Contract:  "contract",
		Namespace: "default",
		Entry:     "update",
		Input:     input,
		Output:    output,
		Key:       cu.setup.KeyHandler,
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
		resp, err := alsdk.Send(*t.tx, cu.setup.Conn)
		if err != nil {
			helper.HandleALError(err, resp, fmt.Sprintf("Error running update transaction %s", t.contractName))
		}
	}
}
