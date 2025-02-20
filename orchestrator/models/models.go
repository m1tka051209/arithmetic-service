package models

import "time"

type Expression struct {
    ID     string
    Status string
    Result float64
}

type Task struct {
    ID            string
    Arg1          float64
    Arg2          float64
    Operation     string
    OperationTime time.Duration
    Status        string
    Result        float64
}