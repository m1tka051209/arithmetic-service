package models

import "time"

type Task struct {
    ID            string        `json:"id"`
    Arg1          float64       `json:"arg1"`
    Arg2          float64       `json:"arg2"`
    Operation     string        `json:"operation"`
    OperationTime time.Duration `json:"-"`
    Status        string        `json:"status"`
    Result        float64       `json:"result,omitempty"`
    ExpressionID  string        `json:"expression_id"`
}

func (t Task) GetOperationTimeMS() int {
    return int(t.OperationTime.Milliseconds())
}