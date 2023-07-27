package structs

type SetupStore struct {
	Identity     string          `json:"identity"`
	Namespace    string          `json:"namespace"`
	KeyData      KeyStore        `json:"keyData"`
	ContractData []ContractStore `json:"contractData"`
}

type KeyStore struct {
	PublicPem   string `json:"publicPem"`
	PublicHash  string `json:"publicHash"`
	PrivatePem  string `json:"privatePem"`
	PrivateHash string `json:"privateHash"`
}
