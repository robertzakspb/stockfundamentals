package accountservice

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	accountdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_mapDbAccountsToAccounts(t *testing.T) {
	id := uuid.New()
	openingDate := time.Now()
	accountType := "IIS_1"
	broker := "VTB"
	holder := "John Doe"

	dbAccount := accountdb.AccountDbModel{
		Id:          id,
		OpeningDate: openingDate,
		Type:        accountType,
		Broker:      broker,
		Holder:      holder,
	}

	mappedAccount := mapDbAccountsToAccounts([]accountdb.AccountDbModel{dbAccount})[0]

	test.AssertEqual(t, id, mappedAccount.Id)
	test.AssertEqual(t, openingDate, mappedAccount.OpeningDate)
	test.AssertEqual(t, string(account.IIS_1), mappedAccount.Type)
	test.AssertEqual(t, broker, mappedAccount.Broker)
	test.AssertEqual(t, holder, mappedAccount.Holder)
}

func Test_mapDbAccountsToAccounts_IIS2(t *testing.T) {
	accountType := "IIS_2"

	dbAccount := accountdb.AccountDbModel{
		Type: accountType,
	}

	mappedAccount := mapDbAccountsToAccounts([]accountdb.AccountDbModel{dbAccount})[0]

	test.AssertEqual(t, string(account.IIS_2), mappedAccount.Type)
}

func Test_mapDbAccountsToAccounts_IIS3(t *testing.T) {
	accountType := "IIS_3"

	dbAccount := accountdb.AccountDbModel{
		Type: accountType,
	}

	mappedAccount := mapDbAccountsToAccounts([]accountdb.AccountDbModel{dbAccount})[0]

	test.AssertEqual(t, string(account.IIS_3), mappedAccount.Type)
}

func Test_mapDbAccountsToAccounts_STANDARD(t *testing.T) {
	accountType := "STANDARD"

	dbAccount := accountdb.AccountDbModel{
		Type: accountType,
	}

	mappedAccount := mapDbAccountsToAccounts([]accountdb.AccountDbModel{dbAccount})[0]

	test.AssertEqual(t, string(account.STANDARD), mappedAccount.Type)
}

