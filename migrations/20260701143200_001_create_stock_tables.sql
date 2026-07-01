-- +goose Up
CREATE TABLE `stockfundamentals/stocks/dividend_forecast`(
    id Uuid,
    figi Text,
    expected_DPS Double,
    currency Text,
    payment_period Text,
    forecast_author Text,
    comment Text,
    payout_date Date,
    PRIMARY KEY (figi, payout_date)
);

CREATE TABLE `stockfundamentals/stocks/stock`(
    figi Text,
    company_name Text,
    is_public Bool,
    isin Text,
    security_type Text,
    country_iso2 Text,
    MIC Text,
    ticker Text,
    issue_size Text,
    sector Text,
    PRIMARY KEY(figi)
);

CREATE TABLE `stockfundamentals/stocks/dividend_payment`(
    id Uuid,
    stock_id Text,
    actual_DPS Int64,
    expected_DPS Int64,
    currency Text,
    announcement_date Date,
    record_date Date,
    payout_date Date,
    payment_period Text,
    type Text,
    regularity Text,
    management_comment Text,
    PRIMARY KEY(stock_id, record_date, actual_DPS)
);

CREATE TABLE `stockfundamentals/stocks/financial_metric`(
    id Uuid,
    figi Text,
    metric Text,
    reporting_period Text,
    year Uint16,
    metric_value Int64,
    metric_currency Text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE `stockfundamentals/stocks/dividend_forecast`;
DROP TABLE `stockfundamentals/stocks/stock`;
DROP TABLE `stockfundamentals/stocks/dividend_payment`;
DROP TABLE `stockfundamentals/stocks/financial_metric`;
