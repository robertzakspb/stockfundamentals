CREATE TABLE company (
    id Utf8 NOT NULL,
    company_name Utf8,
    is_public Bool,
    isin Utf8,
    security_type Utf8,
    country_iso2 Utf8,
    ordinary_share_ticker Utf8,
    ordinary_share_count Uint64,
    preferred_share_ticker Utf8,
    preferred_share_count Uint64,
    PRIMARY KEY (id)
);

CREATE TABLE dividend_payment(
	id Utf8,
	company_id Utf8,
	actual_DPS Double,
	expected_DPS Double,
	currency Utf8,
    ex_div_date Date,
	payout_date Date,
    payment_period Utf8,
	management_comment Utf8,
	PRIMARY KEY(id)
);

CREATE TABLE corporate_financials(
    id Utf8,
	company_id Utf8,
    financial_metric Utf8,
    reporting_period Utf8,
    metric_value Utf8,
    metric_currency Utf8,
)

-- CREATE TABLE dividend_policy(
    
-- )