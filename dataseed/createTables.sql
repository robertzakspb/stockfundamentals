CREATE TABLE `stockfundamentals/stocks/stock` (
    id Utf8,
    company_name Utf8,
    is_public Bool,
    isin Utf8 NOT NULL,
    security_type Utf8 NOT NULL,
    country_iso2 Utf8,
    ticker Utf8,
    share_count Uint64,
    sector Utf8,
    PRIMARY KEY (id, isin, ticker)
);

CREATE TABLE `stockfundamentals/stocks/dividend_payment`(
	id Utf8,
	company_id Utf8,
	actual_DPS Double,
	expected_DPS Double,
	currency Utf8,
    announcement_date Date,
    record_date Date,
	payout_date Date,
    payment_period Utf8,
	management_comment Utf8,
	PRIMARY KEY(id)
);

CREATE TABLE `stockfundamentals/stocks/corporate_financials`(
    id Utf8,
	company_id Utf8,
    financial_metric Utf8,
    reporting_period Utf8,
    metric_value Utf8,
    metric_currency Utf8,
    PRIMARY KEY(id)
)