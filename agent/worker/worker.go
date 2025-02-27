package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Task struct {
    ID            string  `json:"id"`
    Arg1          float64 `json:"arg1"`
    Arg2          float64 `json:"arg2"`
    Operation     string  `json:"operation"`
    OperationTime int     `json:"operation_time"`
}

func StartWorkers(power int) {
    for i := 0; i < power; i++ {
        go func(workerID int) {
            for {
                task, err := getTask()
                if err != nil {
                    log.Printf("Worker %d: Error getting task: %v", workerID, err)
                    time.Sleep(1 * time.Second)
                    continue
                }

                result := calculate(task)
                err = submitResult(task.ID, result)
                if err != nil {
                    log.Printf("Worker %d: Error submitting result: %v", workerID, err)
                }
            }
        }(i)
    }
}

func getTask() (Task, error) {
    resp, err := http.Get("http://localhost:8080/internal/task")
    if err != nil {
        return Task{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        return Task{}, fmt.Errorf("no tasks available")
    }

    var taskResp struct {
        Task Task `json:"task"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
        return Task{}, err
    }

    return taskResp.Task, nil
}

func calculate(task Task) float64 {
    switch task.Operation {
    case "+":
        return task.Arg1 + task.Arg2
    case "-":
        return task.Arg1 - task.Arg2
    case "*":
        return task.Arg1 * task.Arg2
    case "/":
        return task.Arg1 / task.Arg2
    default:
        return task.Arg1
    }
}

func submitResult(taskID string, result float64) error {
    payload := struct {
        ID     string  `json:"id"`
        Result float64 `json:"result"`
    }{
        ID:     taskID,
        Result: result,
    }

    jsonPayload, _ := json.Marshal(payload)
    resp, err := http.Post(
        "http://localhost:8080/internal/task",
        "application/json",
        bytes.NewBuffer(jsonPayload),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return nil
}