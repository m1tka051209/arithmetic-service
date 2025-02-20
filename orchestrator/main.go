package main

import (
	"log"
	"net/http"

	"github.com/m1tka051209/arithmetic-service/orchestrator"
)

func main() {
	http.HandleFunc("/api/v1/calculate", orchestrator.CalculateHandler)
	http.HandleFunc("/api/v1/expressions", orchestrator.ExpressionsHandler)
	http.HandleFunc("/api/v1/expressions/", orchestrator.ExpressionByIDHandler)
	http.HandleFunc("/internal/task", orchestrator.TaskHandler)

	log.Println("Orchestrator запущен на порту :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
