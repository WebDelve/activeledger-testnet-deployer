package main

import (
	"flag"
	"fmt"
	"os"

	"dynamicledger.com/testnet-deployer/bootstrap"
	"dynamicledger.com/testnet-deployer/config"
	"dynamicledger.com/testnet-deployer/contracts"
	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/logging"
	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type CLIFlags struct {
	setupTestnet    bool
	updateContracts bool
	verboseLogging  bool
	headlessMode    bool
}

func main() {
	logger := logging.CreateLogger()

	config := config.LoadConfig(&logger)

	// Set after loading config so we can use it for errors loading config
	logger.SetConfig(config)

	flags := handleFlags()

	setupData := structs.SetupData{}

	setupData.Conn = alsdk.Connection{
		Protocol: alsdk.HTTP,
		Url:      "localhost",
		Port:     "5260",
	}

	setupData.Namespace = config.DefaultNamespace

	if flags.setupTestnet {

		bs := bootstrap.GetBootstrapper(config, &setupData, &logger)
		bs.Bootstrap()
	}

	if flags.updateContracts {
		fHan := files.GetFileHandler(&logger)
		setupData = fHan.ReadSetupData(config)
		conHan := contracts.SetupContractHandler(config, &setupData, &logger)
		conHan.UpdateContracts()
	}

}

func handleFlags() CLIFlags {
	setupTestnetPtr := flag.Bool(
		"t",
		false,
		"Setup a testnet",
	)

	updateContractsPtr := flag.Bool(
		"u",
		false,
		"Update contracts",
	)

	verboseLoggingPtr := flag.Bool(
		"v",
		false,
		"Verbose logging mode, no logs will be output without this flag",
	)

	headlessPtr := flag.Bool(
		"hl",
		false,
		"Run in headless mode, won't ask questions, be careful!",
	)

	flag.Parse()

	flags := CLIFlags{
		setupTestnet:    *setupTestnetPtr,
		updateContracts: *updateContractsPtr,
		verboseLogging:  *verboseLoggingPtr,
		headlessMode:    *headlessPtr,
	}

	if !flags.setupTestnet && !flags.updateContracts {
		help := flag.Usage
		fmt.Printf("Please include a CLI flag\n\n")
		help()
		os.Exit(0)
	}

	return flags
}
