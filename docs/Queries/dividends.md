To fetch the dividend payments of a particular company, execute the following query:

```sql
SELECT company_name, ticker, payout_date, actual_DPS/1000000, currency,
FROM `stockfundamentals/stocks/dividend_payment`
JOIN `stockfundamentals/stocks/stock` ON `stockfundamentals/stocks/stock`.id = `stockfundamentals/stocks/dividend_payment`.stock_id
WHERE ticker = 'NVTK'
ORDER BY payout_date DESC
```
