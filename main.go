package main

import (
	"database/sql"
	"fmt"

	"github.com/compoundinvest/stockfundamentals/Features/company"
	"github.com/compoundinvest/stockfundamentals/Features/portfolio/lot"
	ginconfig "github.com/compoundinvest/stockfundamentals/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+" dbname=%s sslmode=disable", ginconfig.DBHost, ginconfig.DBPort, ginconfig.DBUser, ginconfig.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Failed to initilalize the database: ", err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(ginconfig.CORSMiddleware())
	router.Use(ginconfig.Database(db))

	router.GET("/lots", lot.GetLots)
	router.GET("/security", company.GetCompanyFromDB)

	router.Run("localhost:8080")
}
