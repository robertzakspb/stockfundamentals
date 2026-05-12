package jobs

import (
	"net/http"

	jobservice "github.com/compoundinvest/stockfundamentals/internal/application/jobs"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func StartAllJobs(c *gin.Context) {
	jobservice.StartDailyJobs()
	jobservice.StartHeavyJobs()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "All jobs were successfully started/executed"})
}

func StartDailyJobs(c *gin.Context) {
	jobservice.StartDailyJobs()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Daily jobs were successfully started"})
}

func StartHeavyJobs(c *gin.Context) {
	jobservice.StartHeavyJobs()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Heavy jobs were successfully started/executed"})
}
