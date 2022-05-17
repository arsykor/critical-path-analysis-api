package composites

import (
	"critical-path-analysis-api/internal/adapters/api"
	taskApi "critical-path-analysis-api/internal/adapters/api/task"
	taskDb "critical-path-analysis-api/internal/adapters/db/task"
	taskDom "critical-path-analysis-api/internal/domain/task"
)

type TaskComposite struct {
	Storage taskDom.Storage
	Service taskDom.Service
	Handler api.Handler
}

func NewTaskComposite() *TaskComposite {
	storage := taskDb.NewStorage()
	service := taskDom.NewService(storage)
	handler := taskApi.NewHandler(service)

	return &TaskComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}
}
