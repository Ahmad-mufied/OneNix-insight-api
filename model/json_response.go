package model

type JSONResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
	Data    any    `json:"data,omitempty"`
}
