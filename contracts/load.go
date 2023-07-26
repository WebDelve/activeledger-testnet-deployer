package contracts

import (
	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/structs"

	b64 "encoding/base64"
)

func (ch *ContractHandler) loadContracts() {
	contracts := []structs.Contract{}

	ch.getManifest()

	for k, c := range ch.Manifest.Contracts {
		if c.Exclude {
			continue
		}

		con := ch.readContract(c)

		// Only hash if one isn't set, we only want to update hashes after uploading
		// updated contracts
		var blank string
		if c.Hash == blank {
			hash := getContractHash(con)
			ch.Manifest.Contracts[k].Hash = hash
		}

		contracts = append(contracts, con)
	}

	// Write the manifest file in case there are new hashes
	ch.storeManifest()

	ch.Contracts = contracts
}

func (ch *ContractHandler) readContract(contractMeta structs.ContractMetadata) structs.Contract {
	c := structs.Contract{}
	c.Name = contractMeta.Name
	c.Version = contractMeta.Version

	data := files.ReadFile(contractMeta.Path)

	c.Data = b64.StdEncoding.EncodeToString(data)

	return c
}
