package contracts

import (
	"crypto/sha256"
	"encoding/json"

	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/helper"
	"dynamicledger.com/testnet-deployer/structs"
)

func (ch *ContractHandler) getManifest() {
	var man structs.ContractManifest

	data := files.ReadFile(ch.Config.ContractManifest)
	if err := json.Unmarshal(data, &man); err != nil {
		helper.HandleError(err, "Error parsing contract manifest file")
	}

	ch.Manifest = man
}

func (ch *ContractHandler) updateContractHashes() {

	newMetadata := []structs.ContractMetadata{}

	for _, c := range ch.Contracts {
		con := matchContractManifest(c.Name, ch.Manifest.Contracts)
		con.Hash = getContractHash(c)
		newMetadata = append(newMetadata, con)
	}

	ch.Manifest.Contracts = newMetadata
	ch.storeManifest()

}

func matchContractManifest(name string, manifest []structs.ContractMetadata) structs.ContractMetadata {
	for _, c := range manifest {
		if c.Name == name {
			return c
		}
	}

	return structs.ContractMetadata{}
}

func (ch *ContractHandler) storeManifest() {
	bMan, err := json.Marshal(ch.Manifest)
	if err != nil {
		helper.HandleError(err, "Error marshalling manifest data")
	}

	files.WriteFile(ch.Config.ContractManifest, bMan)
}

func getContractHash(contract structs.Contract) string {
	hasher := sha256.New()
	hasher.Write([]byte(contract.Data))

	hash := hasher.Sum(nil)

	return string(hash)
}
