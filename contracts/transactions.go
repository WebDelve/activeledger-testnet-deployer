package contracts

import (
	"fmt"

	"dynamicledger.com/testnet-deployer/logging"
	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

func buildOnboardTx(
	contract structs.Contract,
	setup *structs.SetupData,
	logger *logging.Logger,
) alsdk.Transaction {

	input := alsdk.DataWrapper{
		"version":   contract.Version,
		"namespace": setup.Namespace,
		"name":      contract.Name,
		"contract":  contract.Data,
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  setup.Identity,
		Contract:  "contract",
		Namespace: "default",
		Input:     input,
		Key:       setup.KeyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		logger.Fatal(err, fmt.Sprintf("Error building contract onboard transaction for contract %s", contract.Name))
	}

	tx := txHan.GetTransaction()

	return tx
}

func buildUpdateTx(
	contract structs.Contract,
	setup *structs.SetupData,
	logger *logging.Logger,
) alsdk.Transaction {

	input := alsdk.DataWrapper{
		"version":   contract.Version,
		"namespace": setup.Namespace,
		"name":      contract.Name,
		"contract":  contract.Data,
	}

	contractId := contract.Id

	txOpts := alsdk.TransactionOpts{
		StreamID:       setup.Identity,
		OutputStreamID: alsdk.StreamID(contractId),
		Contract:       "contract",
		Namespace:      "default",
		Entry:          "update",
		Input:          input,
		Output:         alsdk.DataWrapper{},
		Key:            setup.KeyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		logger.Fatal(
			err,
			fmt.Sprintf("Error building contract update transaction for contract %s", contract.Name),
		)
	}

	tx := txHan.GetTransaction()

	return tx
}

func buildLabelTx(
	contractName string,
	contractId string,
	setup *structs.SetupData,
	logger *logging.Logger,
) alsdk.Transaction {

	input := alsdk.DataWrapper{
		"namespace": setup.Namespace,
		"contract":  contractId,
		"link":      contractName,
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  setup.Identity,
		Contract:  "contract",
		Namespace: "default",
		Entry:     "link",
		Input:     input,
		Key:       setup.KeyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		logger.Fatal(
			err,
			fmt.Sprintf("Error building contract link transaction for contract %s", contractName),
		)
	}

	tx := txHan.GetTransaction()

	return tx
}
