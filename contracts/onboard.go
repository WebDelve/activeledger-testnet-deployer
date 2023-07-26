package contracts

import (
	"fmt"

	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/helper"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

func (ch ContractHandler) OnboardContracts() {
	fmt.Println("Onboarding contracts...")

	data := []files.ContractStore{}

	for _, c := range ch.contracts {
		contractData := ch.onboardContract(c)
		data = append(data, contractData)

		// update ch data
	}
}

func (ch *ContractHandler) onboardContract(contract Contract) files.ContractStore {
	fmt.Printf("\nOnboarding %s contract...\n", contract.name)

	input := alsdk.DataWrapper{
		"version":   contract.version,
		"namespace": ch.setup.Namespace,
		"name":      contract.name,
		"contract":  contract.data,
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  ch.setup.Identity,
		Contract:  "contract",
		Namespace: "default",
		Input:     input,
		Key:       ch.setup.KeyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		helper.HandleError(err, fmt.Sprintf("Error building contract onboard transaction for contract %s", contract.name))
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(*tx, ch.setup.Conn)
	if err != nil {
		helper.HandleALError(err, resp, fmt.Sprintf("Error running onboarding transaction for contract %s", contract.name))
	}

	fmt.Printf("%s contract onboarded.\n", contract.name)

	data := files.ContractStore{
		Name: contract.name,
		ID:   resp.Streams.New[0].ID,
	}

	ch.labelContract(contract, data.ID)

	return data
}
