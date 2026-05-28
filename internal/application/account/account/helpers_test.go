package accountservice

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_FindAccountById_Positive(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	accounts := []account.Account{
		{Id: id1},
		{Id: id2},
	}

	targetAccount, err := FindAccountById(id2, accounts)

	test.AssertNoError(t, err)
	test.AssertEqual(t, id2,targetAccount.Id)
}

func Test_FindAccountById_Negative(t *testing.T) {
	id1, id2, id3 := uuid.New(), uuid.New(), uuid.New()
	accounts := []account.Account{
		{Id: id1},
		{Id: id2},
	}

	_, err := FindAccountById(id3, accounts)

	test.AssertError(t, err)
}
