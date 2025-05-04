package security

type SecurityType string

const (
	Unspecified       SecurityType = "unspecified"
	OrdinaryShare     SecurityType = "ordinaryStock"
	PreferredShare    SecurityType = "preferredStock"
	DepositoryReceipt SecurityType = "depositoryReceipt"
)

var(
	securityTypeMap = map[string]SecurityType {
		"unspecified": Unspecified,
		"ordinaryStock": OrdinaryShare,
		"preferredStock": PreferredShare,
		"depositoryReceipt": DepositoryReceipt,
	}
)

type Security interface {
	GetId() string
	GetCompanyName() string
	GetIsPublic() bool
	GetIsin() string
	GetSecurityType() SecurityType
	GetCountry() string
	GetTicker() string
	GetIssueSize() int
	GetSector() string
	GetFigi() string
}

type Stock struct {
	Id           string
	Isin         string
	Figi         string
	CompanyName  string
	IsPublic     bool
	SecurityType SecurityType 
	Country      string       
	Ticker       string
	IssueSize    int
	Sector       string
}

type StockDbModel struct {
	Id           string `sql:"id"`
	Isin         string `sql:"isin"`
	Figi         string `sql:"figi"`
	CompanyName  string `sql:"company_name"`
	IsPublic     bool   `sql:"is_public"`
	SecurityType string `sql:"security_type"`
	Country      string `sql:"country_iso2"`
	Ticker       string `sql:"ticker"`
	IssueSize    int `sql:"issue_size"`
	Sector       string `sql:"sector"`
}

//Implementing the Security interface
func (stock Stock) GetId() string {
	return stock.Id
}

func (stock Stock) GetCompanyName() string {
	return stock.CompanyName
}

func (stock Stock) GetIsPublic() bool {
	return stock.IsPublic
}

func (stock Stock) GetIsin() string {
	return stock.Isin
}

func (stock Stock) GetFigi() string {
	return stock.Figi
}

func (stock Stock) GetSecurityType() SecurityType {
	return stock.SecurityType
}

func (stock Stock) GetCountry() string {
	return stock.Country
}

func (stock Stock) GetTicker() string {
	return stock.Ticker
}

func (stock Stock) GetIssueSize() int {
	return stock.IssueSize
}

func (stock Stock) GetSector() string {
	return stock.Sector
}
