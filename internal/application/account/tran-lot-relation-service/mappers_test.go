package tranlotrelationservice

import (
	"testing"
	"time"

	tranlotrelation "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/tran-lot-relation"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_mapTranLotRelationsToDbModel(t *testing.T) {
	stockLotId, bondLotIt := uuid.New(), uuid.New()
	date := time.Now()
	quantity := 25.4
	relation := tranlotrelation.TransactionLotRelation{
		StockLotId: stockLotId,
		BondLotId:  bondLotIt,
		Date:       date,
		Quantity:   quantity,
	}

	mappedDbModels := mapTranLotRelationsToDbModel([]tranlotrelation.TransactionLotRelation{relation})

	test.AssertEqual(t, 1, len(mappedDbModels))
	test.AssertEqual(t, stockLotId, mappedDbModels[0].StockLotId)
	test.AssertEqual(t, bondLotIt, mappedDbModels[0].BondLotId)
	test.AssertEqual(t, date, mappedDbModels[0].Date)
	test.AssertEqual(t, quantity, mappedDbModels[0].Quantity)
}
