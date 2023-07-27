package config

import (
	"encoding/json"

	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/logging"
	"dynamicledger.com/testnet-deployer/structs"
)

func LoadConfig(logger *logging.Logger) *structs.Config {
	var c structs.Config

	fHan := files.GetFileHandler(logger)
	data := fHan.ReadFile("./config.json")

	if err := json.Unmarshal(data, &c); err != nil {
		logger.Fatal(err, "Error unmarshalling content of config file")
	}

	return &c
}
