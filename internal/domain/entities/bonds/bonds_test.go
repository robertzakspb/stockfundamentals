package bonds

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_Validate(t *testing.T) {
	var id uuid.UUID
	bond := Bond{
		Id: id,
	}

	err := bond.Validate()
	test.AssertError(t, err)
}
