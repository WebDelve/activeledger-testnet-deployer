package contracts

import (
	"fmt"

	"dynamicledger.com/testnet-deployer/helper"
	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

func (ch ContractHandler) OnboardContracts() {
	fmt.Println("Onboarding contracts...")

	data := []structs.ContractStore{}

	for _, c := range ch.Contracts {
		contractData := ch.onboardContract(c)
		data = append(data, contractData)
	}
}

func (ch *ContractHandler) onboardContract(contract structs.Contract) structs.ContractStore {
	fmt.Printf("\nOnboarding %s contract...\n", contract.Name)

	input := alsdk.DataWrapper{
		"version":   contract.Version,
		"namespace": ch.Setup.Namespace,
		"name":      contract.Name,
		"contract":  contract.Data,
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  ch.Setup.Identity,
		Contract:  "contract",
		Namespace: "default",
		Input:     input,
		Key:       ch.Setup.KeyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		helper.HandleError(err, fmt.Sprintf("Error building contract onboard transaction for contract %s", contract.Name))
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(*tx, ch.Setup.Conn)
	if err != nil {
		helper.HandleALError(err, resp, fmt.Sprintf("Error running onboarding transaction for contract %s", contract.Name))
	}

	fmt.Printf("%s contract onboarded.\n", contract.Name)

	data := structs.ContractStore{
		Name: contract.Name,
		ID:   resp.Streams.New[0].ID,
	}

	ch.labelContract(contract, data.ID)

	return data
}
