package models

import "time"

type Task struct {
    ID            string
    Arg1          float64
    Arg2          float64
    Operation     string // "+", "-", "*", "/"
    OperationTime time.Duration
    Status        string // "pending", "in_progress", "completed"
    Result        float64
}