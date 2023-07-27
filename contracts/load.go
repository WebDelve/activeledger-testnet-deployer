package contracts

import (
	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/structs"

	b64 "encoding/base64"
)

func (ch *ContractHandler) loadContracts() {
	contracts := []structs.Contract{}

	for _, c := range ch.Manifest.Contracts {
		if c.Exclude {
			continue
		}

		con := ch.readContract(c)

		// Checks if there are missing hashes and adds them
		ch.setHashes(true)

		contracts = append(contracts, con)
	}

	ch.Contracts = contracts
}

func (ch *ContractHandler) readContract(contractMeta structs.ContractMetadata) structs.Contract {
	c := structs.Contract{}
	c.Name = contractMeta.Name
	c.Version = contractMeta.Version
	c.Id = contractMeta.ID

	data := files.ReadFile(contractMeta.Path)

	c.Data = b64.StdEncoding.EncodeToString(data)

	return c
}
