package worker

import (
	"time"

	"github.com/mitka051209/arithmetic-service/agent/client"
)

func StartWorkers(workers int) {
	for i := 0; i < workers; i++ {
		go processTasks()
	}
}

func processTasks() {
	for {
		task := client.FetchTask()
		if task != nil {
			result := calculate(task)
			client.SendResult(task.ID, result)
		}
		time.Sleep(1 * time.Second)
	}
}

func calculate(task *client.Task) float64 {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

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
		return 0
	}
}
