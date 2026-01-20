package models

// Define a estrutura do erro padronizado com base na RFC 9457.
type ErrorResponse struct {
	Type     string `json:"type"`
	Status   int    `json:"status"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
	Code     string `json:"code"`
}
