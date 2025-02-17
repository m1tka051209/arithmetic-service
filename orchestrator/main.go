package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mitka051209/arithmetic-service/orchestrator/handlers"
	"github.com/mitka051209/arithmetic-service/orchestrator/task_manager"
)

func main() {
	r := mux.NewRouter()
	taskManager := task_manager.NewTaskManager() // Initialize the task manager

	// Endpoints
	r.HandleFunc("/api/v1/calculate", handlers.CalculateHandler(taskManager)).Methods("POST")
	r.HandleFunc("/api/v1/expressions", handlers.ListExpressionsHandler(taskManager)).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", handlers.GetExpressionHandler(taskManager)).Methods("GET")
	r.HandleFunc("/internal/task", handlers.GetTaskHandler(taskManager)).Methods("GET")
	r.HandleFunc("/internal/task", handlers.SubmitResultHandler(taskManager)).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port
	}
	log.Printf("Orchestrator listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
