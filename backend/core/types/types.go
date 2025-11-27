package types

type ApiEnpointScope struct {
	Insert   bool `json:"insert"`
	Update   bool `json:"update"`
	Delete   bool `json:"delete"`
	Retrieve bool `json:"retrieve"`
}

type ApiDiscovery struct {
	ViewPath  string   `json:"viewPath"`
	EndPoints []string `json:"endpoints"`
}
