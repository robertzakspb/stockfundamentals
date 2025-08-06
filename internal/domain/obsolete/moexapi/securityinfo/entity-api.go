//DISCLAIMER: THIS PACKAGE IS NO LONGER IN USE, IT MAY EVENTUALLY BE DELETED
package securityinfo

type MoexSecuritiesDTO struct {
	Securities struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"securities"`
}
