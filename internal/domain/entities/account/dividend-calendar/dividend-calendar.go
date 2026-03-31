package dividendcalendar

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/google/uuid"
)

type DividendCalendar struct {
	AccountIds        uuid.UUIDs
	FuturePayouts []dividend.Payout
}
