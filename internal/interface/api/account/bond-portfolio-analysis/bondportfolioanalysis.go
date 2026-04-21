package bondportfolioanalysis

import (
	"net/http"

	bondportfolioanalysis "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio-analysis"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func GetAccountBondPortfolioOverview(c *gin.Context) {
	filters := []ydbfilter.YdbFilter{
		{
			YqlColumnName:  "account_id",
			Condition:      ydbfilter.Equal,
			ConditionValue: types.UuidValue(uuid.MustParse("129274f9-ee80-4e74-aa1c-fea578bac6e6")),
		},
	}
	accountBondPortfolioOverview, err := bondportfolioanalysis.GeneratePortfolioOverview(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, accountBondPortfolioOverview)
}
