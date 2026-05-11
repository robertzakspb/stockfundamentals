package shared

type ResponseError struct {
	Errors  []string `json:"errors"`
	Message string   `json:"message"`
}

type StringResponse struct {
	Message string `json:"message"`
}

type NumberResponse struct {
	Value any `json:"value"`
}
