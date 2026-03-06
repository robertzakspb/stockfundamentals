package bondservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	pb "opensource.tbank.ru/invest/invest-go/proto"
)

func mapTinkoffBondToBond(tinkoffBond *pb.Bond) bonds.Bond {
	bond := bonds.Bond{
		Figi:     tinkoffBond.Figi,
		Isin:     tinkoffBond.Isin,
		Lot:      int(tinkoffBond.Lot),
		Currency: tinkoffBond.Currency,
	}

	return bond
}
