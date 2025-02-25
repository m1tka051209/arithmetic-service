package main

import (
    "log"
    "net/http"
    
    "github.com/m1tka051209/arithmetic-service/orchestrator/api"
    "github.com/m1tka051209/arithmetic-service/orchestrator/task_manager"
)

func main() {
    tm := task_manager.NewTaskManager()
    handlers := api.NewHandlers(tm)
    
    http.HandleFunc("/api/v1/calculate", handlers.CalculateHandler)
    http.HandleFunc("/api/v1/expressions", handlers.ExpressionsHandler) // Исправленное имя метода
    
    log.Println("🚀 Orchestrator запущен на :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}