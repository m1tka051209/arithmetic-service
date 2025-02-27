package models

type Expression struct {
    ID     string `json:"id"`
    Status string `json:"status"`
    Result float64 `json:"result,omitempty"`
}