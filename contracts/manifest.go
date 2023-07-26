package contracts

import (
	"crypto/sha256"
	"encoding/json"

	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/helper"
)

func (ch *ContractHandler) getManifest() {
	var man contractManifest

	data := files.ReadFile(ch.config.ContractManifest)
	if err := json.Unmarshal(data, &man); err != nil {
		helper.HandleError(err, "Error parsing contract manifest file")
	}

	ch.manifest = man
}

func (ch *ContractHandler) updateContractHashes() {

	newMetadata := []contractMetadata{}

	for _, c := range ch.contracts {
		con := matchContractManifest(c.name, ch.manifest.Contracts)
		con.Hash = getContractHash(c)
		newMetadata = append(newMetadata, con)
	}

	ch.manifest.Contracts = newMetadata
	ch.storeManifest()

}

func matchContractManifest(name string, manifest []contractMetadata) contractMetadata {
	for _, c := range manifest {
		if c.Name == name {
			return c
		}
	}

	return contractMetadata{}
}

func (ch *ContractHandler) storeManifest() {
	bMan, err := json.Marshal(ch.manifest)
	if err != nil {
		helper.HandleError(err, "Error marshalling manifest data")
	}

	files.WriteFile(ch.config.ContractManifest, bMan)
}

func getContractHash(contract Contract) string {
	hasher := sha256.New()
	hasher.Write([]byte(contract.data))

	hash := hasher.Sum(nil)

	return string(hash)
}
