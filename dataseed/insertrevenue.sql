-- CREATE TABLE corporate_financials(
    id Utf8,
	company_id Utf8,
    financial_metric Utf8,
    reporting_period Utf8,
    metric_value Utf8,
    metric_currency Utf8,
)
-- Jedinstvo iz Sevojna
INSERT INTO corporate_financials(id, security_id, financial_metric, reporting_period, metric_value, metric_currency) VALUES
('a0de3d29-7ed2-4696-ba41-c8f4b80ade07', 'dd194350-4c61-4643-8c74-1120ceca8fae', 'revenue', 'fy2009', 4110977000, 'RSD'),

--Metalac
('7d7890da-e7f7-41ce-b4ec-0e2fb05f388b', '8f0161f0-083d-431a-9fff-89deb073ce0f', 'revenue', 'fy2014', 7824790000, 'RSD'),
('1127aaff-ff53-41d3-b24a-f0c7bb36d618', '8f0161f0-083d-431a-9fff-89deb073ce0f', 'revenue', 'fy2015', 8284390000, 'RSD'),
('a1dc421a-63b3-4568-b439-e1f28b501206', '8f0161f0-083d-431a-9fff-89deb073ce0f', 'revenue', 'fy2016', 8700000000, 'RSD'),
('7b06ef91-7ec2-4bcd-a45b-d8498042633d', '8f0161f0-083d-431a-9fff-89deb073ce0f', 'revenue', 'fy2017', 9341870000, 'RSD'),
('f4715455-e936-4768-a8e4-de28ec2c6dad', '8f0161f0-083d-431a-9fff-89deb073ce0f', 'revenue', 'fy2018', 9994852000, 'RSD'),