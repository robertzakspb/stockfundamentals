-- +goose Up
UPSERT INTO `stockfundamentals/stocks/stock` (figi, company_name, is_public, isin, security_type, country_iso2, MIC, ticker, issue_size, sector) VALUES 
('BBG000BMX476', 'Dunav Osiguranje', true, 'RSDNOSE74915', 'commonStock', 'RS', 'XBEL', 'DNOS', 15189202, 'Insurance'),
('BBG000BS7XH7', 'Jedinstvo iz Sevojna', true, 'RSJESVE87017', 'commonStock', 'RS', 'XBEL', 'JESV', 232703, 'Construction'),
('BBG000HP5RC7', 'Metalac', true, 'RSMETAE71629', 'commonStock', 'RS', 'XBEL', 'MTLC', 2040000, 'Utensils'),
('BBG0015L55D4', 'NIS', true, 'RSNISHE79420', 'commonStock', 'RS', 'XBEL', 'NIIS', 163060400, 'Oil'),
('BBG000HGH3F4', 'Impol Seval', true, 'RSIMPLE20713', 'commonStock', 'RS', 'XBEL', 'IMPL', 942287, 'Metallurgy')

-- +goose Down

