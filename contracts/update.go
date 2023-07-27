package contracts

import (
	"fmt"
	"os"

	"dynamicledger.com/testnet-deployer/logging"
	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type ContractUpdater struct {
	manifest          structs.ContractManifest
	contractData      []structs.Contract
	contractsToUpdate []structs.Contract
	newContractsMeta  []structs.ContractStore
	contractsToUpload []structs.Contract
	setup             *structs.SetupData
	config            *structs.Config
	transactions      []contractTx
	logger            *logging.Logger
}

type contractTx struct {
	tx           alsdk.Transaction
	contractName string
	update       bool
}

func (ch *ContractHandler) getContractUpdater() ContractUpdater {

	u := ContractUpdater{
		manifest:          ch.Manifest,
		contractData:      ch.Contracts,
		contractsToUpdate: []structs.Contract{},
		newContractsMeta:  []structs.ContractStore{},
		contractsToUpload: []structs.Contract{},
		setup:             ch.Setup,
		config:            ch.Config,
		transactions:      []contractTx{},
		logger:            ch.Logger,
	}

	return u

}

func (cu *ContractUpdater) GetChangedContracts() []structs.Contract {
	return cu.contractsToUpdate
}

func (cu *ContractUpdater) GetNewContracts() []structs.Contract {
	return cu.contractsToUpload
}

func (cu *ContractUpdater) GetNewMetadata() []structs.ContractStore {
	return cu.newContractsMeta
}

func (cu *ContractUpdater) Update() {
	cu.logger.Info("Updating changed contracts...")
	cu.logger.Info("Finding changed contracts...")
	cu.findChangedContracts()
	cu.logger.Info("Finding new contracts...")
	cu.findNewContracts()

	if len(cu.contractsToUpload) <= 0 && len(cu.contractsToUpdate) <= 0 {
		cu.logger.Info("No contracts to upload or update, quitting...")
		os.Exit(0)
	}

	cu.logger.Info("Contracts found, creating new transactions...")
	cu.createUpdateTxs()
	cu.createNewContractTxs()

	// Run them all
	cu.logger.Info("Transactions created, running them...")
	cu.runTransactions()

	cu.logger.Info("Contracts updated/uploaded")
}

func (cu *ContractUpdater) findNewContracts() {
	for _, c := range cu.contractData {
		if cu.contractNotOnboarded(c) {
			cu.contractsToUpload = append(cu.contractsToUpload, c)
		}
	}
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
		if c.Name != contract.Name {
			continue
		}

		// Skip excluded or not onboarded contracts
		if c.Exclude || !c.Onboarded {
			continue
		}

		isChanged := c.Hash != hash
		if isChanged {
			cu.logger.Info("Found updated contract: " + c.Name)
		}

		return true
	}

	return false
}

func (cu *ContractUpdater) contractNotOnboarded(contract structs.Contract) bool {
	for _, c := range cu.manifest.Contracts {
		if c.Name != contract.Name {
			continue
		}

		// skip excluded and onboarded
		if c.Exclude || c.Onboarded {
			continue
		}

		cu.logger.Info("Found new contract: " + c.Name)

		return true
	}

	return false
}

func (cu *ContractUpdater) requestVersionUpdate(currVersion string) string {
	newVersion := cu.logger.GetUserInput(fmt.Sprintf(
		"Current version is set to \"%s\", enter new version (leave blank to use existing): ",
		currVersion,
	))

	var blank string
	if newVersion == blank {
		cu.logger.Info(fmt.Sprintf("Version unchanged, will use %s", currVersion))
		return currVersion
	}

	return newVersion
}

func (cu *ContractUpdater) createUpdateTxs() {
	contracts := cu.contractsToUpdate

	// Build the update transactions
	for i, contract := range contracts {

		newVersionNum := cu.requestVersionUpdate(contract.Version)
		if newVersionNum != contract.Version {
			contract.Version = newVersionNum
			cu.contractsToUpdate[i].Version = contract.Version
		}

		cu.buildContractUpdateTx(contract)
	}

}

func (cu *ContractUpdater) createNewContractTxs() {
	contracts := cu.contractsToUpload

	for _, contract := range contracts {
		cu.buildNewContractTx(contract)
	}
}

func (cu *ContractUpdater) buildNewContractTx(contract structs.Contract) {
	tx := buildOnboardTx(contract, cu.setup, cu.logger)
	txData := contractTx{
		tx:           tx,
		contractName: contract.Name,
		update:       false,
	}

	cu.transactions = append(cu.transactions, txData)
}

func (cu *ContractUpdater) labelContract(contractName string, contractId string) {

	cu.logger.Info(fmt.Sprintf("Labeling contract %s..", contractName))

	tx := buildLabelTx(contractName, contractId, cu.setup, cu.logger)

	resp, err := alsdk.Send(tx, cu.setup.Conn)
	if err != nil {
		cu.logger.ActiveledgerError(
			err,
			resp,
			fmt.Sprintf("Error running contract link transaction for contract %s", contractName),
		)
	}

	cu.logger.Info(fmt.Sprintf("Link created for contract %s.", contractName))

}

func (cu *ContractUpdater) buildContractUpdateTx(contract structs.Contract) {
	tx := buildUpdateTx(contract, cu.setup, cu.logger)
	txData := contractTx{
		tx:           tx,
		contractName: contract.Name,
		update:       true,
	}

	cu.transactions = append(cu.transactions, txData)
}

func (cu *ContractUpdater) runTransactions() {
	for _, t := range cu.transactions {
		resp, err := alsdk.Send(t.tx, cu.setup.Conn)
		if err != nil {
			cu.logger.ActiveledgerError(err, resp, fmt.Sprintf("Error running update transaction %s", t.contractName))
		}

		if !t.update {
			id := resp.Streams.New[0].ID
			data := structs.ContractStore{
				Name: t.contractName,
				ID:   id,
			}

			cu.newContractsMeta = append(cu.newContractsMeta, data)
			cu.labelContract(t.contractName, id)
		}

	}
}
