-- CREATE TABLE corporate_financials(
    id Utf8,
	company_id Utf8,
    financial_metric Utf8,
    reporting_period Utf8,
    metric_value Utf8,
    metric_currency Utf8,
)
-- Jedinstvo iz Sevojna
INSERT INTO corporate_financials(id, company_id, financial_metric, reporting_period, metric_value, metric_currency) VALUES
('a0de3d29-7ed2-4696-ba41-c8f4b80ade07', 'dd194350-4c61-4643-8c74-1120ceca8fae', 'revenue', 'fy2009', 4110977000, 'RSD'),