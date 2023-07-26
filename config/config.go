package config

import (
	"encoding/json"

	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/helper"
	"dynamicledger.com/testnet-deployer/structs"
)

func LoadConfig() *structs.Config {
	var c structs.Config

	data := files.ReadFile("./config.json")

	if err := json.Unmarshal(data, &c); err != nil {
		helper.HandleError(err, "Error unmarshalling content of config file")
	}

	return &c
}
