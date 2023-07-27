package contracts

import (
	"fmt"

	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

func (ch *ContractHandler) OnboardContracts() {
	ch.Logger.Info("Onboarding contracts...")

	data := []structs.ContractStore{}

	for _, c := range ch.Contracts {
		contractData := ch.onboardContract(c)
		data = append(data, contractData)
	}

	ch.setHashes(false)

	ch.Store = data
}

func (ch *ContractHandler) onboardContract(contract structs.Contract) structs.ContractStore {
	ch.Logger.Info(fmt.Sprintf("Onboarding %s contract...", contract.Name))

	tx := buildOnboardTx(contract, ch.Setup, ch.Logger)

	resp, err := alsdk.Send(tx, ch.Setup.Conn)
	if err != nil {
		ch.Logger.ActiveledgerError(err, resp, fmt.Sprintf("Error running onboarding transaction for contract %s", contract.Name))
	}

	ch.Logger.Info(fmt.Sprintf("%s contract onboarded.", contract.Name))

	data := structs.ContractStore{
		Name: contract.Name,
		ID:   resp.Streams.New[0].ID,
	}

	ch.addIDToManifest(contract.Name, data.ID)
	ch.updateOnboardedStatus(data.ID)

	ch.labelContract(contract, data.ID)

	return data
}
