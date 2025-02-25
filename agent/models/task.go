package models

// import "time"

// Task - структура задачи для вычислительного агента
type Task struct {
    ID            string    `json:"id"`
    Arg1          float64   `json:"arg1"`
    Arg2          float64   `json:"arg2"`
    Operation     string    `json:"operation"`
    OperationTime int       `json:"operation_time"` // Время выполнения в миллисекундах
}