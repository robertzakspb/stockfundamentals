package lot

import (
	"database/sql"
)

func GetLotsFromDB(db *sql.DB) ([]Lot, error) {
	rows, err := db.Query("SELECT company.name, company.ordinary_share_ticker, lot.quantity, lot.cost_basis, broker_account.name, lot.opening_date FROM lot JOIN position on lot.position_id = position.id JOIN company on position.security_id = company.id JOIN broker_account on position.account_id = broker_account.id")
	if err != nil {
		return []Lot{}, err
	}
	defer rows.Close()

	lots := []Lot{}
	for rows.Next() {
		var lot Lot
		if err := rows.Scan(&lot.CompanyName, &lot.Ticker, &lot.Quantity, &lot.OpeningPrice, &lot.BrokerName); err != nil {
			continue
		}
		lots = append(lots, lot)
	}

	if err = rows.Err(); err != nil {
		return lots, err
	}

	return lots, nil
}

// func AddLot(lot Lot, db *sql.DB) error {
// _, err := company.GetCompanyFromDB(lot.Ticker, db)

// return err
// }
