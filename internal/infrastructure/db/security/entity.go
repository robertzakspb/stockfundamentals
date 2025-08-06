package security

import "github.com/google/uuid"



type StockDbModel struct {
	Id           uuid.UUID `sql:"id"`
	Isin         string    `sql:"isin"`
	Figi         string    `sql:"figi"`
	CompanyName  string    `sql:"company_name"`
	IsPublic     bool      `sql:"is_public"`
	SecurityType string    `sql:"security_type"`
	Country      string    `sql:"country_iso2"`
	Ticker       string    `sql:"ticker"`
	IssueSize    int64     `sql:"issue_size"`
	Sector       string    `sql:"sector"`
	MIC          string    `sql:"MIC"`
}

