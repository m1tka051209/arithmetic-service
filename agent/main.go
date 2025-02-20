package main

import (
	"log"
	"os"

	"github.com/mitka051209/arithmetic-service/agent/worker"
)

func main() {
	log.Println("Агент запущен")
	worker.StartWorkers(os.Getenv("COMPUTING_POWER"))
}
