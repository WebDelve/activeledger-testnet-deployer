package structs

type Config struct {
	DefaultIdentity   string `json:"identity"`
	DefaultNamespace  string `json:"namespace"`
	ContractFolder    string `json:"contractDir"`
	ContractManifest  string `json:"contractManifest"`
	SetupDataSaveFile string `json:"setupDataSaveFile"`
	TestnetFolder     string `json:"testnetFolder"`
}
