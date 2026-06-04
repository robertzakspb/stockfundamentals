package transaction

func (t *Transaction) IsBuyOrder() bool {
	return t.Side == Buy
}

func (t *Transaction) IsSellOrder() bool {
	return t.Side == Sell
}

func (t *Transaction) IsDepositOrWithdrawal() bool {
	return t.Type == Deposit || t.Type == Withdrawal
}

func (t *Transaction) IsDeposit() bool {
	return t.Type == Deposit
}

func (t *Transaction) IsWithdrawal() bool {
	return t.Type == Withdrawal
}
