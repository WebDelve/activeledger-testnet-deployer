package config

import (
	"encoding/json"

	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/helper"
)

func LoadConfig() *helper.Config {
	var c helper.Config

	data := files.ReadFile("./config.json")

	if err := json.Unmarshal(data, &c); err != nil {
		helper.HandleError(err, "Error unmarshalling content of config file")
	}

	return &c
}
