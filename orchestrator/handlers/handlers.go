package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mitka051209/arithmetic-service/orchestrator/models"
	"github.com/mitka051209/arithmetic-service/orchestrator/task_manager"
)

// Dependency injection for the TaskManager
func CalculateHandler(taskManager *task_manager.TaskManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Expression string `json:"expression"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
			return
		}

		if requestBody.Expression == "" {
			http.Error(w, "Expression is required", http.StatusUnprocessableEntity)
			return
		}

		// Generate a unique ID for the expression
		expressionID := uuid.New().String()

		// Add the expression to the task manager
		expression := &models.Expression{
			ID:     expressionID,
			Status: "pending",
			Result: nil, // Initialize result to nil
		}

		taskManager.AddExpression(expression) // Method to add the expression to the system

		// Mock parsing expression to tasks
		var arg1, arg2 float64 = 2, 2
		task := &models.Task{
			ID:            uuid.New().String(),
			ExpressionID:  expressionID,
			Arg1:          &arg1,
			Arg2:          &arg2,
			Operation:     "+",
			OperationTime: 1000,
		}
		taskManager.AddTask(task)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": expressionID})
	}
}

func ListExpressionsHandler(taskManager *task_manager.TaskManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expressions := taskManager.GetAllExpressions()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]models.Expression{"expressions": expressions})
	}
}

func GetExpressionHandler(taskManager *task_manager.TaskManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		expressionID := vars["id"]

		expression, found := taskManager.GetExpression(expressionID)
		if !found {
			http.Error(w, "Expression not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"expression": expression})
	}
}

func GetTaskHandler(taskManager *task_manager.TaskManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, found := taskManager.GetNextTask()
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
	}
}

func SubmitResultHandler(taskManager *task_manager.TaskManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var taskResult struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}

		if err := json.NewDecoder(r.Body).Decode(&taskResult); err != nil {
			http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
			return
		}

		err := taskManager.CompleteTask(taskResult.ID, taskResult.Result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound) // Or another appropriate status
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
