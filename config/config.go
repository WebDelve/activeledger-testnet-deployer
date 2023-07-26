package config

import (
	"encoding/json"

	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/helper"
)

type Config struct {
	DefaultIdentity   string `json:"identity"`
	DefaultNamespace  string `json:"namespace"`
	ContractFolder    string `json:"contractDir"`
	ContractManifest  string `json:"contractManifest"`
	SetupDataSaveFile string `json:"setupDataSaveFile"`
	TestnetFolder     string `json:"testnetFolder"`
}

func LoadConfig() *Config {
	var c Config

	data := files.ReadFile("./config.json")

	if err := json.Unmarshal(data, &c); err != nil {
		helper.HandleError(err, "Error unmarshalling content of config file")
	}

	return &c
}
