package bondsdb

import (
	"context"
	"errors"
	"path"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func SaveBonds(bonds []BondDbModel) error {
	dbConnection, err := db.MakeYdbDriver()
	if err != nil {
		return err
	}
	defer dbConnection.Close(context.TODO())

	ydbBonds := []types.Value{}
	for _, bond := range bonds {
		ydbBond := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(bond.Id)),
			types.StructFieldValue("figi", types.TextValue(bond.Figi)),
			types.StructFieldValue("isin", types.TextValue(bond.Isin)),
			types.StructFieldValue("lot", types.Int64Value(int64(bond.Lot))),
			types.StructFieldValue("currency", types.TextValue(bond.Currency)),
			types.StructFieldValue("name", types.TextValue(bond.Name)),
			types.StructFieldValue("country_of_risk", types.TextValue(bond.CountryOfRisk)),
			types.StructFieldValue("real_exchange", types.TextValue(bond.RealExchange)),
			types.StructFieldValue("coupon_count_per_year", types.Int64Value(int64(bond.CouponCountPerYear))),
			types.StructFieldValue("maturity_date", bond.CorrectMaturityDate()),
			types.StructFieldValue("nominal_value", types.DoubleValue(bond.NominalValue)),
			types.StructFieldValue("nominal_currency", types.TextValue(bond.NominalCurrency)),
			types.StructFieldValue("initial_nominal_value", types.DoubleValue(bond.InitialNominalValue)),
			types.StructFieldValue("initial_nominal_currency", types.TextValue(bond.InitialNominalCurrency)),
			types.StructFieldValue("registration_date", ydbhelper.ConvertToOptionalYDBdate(bond.RegistrationDate)),
			types.StructFieldValue("placement_date", ydbhelper.ConvertToOptionalYDBdate(bond.PlacementDate)),
			types.StructFieldValue("placement_price", types.DoubleValue(bond.PlacementPrice)),
			types.StructFieldValue("placement_currency", types.TextValue(bond.PlacementCurrency)),
			types.StructFieldValue("accumulated_coupon_income", types.DoubleValue(bond.AccruedInterest)),
			types.StructFieldValue("issue_size", types.Int64Value(int64(bond.IssueSize))),
			types.StructFieldValue("issue_size_plan", types.Int64Value(int64(bond.IssueSizePlan))),
			types.StructFieldValue("has_floating_coupon", types.BoolValue(bond.HasFloatingCoupon)),
			types.StructFieldValue("is_perpetual", types.BoolValue(bond.IsPerpetual)),
			types.StructFieldValue("has_amortization", types.BoolValue(bond.HasAmortization)),
			types.StructFieldValue("is_available_for_iis", types.BoolValue(bond.IsAvailableForIis)),
			types.StructFieldValue("is_for_qualified_investors", types.BoolValue(bond.IsForQualifiedInvestors)),
			types.StructFieldValue("is_subordinated", types.BoolValue(bond.IsSubordinated)),
			types.StructFieldValue("risk_level", types.TextValue(bond.RiskLevel)),
			types.StructFieldValue("bond_type", types.TextValue(bond.BondType)),
			types.StructFieldValue("call_option_exercise_date", ydbhelper.ConvertToOptionalYDBdate(bond.CallOptionExerciseDate)),
		)
		ydbBonds = append(ydbBonds, ydbBond)
	}

	tableName := path.Join(dbConnection.Name(), db.BOND_DIRECTORY_PREFIX, db.BOND_TABLE_NAME)
	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbBonds...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save bonds to the database")
	}

	return nil
}

func SaveCoupons(coupons *[]CouponDbModel) error {
	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		return err
	}
	defer dbConnection.Close(context.TODO())

	ydbCoupons := []types.Value{}
	for _, c := range *coupons {
		ydbCoupon := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(c.Id)),
			types.StructFieldValue("figi", types.TextValue(c.Figi)),
			types.StructFieldValue("coupon_date", ydbhelper.ConvertToYdbDate(c.CouponDate)),
			types.StructFieldValue("coupon_number", types.Int64Value(int64(c.CouponNumber))),
			types.StructFieldValue("record_date", ydbhelper.ConvertToYdbDate(c.RecordDate)),
			types.StructFieldValue("per_bond_amount", types.DoubleValue(c.PerBondAmount)),
			types.StructFieldValue("coupon_type", types.TextValue(c.CouponType)),
			types.StructFieldValue("coupon_start_date", ydbhelper.ConvertToYdbDate(c.CouponStartDate)),
			types.StructFieldValue("coupon_end_date", ydbhelper.ConvertToYdbDate(c.CouponEndDate)),
			types.StructFieldValue("coupon_period", types.Int64Value(int64(c.CouponPeriod))),
		)
		ydbCoupons = append(ydbCoupons, ydbCoupon)
	}
	tableName := path.Join(dbConnection.Name(), db.BOND_DIRECTORY_PREFIX, db.COUPON_TABLE_NAME)
	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbCoupons...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save coupons to the database")
	}

	return nil
}
