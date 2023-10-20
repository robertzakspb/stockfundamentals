package securityinfo

type SecurityType string

const (
	OrdinaryShare     SecurityType = "ordinaryStock"
	PreferredShare    SecurityType = "preferredStock"
	DepositoryReceipt SecurityType = "depositoryReceipt"
)

type MoexSecurity struct {
	Ticker       string
	ShortName    string
	FullName     string
	EnglishName  string
	LotSize      float64
	FaceValue    float64
	SharesIssued float64
	ISIN         string
	SecurityType SecurityType
}
