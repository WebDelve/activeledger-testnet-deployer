package bootstrap

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"dynamicledger.com/testnet-deployer/contracts"
	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/logging"
	"dynamicledger.com/testnet-deployer/structs"

	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type Bootstrapper struct {
	config       *structs.Config
	setup        *structs.SetupData
	contractData []structs.ContractStore
	logger       *logging.Logger
	fileHandler  files.FileHandler
}

func GetBootstrapper(
	config *structs.Config,
	setupData *structs.SetupData,
	logger *logging.Logger,
) Bootstrapper {

	fHan := files.GetFileHandler(logger)

	return Bootstrapper{
		config,
		setupData,
		[]structs.ContractStore{},
		logger,
		fHan,
	}

}

func (b *Bootstrapper) Bootstrap() {
	b.setup.Folder = b.getFolder(b.config.TestnetFolder)

	b.createTestnet()
	b.createIden()
	b.createNamespace()
	b.onboardSmartContracts()

	b.fileHandler.SaveSetupData(b.setup, b.contractData, b.config.SetupDataSaveFile)

	b.logger.Info("\n\nBootstrapping complete\n\n")
}

func (b *Bootstrapper) getFolder(configuredFolder string) string {

	folder := b.logger.GetUserInput(
		fmt.Sprintf(
			"Input base folder name, leave blank for default (default = %s): ",
			configuredFolder,
		),
	)

	var blank string
	if folder == blank {
		folder = configuredFolder
	}

	// Check if folder exists
	_, err := os.Stat(folder)

	// No error, folder exists, delete it if user agrees
	if err == nil {

		deleteFolder := b.logger.GetUserInput(
			fmt.Sprintf(
				"Folder \"%s\" exists, do you want to delete it? (Y/n): ",
				folder,
			),
		)

		if deleteFolder == blank || deleteFolder == "Y" {
			if err := os.RemoveAll(folder); err != nil {
				b.logger.Fatal(err, "Error removing folder")
			}
		} else {
			b.logger.Info("Won't delete folder, exiting instead...")
			os.Exit(0)
		}
	}

	// Check if the error is NOT that the folder doesn't exist
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// Some other error happened, handle it
		b.logger.Fatal(err, "Error checking if folder exists")
	}

	// Create folder, we've either deleted it or it never existed
	if err := os.Mkdir(folder, 0755); err != nil {
		b.logger.Fatal(err, "Error creating base folder")
	}

	return folder
}

func (b *Bootstrapper) createTestnet() {

	cmd := exec.Command("activeledger", "--testnet")
	cmd.Dir = b.setup.Folder

	b.logger.Info("Creating testnet")
	_, err := cmd.Output()
	if err != nil {
		b.logger.Fatal(err, "Error running activeledger --testnet, is activeledger installed?")
	}
	b.logger.Info("Testnet created")

	b.logger.Info("\nLeave this process running and start the testnet, navigate to the created folder and run 'node testnet'. When done return here and press enter to continue.")
	fmt.Scanln()
}

func (b *Bootstrapper) createIden() {
	b.logger.Info("Creating and onboarding identity...")

	keyhandler, err := alsdk.GenerateRSA()
	if err != nil {
		b.logger.Fatal(err, "Error generating RSA key")
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
		b.logger.Fatal(err, "Error building onboarding transaction")
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(tx, b.setup.Conn)
	if err != nil {
		b.logger.ActiveledgerError(err, resp, "Error sending identity onboarding transaction")
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
		b.logger.Fatal(errors.New("Blank streamID"), "Identity transaction didn't error, but we got no stream ID")
	}

	b.setup.Identity = alsdk.StreamID(streamId)
}

func (b *Bootstrapper) createNamespace() {
	b.logger.Info("Creating Namespace...")

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
		b.logger.Fatal(err, "Error building namespace transaction")
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(tx, b.setup.Conn)
	if err != nil {
		b.logger.ActiveledgerError(err, resp, "Error sending namespace transaction")
	}

	b.logger.Info(fmt.Sprintf("Namespace '%s' created\n", b.setup.Namespace))
}

func (b *Bootstrapper) onboardSmartContracts() {
	b.logger.Info("Onboarding contracts...")

	conHan := contracts.SetupContractHandler(b.config, b.setup, b.logger)
	conHan.OnboardContracts()
	b.contractData = conHan.GetContractData()

	b.logger.Info("Onboarding complete")
}
