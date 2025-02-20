package orchestrator

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/m1tka051209/arithmetic-service/orchestrator/task_manager"
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusUnprocessableEntity)
		return
	}

	tasks, err := task_manager.ParseExpression(req.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exprID := task_manager.GenerateID()
	task_manager.SaveExpression(exprID, tasks)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": exprID})
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	expressions := task_manager.GetAllExpressions()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": expressions})
}

func ExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	expr, exists := task_manager.GetExpressionByID(id)
	if !exists {
		http.Error(w, "Выражение не найдено", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		task, exists := task_manager.GetNextTask()
		if !exists {
			http.Error(w, "Нет доступных задач", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
	case http.MethodPost:
		var req struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusUnprocessableEntity)
			return
		}
		task_manager.SaveTaskResult(req.ID, req.Result)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
