package structs

type ContractStore struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Hash string `json:"hash"`
}

type Contract struct {
	Name    string
	Id      string
	Data    string
	Version string
}

type ContractManifest struct {
	Contracts []ContractMetadata `json:"contracts"`
}

type ContractMetadata struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Path      string `json:"path"`
	Version   string `json:"version"`
	Hash      string `json:"hash"`
	Exclude   bool   `json:"exclude"`
	Onboarded bool   `json:"onboarded"`
}
