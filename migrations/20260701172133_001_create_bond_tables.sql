-- +goose Up
CREATE TABLE `bonds/bond`(
    id Uuid,
    figi Text,
    isin Text,
    lot Int64,
    currency Text,
    name Text,
    country_of_risk Text,
    real_exchange Text,
    coupon_count_per_year Int64,
    maturity_date Date,
    nominal_value Double,
    nominal_currency Text,
    initial_nominal_value Double,
    initial_nominal_currency Text,
    registration_date Date,
    placement_date Date,
    placement_price Double,
    placement_currency Text,
    accumulated_coupon_income Double,
    issue_size Int64,
    issue_size_plan Int64,
    has_floating_coupon Bool,
    is_perpetual Bool,
    has_amortization Bool,
    is_available_for_iis Bool,
    is_for_qualified_investors Bool,
    is_subordinated Bool,
    risk_level Text,
    bond_type Text,
    call_option_exercise_date Date,
    PRIMARY KEY (figi, isin)
);

CREATE TABLE `bonds/coupon`(
    id Uuid,
    figi Text,
    coupon_date Date,
    record_date Date,
    coupon_number Int64,
    per_bond_amount Double,
    coupon_type Text,
    coupon_start_date Date,
    coupon_end_date Date,
    coupon_period Int64,
    PRIMARY KEY(figi, coupon_date)
);

-- +goose Down
DROP TABLE `bonds/bond`;
DROP TABLE `bonds/coupon`;


