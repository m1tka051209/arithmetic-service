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
					log.Printf("Worker %d: %v", workerID, err)
					time.Sleep(2 * time.Second)
					continue
				}

				log.Printf("Worker %d: Processing task %s", workerID, task.ID)
				time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
				result := calculate(task)

				if err := submitResult(task.ID, result); err != nil {
					log.Printf("Worker %d: Submit error: %v", workerID, err)
				}
			}
		}(i)
	}
}

func getTask() (Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return Task{}, fmt.Errorf("failed to fetch task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Task{}, fmt.Errorf("no tasks available")
	}

	var response struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Task{}, fmt.Errorf("failed to decode task: %w", err)
	}

	return response.Task, nil
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
		if task.Arg2 == 0 {
			return 0
		}
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

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(
		"http://localhost:8080/internal/task",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("post failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}