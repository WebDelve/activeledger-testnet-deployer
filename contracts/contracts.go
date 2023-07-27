package contracts

import (
	"fmt"

	"dynamicledger.com/testnet-deployer/logging"
	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type ContractHandler struct {
	Setup     *structs.SetupData
	Config    *structs.Config
	Manifest  structs.ContractManifest
	Store     []structs.ContractStore
	Contracts []structs.Contract
	Logger    *logging.Logger
}

func SetupContractHandler(
	config *structs.Config,
	setup *structs.SetupData,
	logger *logging.Logger,
) ContractHandler {

	ch := ContractHandler{
		Setup:    setup,
		Config:   config,
		Manifest: structs.ContractManifest{},
		Store:    []structs.ContractStore{},
		Logger:   logger,
	}

	ch.Logger.Info("Loading Contract manifest...")
	ch.getManifest()
	ch.Logger.Info("Manifest loaded")

	ch.Logger.Info("Loading contracts...")
	ch.loadContracts()

	return ch
}

func (ch *ContractHandler) GetContractData() []structs.ContractStore {
	return ch.Store
}

func (ch *ContractHandler) UpdateContracts() {
	updater := ch.getContractUpdater()
	updater.Update()

	changedContracts := updater.GetChangedContracts()

	newContracts := updater.GetNewContracts()
	ch.Contracts = append(ch.Contracts, newContracts...)

	newMetadata := updater.GetNewMetadata()

	for _, cMeta := range newMetadata {
		ch.Store = append(ch.Store, cMeta)

		ch.addIDToManifest(cMeta.Name, cMeta.ID)
		ch.updateOnboardedStatus(cMeta.ID)
	}

	ch.mergeInChangedContracts(changedContracts)
	ch.setHashes(false)
}

func (ch *ContractHandler) mergeInChangedContracts(changedContracts []structs.Contract) {
	for _, changed := range changedContracts {

		for i, contract := range ch.Contracts {

			if contract.Id == changed.Id {

				if contract.Version != changed.Version {
					ch.updateVersion(contract.Id, changed.Version)
				}

				ch.Contracts[i] = changed
				break
			}
		}
	}
}

func (ch *ContractHandler) labelContract(contract structs.Contract, contractId string) {
	ch.Logger.Info(fmt.Sprintf("\nLabeling contract %s..\n", contract.Name))

	tx := buildLabelTx(contract.Name, contractId, ch.Setup, ch.Logger)

	resp, err := alsdk.Send(tx, ch.Setup.Conn)
	if err != nil {
		ch.Logger.ActiveledgerError(err, resp, fmt.Sprintf("Error running contract lin transaction for contract %s", contract.Name))
	}

	ch.Logger.Info(fmt.Sprintf("Link created for contract %s.\n", contract.Name))
}
