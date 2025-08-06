package securityinfo

type SecurityType string

const (
	OrdinaryShare     SecurityType = "commonStock"
	PreferredShare    SecurityType = "preferredStock"
	DepositoryReceipt SecurityType = "depositoryReceipt"
)

type MoexSecurity struct {
	Ticker       string `sql:"share_ticker"`
	ShortName    string `sql:"company_name"`
	FullName     string
	EnglishName  string
	LotSize      int64
	FaceValue    float64
	SharesIssued int64        `sql:"share_count"`
	ISIN         string       `sql:"isin"`
	SecurityType SecurityType `sql:"security_type"`
}
