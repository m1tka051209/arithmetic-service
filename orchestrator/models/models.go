package models

type Expression struct {
	ID     string          `json:"id"`
	Status string          `json:"status"` // "pending", "processing", "completed", "error"
	Result *float64        `json:"result"`
}

type Task struct {
	ID            string      `json:"id"`
	ExpressionID  string      `json:"expression_id"`
	Arg1          *float64    `json:"arg1"`
	Arg2          *float64    `json:"arg2"`
	Operation     string      `json:"operation"` // "+", "-", "*", "/"
	OperationTime int         `json:"operation_time"` // in milliseconds
}