package bondportfolioapi

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"

	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func AddBondPositionLotToAccount(c *gin.Context) {
	bodyReader := c.Request.Body
	defer bodyReader.Close()

	jsonData, err := io.ReadAll(bodyReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		logger.Log("Failed to read the bond position lot json from the POST payload: "+err.Error(), logger.ERROR)
		return
	}

	dto := bondPositionLotDto{}
	err = json.Unmarshal(jsonData, &dto)
	if err != nil {
		logger.Log("Failed to unmarshal the dividend forecast json in the POST payload: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	domain := mapBondLotDtoToDomain(dto)

	err = bondportfolio.SaveBondPositionLot(domain)

	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, shared.StringResponse{Message: "The position lot has been successfully saved to the database"})
}

func GetAccountPositionLots(c *gin.Context) {
	filters, err := ydbfilter.MapQueryFiltersToYdb[bondPositionLotDto](c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	lots, err := bondportfolio.GetFilteredPositionLots(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	lots, err = bondportfolio.CalculateYtmForLots(lots)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	mappedLots := []bondPositionLotDto{}
	for _, lot := range lots {
		mappedLot := mapBondLotToDto(lot)
		mappedLots = append(mappedLots, mappedLot)
	}

	slices.SortFunc(mappedLots, func(dto1, dto2 bondPositionLotDto) int {
		if dto1.Ytm > dto2.Ytm {
			return -1
		} else {
			return 1
		}
	})

	c.JSON(http.StatusOK, mappedLots)
}

func GetAccountBondTimeline(c *gin.Context) {
	items, err := bondportfolio.GetAccountTimeline()
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	dtos := mapTimeLineItemsToDtos(items)
	c.JSON(http.StatusOK, dtos)
}
