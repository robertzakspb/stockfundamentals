package bondsapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/bond-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type bondPositionLotDto struct {
	Figi             string    `json:"figi"`
	Isin             string    `json:"isin"`
	OpeningDate      time.Time `json:"openingDate"`
	ModificationDate time.Time `json:"modificationDate"`
	AccountId        string    `json:"accountId"`
	Quantity         float64   `json:"quantity"`
	PricePerUnit     float64   `json:"pricePerUnit"`
}

func AddBondPositionLotToAccount(c *gin.Context) {
	bodyReader := c.Request.Body
	defer bodyReader.Close()

	jsonData, err := io.ReadAll(bodyReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		logger.Log("Failed to read the bond position lot json from the POST payload: "+err.Error(), logger.ERROR)
		return
	}

	dto := bondPositionLotDto{}
	err = json.Unmarshal(jsonData, &dto)
	if err != nil {
		logger.Log("Failed to unmarshal the dividend forecast json in the POST payload: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	domain := mapBondLotDtoToDomain(dto)

	err = bondportfolio.SaveBondPositionLot(domain)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, "The position lot has been successfully saved to the database")
	}
}

func GetAccountPositionLots(c *gin.Context) {
	lots, err := bondportfolio.GetAllPositionLots()

	mappedLots := []bondPositionLotDto{}
	for _, lot := range lots {
		mappedLot := mapBondLotToDto(lot)
		mappedLots = append(mappedLots, mappedLot)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, mappedLots)
	}
}

func mapBondLotDtoToDomain(dto bondPositionLotDto) bonds.BondLot {
	accountId, _ := uuid.Parse(dto.AccountId)
	domain := bonds.BondLot{
		Id:               uuid.New(),
		Figi:             dto.Figi,
		Isin:             dto.Isin,
		OpeningDate:      dto.OpeningDate,
		ModificationDate: dto.ModificationDate,
		AccountId:        accountId,
		Quantity:         dto.Quantity,
		PricePerUnit:     dto.PricePerUnit,
	}

	return domain
}

func mapBondLotToDto(lot bonds.BondLot) bondPositionLotDto {
	dto := bondPositionLotDto{
		Figi:             lot.Figi,
		Isin:             lot.Isin,
		OpeningDate:      lot.OpeningDate,
		ModificationDate: lot.ModificationDate,
		AccountId:        lot.AccountId.String(),
		Quantity:         lot.Quantity,
		PricePerUnit:     lot.PricePerUnit,
	}

	return dto
}
