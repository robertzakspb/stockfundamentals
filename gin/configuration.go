package ginconfig

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

const (
	DBHost = "localhost"
	DBPort = 5432
	DBUser = "postgres"
	DBName = "security_fundamentals"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Database(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}
