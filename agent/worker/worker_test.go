package worker

import (
    "bytes"
    "encoding/json"
    "net/http"
)

// type Task struct {
//     ID            string  `json:"id"`
//     Arg1          float64 `json:"arg1"`
//     Arg2          float64 `json:"arg2"`
//     Operation     string  `json:"operation"`
//     OperationTime int     `json:"operation_time"`
// }

func FetchTask() *Task {
    resp, err := http.Get("http://localhost:8080/internal/task")
    if err != nil {
        return nil
    }
    defer resp.Body.Close()

    var task Task
    if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
        return nil
    }
    return &task
}

func SendResult(taskID string, result float64) {
    payload := map[string]interface{}{
        "id":     taskID,
        "result": result,
    }
    
    jsonData, _ := json.Marshal(payload)
    http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(jsonData))
}