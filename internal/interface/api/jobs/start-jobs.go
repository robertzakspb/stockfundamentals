package jobs

import (
	"net/http"

	jobservice "github.com/compoundinvest/stockfundamentals/internal/application/jobs"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func StartAllJobs(c *gin.Context) {
	go func() {
		jobservice.StartDailyJobs()
		jobservice.StartHeavyJobs()
	}()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "All jobs were successfully started/executed"})
}

func StartDailyJobs(c *gin.Context) {
	go jobservice.StartDailyJobs()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Daily jobs were successfully started"})
}

func StartHeavyJobs(c *gin.Context) {
	go jobservice.StartHeavyJobs()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Heavy jobs were successfully started/executed"})
}
