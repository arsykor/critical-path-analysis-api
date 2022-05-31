package composites

import (
	"critical-path-analysis-api/internal/adapters/api"
	taskApi "critical-path-analysis-api/internal/adapters/api/task"
	taskDb "critical-path-analysis-api/internal/adapters/db/task"
	taskDom "critical-path-analysis-api/internal/domain/task"
	"critical-path-analysis-api/pkg/client/postgresql"
)

type TaskComposite struct {
	Storage taskDom.Storage
	Service taskDom.Service
	Handler api.Handler
}

func NewTaskComposite(postgresqlClient postgresql.Client) *TaskComposite {
	storage := taskDb.NewStorage(postgresqlClient)
	service := taskDom.NewService(storage)
	handler := taskApi.NewHandler(service)

	return &TaskComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}
}
