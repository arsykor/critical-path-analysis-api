package main

import (
	"context"
	"critical-path-analysis-api/internal/composites"
	"critical-path-analysis-api/internal/config"
	"critical-path-analysis-api/pkg/client/postgresql"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		return
	}

	postgresqlClient, err := postgresql.NewClient(context.Background(), 5, conf)
	if err != nil {
		fmt.Println(err)
	}

	taskComp := composites.NewTaskComposite(postgresqlClient)

	router := gin.Default()
	taskComp.Handler.Register(router, "task")
	addr := fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port)
	router.Run(addr)
}
