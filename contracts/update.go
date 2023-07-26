package contracts

import (
	"fmt"

	"dynamicledger.com/testnet-deployer/config"
	"dynamicledger.com/testnet-deployer/helper"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type ContractUpdater struct {
	manifest          contractManifest
	contractData      []Contract
	contractsToUpdate []Contract
	setup             *helper.SetupData
	config            *config.Config
	transactions      []contractUpdateTx
}

type contractUpdateTx struct {
	tx           *alsdk.Transaction
	contractName string
}

func (ch *ContractHandler) getContractUpdater() ContractUpdater {

	u := ContractUpdater{
		manifest:          ch.manifest,
		contractData:      ch.contracts,
		contractsToUpdate: []Contract{},
		setup:             ch.setup,
		config:            ch.config,
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

func (cu *ContractUpdater) contractChanged(contract Contract) bool {
	hash := getContractHash(contract)

	for _, c := range cu.manifest.Contracts {
		if c.Name == contract.name {
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

func (cu *ContractUpdater) buildContractUpdateTx(contract Contract) {
	input := alsdk.DataWrapper{
		"version":   contract.version,
		"namespace": cu.setup.Namespace,
		"name":      contract.name,
		"contract":  contract.data,
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  cu.setup.Identity,
		Contract:  "contract",
		Namespace: "default",
		Entry:     "update",
		Input:     input,
		Key:       cu.setup.KeyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		helper.HandleError(err, fmt.Sprintf("Error building contract update transaction for contract %s", contract.name))
	}

	tx := txHan.GetTransaction()

	txData := contractUpdateTx{
		tx:           tx,
		contractName: contract.name,
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
