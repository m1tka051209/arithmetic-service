package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/m1tka051209/arithmetic-service/orchestrator/models"
	"github.com/m1tka051209/arithmetic-service/orchestrator/task_manager"
)

type Handlers struct {
    tm *task_manager.TaskManager
}

func NewHandlers(tm *task_manager.TaskManager) *Handlers {
    return &Handlers{tm: tm}
}


// CalculateHandler — добавление нового выражения
func (h *Handlers) CalculateHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Expression string `json:"expression"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusUnprocessableEntity, "invalid request body")
        return
    }

    tasks, err := h.tm.ParseExpression(req.Expression)
    if err != nil {
        h.respondError(w, http.StatusUnprocessableEntity, err.Error())
        return
    }

    exprID := h.tm.GenerateID()
    h.tm.SaveExpression(exprID, tasks)

    h.respondJSON(w, http.StatusCreated, map[string]string{"id": exprID})
}

func (h *Handlers) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
    expressions := h.tm.GetAllExpressions()
    h.respondJSON(w, http.StatusOK, map[string][]models.Expression{
        "expressions": expressions, // Ключ "expressions" в JSON
    })
}

// GetExpressionHandler — получение выражения по ID
func (h *Handlers) GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
    expr, exists := h.tm.GetExpressionByID(id)
    if !exists {
        h.respondError(w, http.StatusNotFound, "expression not found")
        return
    }
    h.respondJSON(w, http.StatusOK, map[string]models.Expression{"expression": expr})
}

func (h *Handlers) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
    task, exists := h.tm.GetNextTask()
    if !exists {
        h.respondError(w, http.StatusNotFound, "no tasks available")
        return
    }

    // Преобразуем задачу в требуемый формат ответа
    response := struct {
        Task struct {
            ID            string  `json:"id"`
            Arg1          float64 `json:"arg1"`
            Arg2          float64 `json:"arg2"`
            Operation     string  `json:"operation"`
            OperationTime int     `json:"operation_time"`
        } `json:"task"`
    }{
        Task: struct {
            ID            string  `json:"id"`
            Arg1          float64 `json:"arg1"`
            Arg2          float64 `json:"arg2"`
            Operation     string  `json:"operation"`
            OperationTime int     `json:"operation_time"`
        }{
            ID:            task.ID,
            Arg1:          task.Arg1,
            Arg2:          task.Arg2,
            Operation:     task.Operation,
            OperationTime: int(task.OperationTime.Milliseconds()),
        },
    }

    h.respondJSON(w, http.StatusOK, response)
}

// SubmitResultHandler — прием результата от агента
func (h *Handlers) SubmitResultHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ID     string  `json:"id"`
        Result float64 `json:"result"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusUnprocessableEntity, "invalid request body")
        return
    }

    h.tm.SaveTaskResult(req.ID, req.Result)
    w.WriteHeader(http.StatusOK)
}

// Вспомогательные методы
func (h *Handlers) respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Printf("JSON encode error: %v", err)
    }
}

func (h *Handlers) respondError(w http.ResponseWriter, status int, message string) {
    h.respondJSON(w, status, map[string]string{"error": message})
}