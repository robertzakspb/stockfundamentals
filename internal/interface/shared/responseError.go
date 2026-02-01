package shared

type ResponseError struct {
	Errors []string `json:"errors"`
}