package api

import (
    "encoding/json"
    "net/http"
    
    "github.com/m1tka051209/arithmetic-service/orchestrator/task_manager"
)

type Handlers struct {
    tm *task_manager.TaskManager
}

func NewHandlers(tm *task_manager.TaskManager) *Handlers {
    return &Handlers{tm: tm}
}

func (h *Handlers) CalculateHandler(w http.ResponseWriter, r *http.Request) {
    var req struct { 
        Expression string `json:"expression"` 
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
        return
    }

    tasks, err := h.tm.ParseExpression(req.Expression)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    exprID := h.tm.GenerateID()
    h.tm.SaveExpression(exprID, tasks)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"id": exprID})
}

func (h *Handlers) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
    expressions := h.tm.GetAllExpressions()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{"expressions": expressions})
}