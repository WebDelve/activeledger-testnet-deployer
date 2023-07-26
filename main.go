package main

import (
	"flag"
	"fmt"
	"os"

	"dynamicledger.com/testnet-deployer/bootstrap"
	"dynamicledger.com/testnet-deployer/config"
	"dynamicledger.com/testnet-deployer/contracts"
	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type CLIFlags struct {
	setupTestnet    bool
	updateContracts bool
}

func main() {
	flags := handleFlags()

	config := config.LoadConfig()

	setupData := structs.SetupData{}

	setupData.Conn = alsdk.Connection{
		Protocol: alsdk.HTTP,
		Url:      "localhost",
		Port:     "5260",
	}

	setupData.Namespace = config.DefaultNamespace

	if flags.setupTestnet {

		bs := bootstrap.GetBootstrapper(config, &setupData)
		bs.Bootstrap()
	}

	if flags.updateContracts {
		setupData = files.ReadSetupData(config)
		conHan := contracts.SetupContractHandler(config, &setupData)
		conHan.UpdateContracts()
	}

}

func handleFlags() CLIFlags {
	setupTestnetPtr := flag.Bool("t", false, "Setup a testnet")
	updateContractsPtr := flag.Bool("u", false, "Update contracts")

	flag.Parse()

	flags := CLIFlags{
		setupTestnet:    *setupTestnetPtr,
		updateContracts: *updateContractsPtr,
	}

	if !flags.setupTestnet && !flags.updateContracts {
		help := flag.Usage
		fmt.Printf("Please include a CLI flag\n\n")
		help()
		os.Exit(0)
	}

	return flags
}
