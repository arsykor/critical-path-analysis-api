package main

import (
	"critical-path-analysis-api/internal/composites"
	"critical-path-analysis-api/internal/tests"
	"github.com/gin-gonic/gin"
)

func main() {
	tests.InitTestRepository()

	taskComp := composites.NewTaskComposite()
	router := gin.Default()
	taskComp.Handler.Register(router, "task")

	router.Run("localhost:8080")
}
