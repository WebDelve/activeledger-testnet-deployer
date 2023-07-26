package bootstrap

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"dynamicledger.com/testnet-deployer/contracts"
	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/helper"
	"dynamicledger.com/testnet-deployer/structs"

	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type Bootstrapper struct {
	config       *structs.Config
	setup        *structs.SetupData
	contractData []structs.ContractStore
}

func GetBootstrapper(config *structs.Config, setupData *structs.SetupData) Bootstrapper {
	return Bootstrapper{
		config,
		setupData,
		[]structs.ContractStore{},
	}
}

func (b *Bootstrapper) Bootstrap() {
	b.setup.Folder = getFolder(b.config.TestnetFolder)

	b.createTestnet()
	b.createIden()
	b.createNamespace()
	b.onboardSmartContracts()

	files.SaveSetupData(b.setup, b.contractData, b.config.SetupDataSaveFile)

	fmt.Printf("\n\nBootstrapping complete\n\n")
}

func getFolder(configuredFolder string) string {

	fmt.Printf("Input base folder name, leave blank for default (default=%s): ", configuredFolder)

	var folder string
	fmt.Scanln(&folder)

	var blank string
	if folder == blank {
		folder = configuredFolder
	}

	// Check if folder exists
	_, err := os.Stat(folder)

	// No error, folder exists, delete it if user agrees
	if err == nil {

		fmt.Printf("Folder \"%s\" exists, do you want to delete it? (Y/n): ", folder)

		var deleteFolder string
		fmt.Scanln(&deleteFolder)

		if deleteFolder == blank || deleteFolder == "Y" {
			if err := os.RemoveAll(folder); err != nil {
				helper.HandleError(err, "Error removing folder")
			}
		} else {
			fmt.Println("Won't delete folder, exiting instead...")
			os.Exit(0)
		}
	}

	// Folder doesn't exist, create it
	if errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(folder, 0755); err != nil {
			helper.HandleError(err, "Error creating base folder")
		}
	} else if err != nil {
		// Some other error happened, handle it
		helper.HandleError(err, "Error creating folder")
	}

	return folder
}

func (b *Bootstrapper) createTestnet() {

	cmd := exec.Command("activeledger", "--testnet")
	cmd.Dir = b.setup.Folder

	fmt.Println("Creating testnet")
	_, err := cmd.Output()
	if err != nil {
		helper.HandleError(err, "Error running activeledger --testnet, is activeledger installed?")
	}
	fmt.Println("Testnet created")

	fmt.Println("Leave this process running and start the testnet, navigate to the created folder and run 'node testnet'. When done return here and press enter to continue.")
	fmt.Scanln()
}

func (b *Bootstrapper) createIden() {
	fmt.Println("Creating and onboarding identity...")

	keyhandler, err := alsdk.GenerateRSA()
	if err != nil {
		helper.HandleError(err, "Error generating RSA key")
	}

	b.setup.KeyHandler = keyhandler

	pubpem := keyhandler.GetPublicPEM()

	input := alsdk.DataWrapper{
		"publicKey": pubpem,
		"type":      "rsa",
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  alsdk.StreamID(b.config.DefaultIdentity),
		Contract:  "onboard",
		Key:       keyhandler,
		Namespace: "default",
		SelfSign:  true,
		Input:     input,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		helper.HandleError(err, "Error building onboarding transaction")
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(tx, b.setup.Conn)
	if err != nil {
		helper.HandleALError(err, resp, "Error sending identity onboarding transaction")
	}

	var streamId string
	for _, v := range resp.Streams.New {
		if v.ID != streamId {
			streamId = v.ID
			break
		}
	}

	var blank string
	if streamId == blank {
		helper.HandleError(errors.New("Blank streamID"), "Identity transaction didn't error, but we got no stream ID")
	}

	b.setup.Identity = alsdk.StreamID(streamId)
}

func (b *Bootstrapper) createNamespace() {
	fmt.Println("Creating Namespace...")

	input := alsdk.DataWrapper{
		"namespace": b.setup.Namespace,
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  b.setup.Identity,
		Contract:  "namespace",
		Namespace: "default",
		Input:     input,
		Key:       b.setup.KeyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		helper.HandleError(err, "Error building namespace transaction")
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(tx, b.setup.Conn)
	if err != nil {
		helper.HandleALError(err, resp, "Error sending namespace transaction")
	}

	fmt.Printf("Namespace '%s' created\n", b.setup.Namespace)
}

func (b *Bootstrapper) onboardSmartContracts() {
	fmt.Println("Onboarding contracts...")

	conHan := contracts.SetupContractHandler(b.config, b.setup)
	conHan.OnboardContracts()
	b.contractData = conHan.GetContractData()

	fmt.Println("Onboarding complete")
}
