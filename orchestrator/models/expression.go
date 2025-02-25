package models

type Expression struct {
    ID     string
    Status string // "pending", "completed", "error"
    Result float64
}