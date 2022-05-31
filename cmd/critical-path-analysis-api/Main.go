package main

import (
	"context"
	"critical-path-analysis-api/internal/composites"
	"critical-path-analysis-api/internal/config"
	"critical-path-analysis-api/internal/tests"
	"critical-path-analysis-api/pkg/client/postgresql"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	tests.InitTestRepository()

	conf := config.NewConfig()

	postgresqlClient, err := postgresql.NewClient(context.Background(), 5, conf)
	if err != nil {
		fmt.Println(err)
	}

	taskComp := composites.NewTaskComposite(postgresqlClient)
	router := gin.Default()
	taskComp.Handler.Register(router, "task")

	router.Run("localhost:8080")
}
