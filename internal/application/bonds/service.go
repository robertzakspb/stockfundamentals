package bondservice

import (
	"context"
	"errors"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/compoundinvest/invest-core/quote/bondquote"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
	pb "opensource.tbank.ru/invest/invest-go/proto"
)

func ImportAllBondsAndCoupons() error {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	client, err := tinkoff.NewClient(ctx, config, nil)
	if err != nil {
		logger.Log("Failed to initialize the Tinkoff API client: ", logger.ALERT)
		return err
	}

	bondService := client.NewInstrumentsServiceClient()
	response, err := bondService.Bonds(pb.InstrumentStatus_INSTRUMENT_STATUS_ALL)
	if response == nil {
		logger.Log("Unexpectedly received a nil response from Tinkoff API", logger.ALERT)
	}

	dbBonds := []bondsdb.BondDbModel{}
	for _, tinkoffBond := range response.Instruments {
		if tinkoffBond.MaturityDate.AsTime().Before(time.Now()) {
			//No need to import historical bonds that have matured
			continue
		}
		bond := mapTinkoffBondToBond(tinkoffBond)
		validationErr := bond.Validate()
		if validationErr != nil {
			logger.Log(validationErr.Error(), logger.WARNING)
		}
		dbBond := mapBondToDbBond(bond)
		dbBonds = append(dbBonds, dbBond)
	}

	err = bondsdb.SaveBonds(dbBonds)
	if err != nil {
		return err
	}

	go importAllCoupons()

	return nil
}

func GetAllBonds() ([]bonds.Bond, error) {
	return GetFilteredBonds([]ydbfilter.YdbFilter{})
}

func GetFilteredBonds(filters []ydbfilter.YdbFilter) ([]bonds.Bond, error) {
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "maturity_date",
		Condition:      ydbfilter.GreaterThan,
		ConditionValue: shared.ConvertToYdbDate(time.Now()),
	}
	filters = append(filters, filter)
	bondList, err := bondsdb.GetAllBonds(filters)
	if err != nil {
		return []bonds.Bond{}, err
	}

	if len(bondList) == 0 {
		return []bonds.Bond{}, errors.New("Found zero bonds in the DB")
	}

	mappedBonds := []bonds.Bond{}
	for _, dbBond := range bondList {
		mappedBond := mapDbBondToBond(dbBond)
		mappedBonds = append(mappedBonds, mappedBond)
	}

	return mappedBonds, nil
}

func GetBondByFigi(figi string) (bonds.Bond, error) {
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "figi",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue(figi),
	}
	bondList, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return bonds.Bond{}, err
	}
	if len(bondList) == 0 {
		return bonds.Bond{}, errors.New("Found zero bonds with the specificed figi")
	}

	mappedBond := mapDbBondToBond(bondList[0])

	return mappedBond, nil
}

func GetBondsByFigi(figis []string) ([]bonds.Bond, error) {
	ydbFigis := []types.Value{}
	for _, figi := range figis {
		ydbFigis = append(ydbFigis, types.TextValue(figi))
	}

	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "figi",
		Condition:      ydbfilter.Contains,
		ConditionValue: types.ListValue(ydbFigis...),
	}
	bondList, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return []bonds.Bond{}, err
	}
	if len(bondList) == 0 {
		return []bonds.Bond{}, errors.New("Found zero bonds with the specificed figis")
	}

	mappedBonds := []bonds.Bond{}
	for _, dbBond := range bondList {
		mappedBond := mapDbBondToBond(dbBond)
		mappedBonds = append(mappedBonds, mappedBond)
	}

	return mappedBonds, nil
}

func GetBondByIsin(isin string) (bonds.Bond, error) {
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "isin",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue(isin),
	}
	bondList, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return bonds.Bond{}, err
	}
	if len(bondList) == 0 {
		return bonds.Bond{}, errors.New("Found zero bonds with the specificed ISIN")
	}

	mappedBond := mapDbBondToBond(bondList[0])

	return mappedBond, nil
}

func GetBondsByIsin(isins []string) ([]bonds.Bond, error) {
	ydbIsins := []types.Value{}
	for _, ydbIsin := range isins {
		ydbIsins = append(ydbIsins, types.TextValue(ydbIsin))
	}
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "isin",
		Condition:      ydbfilter.Contains,
		ConditionValue: types.ListValue(ydbIsins...),
	}
	bondList, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return []bonds.Bond{}, err
	}
	if len(bondList) == 0 {
		return []bonds.Bond{}, errors.New("Found zero bonds with the specificed ISINs")
	}

	mappedBonds := []bonds.Bond{}
	for _, bond := range bondList {
		mappedBond := mapDbBondToBond(bond)
		mappedBonds = append(mappedBonds, mappedBond)
	}

	return mappedBonds, nil
}

func GetCouponsByFigi(figi string) ([]bonds.Coupon, error) {
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "figi",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue(figi),
	}

	coupons, err := bondsdb.GetBondCoupons([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return []bonds.Coupon{}, err
	}

	mappedCoupons := []bonds.Coupon{}
	for _, coupon := range coupons {
		mappedCoupon := mapCouponDbModelToDomain(coupon)
		mappedCoupons = append(mappedCoupons, mappedCoupon)
	}
	return mappedCoupons, nil
}

func GetCouponsByFigis(figis []string) ([]bonds.Coupon, error) {
	ydbFigis := []types.Value{}
	for _, figi := range figis {
		ydbFigis = append(ydbFigis, types.TextValue(figi))
	}
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "figi",
		Condition:      ydbfilter.Contains,
		ConditionValue: types.ListValue(ydbFigis...),
	}

	coupons, err := bondsdb.GetBondCoupons([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return []bonds.Coupon{}, err
	}

	mappedCoupons := []bonds.Coupon{}
	for _, coupon := range coupons {
		mappedCoupon := mapCouponDbModelToDomain(coupon)
		mappedCoupons = append(mappedCoupons, mappedCoupon)
	}
	return mappedCoupons, nil
}

func PopulateBondCoupons(bondList []bonds.Bond) []bonds.Bond {
	figis := []string{}
	for _, bond := range bondList {
		figis = append(figis, bond.Figi)
	}
	coupons, err := GetCouponsByFigis(figis)
	if err != nil {
		logger.Log("Failed to fetch coupons for the provided bonds", logger.ERROR)
		return bondList
	}
	bondsWithCoupons := MatchCouponsWithBonds(coupons, bondList)
	return bondsWithCoupons
}

func MatchCouponsWithBonds(coupons []bonds.Coupon, bonds []bonds.Bond) []bonds.Bond {
	for _, coupon := range coupons {
		for i, b := range bonds {
			if coupon.Figi == b.Figi {
				bonds[i].Coupons = append(b.Coupons, coupon)
			}
		}
	}
	return bonds
}

func CalculateYtmForBonds(bondList []bonds.Bond, quotes []bondquote.TinkoffBondQuote) []bonds.Bond {
	for _, quote := range quotes {
		for i, b := range bondList {
			if quote.Figi() == b.Figi {
				ytm, err := b.CalcYieldToMaturity(b.Coupons, quote.QuoteAsPercentage())
				if err != nil {
					logger.Log(err.Error(), logger.ERROR)
					continue
				}
				bondList[i].YieldToMaturity = ytm

				if b.CallOptionExerciseDate.IsZero() == false {
					yieldToCallOption, err := b.CalcYieldToCallOption(b.Coupons, quote.QuoteAsPercentage())
					if err != nil {
						logger.Log(err.Error(), logger.ERROR)
						continue
					}
					bondList[i].YieldToCallOption = yieldToCallOption
				}
			}
		}
	}
	return bondList
}

func UpdateAllBondsAci() error {
	bondList, err := GetAllBonds()
	bondList = PopulateBondCoupons(bondList)

	if err != nil {
		return err
	}

	for i, bond := range bondList {
		aci, err := bonds.AccruedInterest(bond, time.Now())
		if err != nil {
			logger.Log(err.Error(), logger.WARNING)
			continue
		}
		bondList[i].AccruedInterest = aci
	}

	dbBonds := []bondsdb.BondDbModel{}
	for _, bond := range bondList {
		dbBond := mapBondToDbBond(bond)
		dbBonds = append(dbBonds, dbBond)
	}

	err = bondsdb.SaveBonds(dbBonds)
	if err != nil {
		return err
	}

	return nil
}

func AllCurrencyPairsInBondList(bondList []bonds.Bond) map[string]string {
	pairs := []string{}

	for _, bond := range bondList {
		if bond.Currency != bond.NominalCurrency {
			foundPair := false
			for _, pair := range pairs {
				if pair == bond.NominalCurrency+"/"+bond.Currency {
					foundPair = true
				}
			}
			if !foundPair {
				pairs = append(pairs, bond.NominalCurrency+"/"+bond.Currency)
			}
		}
	}

	pairMap := map[string]string {}
	for _, pair := range pairs {
		split := strings.Split(pair, "/")
		pairMap[split[0]] = split[1]
	}

	return pairMap
}
