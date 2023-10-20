package company

type Company struct {
	ID                   string
	Name                 string
	SecurityType         string
	Country              string
	IsPublic             bool
	OrdinaryShareCount   int
	OrdinaryShareTicker  string
	PreferredShareCount  int
	PreferredShareTicker string
}
