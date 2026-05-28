package accountservice

import (
	"errors"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	"github.com/google/uuid"
)

func FindAccountById(id uuid.UUID, accounts []account.Account) (account.Account, error) {
	for i := range accounts {
		if accounts[i].Id == id {
			return accounts[i], nil
		}
	}

	return account.Account{}, errors.New("Failed to find an account with ID " + id.String() + " in the provided list")
}
