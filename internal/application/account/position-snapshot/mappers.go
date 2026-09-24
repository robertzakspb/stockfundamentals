package positionsnapshot

import positionsnapshotdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/position-snapshot"

func mapPositionSnapshotsToDbModels(snapshots []StockPositionSnapshot) []positionsnapshotdb.StockPositionSnapshotDbModel {
	dbModels := make([]positionsnapshotdb.StockPositionSnapshotDbModel, len(snapshots))

	for i := range snapshots {
		dbModel := positionsnapshotdb.StockPositionSnapshotDbModel{
			Figi:      snapshots[i].Figi,
			Isin:      snapshots[i].Isin,
			AccountId: snapshots[i].AccountId,
			Date:      snapshots[i].Date,
			Quantity:  snapshots[i].Quantity,
		}
		dbModels[i] = dbModel
	}
	return dbModels
}
