package api

import (
    "encoding/json"
    "net/http"
    "fmt"
    "github.com/m1tka051209/arithmetic-service/orchestrator/task_manager"
    "github.com/m1tka051209/arithmetic-service/orchestrator/models"
)

func validateTasks(tasks []models.Task) error {
    for _, task := range tasks {
        if task.Operation == "/" && task.Arg2 == 0 {
            return fmt.Errorf("деление на ноль в задаче %s", task.ID)
        }
    }
    return nil
}

type Handlers struct {
    tm *task_manager.TaskManager
}

func NewHandlers(tm *task_manager.TaskManager) *Handlers {
    return &Handlers{tm: tm}
}

func (h *Handlers) CalculateHandler(w http.ResponseWriter, r *http.Request) {
    var req struct { Expression string `json:"expression"` }
    
    // Чтение и валидация запроса
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusUnprocessableEntity, "неверный формат запроса")
        return
    }

    // Парсинг выражения
    tasks, err := h.tm.ParseExpression(req.Expression)
    if err != nil {
        respondWithError(w, http.StatusUnprocessableEntity, err.Error())
        return
    }

    // Проверка деления на ноль
    if err := validateTasks(tasks); err != nil {
        respondWithError(w, http.StatusUnprocessableEntity, err.Error())
        return
    }

    // Сохранение выражения
    exprID := h.tm.GenerateID()
    h.tm.SaveExpression(exprID, tasks)

    // Успешный ответ
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"id": exprID})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *Handlers) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
    expressions := h.tm.GetAllExpressions()
    response := map[string]interface{}{
        "expressions": expressions,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}