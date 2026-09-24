package aumapi

import "time"

type AumDto struct {
	TotalAum float64   `json:"totalAUM"`
	Currency string    `json:"currency"`
	Date     time.Time `json:"date"`
}
