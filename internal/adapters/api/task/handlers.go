package task

import (
	"critical-path-analysis-api/internal/adapters/api"
	"critical-path-analysis-api/internal/domain/task"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

type handler struct {
	service task.Service
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHandler(service task.Service) api.Handler {
	return &handler{service: service}
}

func (h *handler) Register(router *gin.Engine, rootName string) {
	tasksRoot := router.Group("/" + rootName)
	{
		tasksRoot.POST("/create", h.CreateTask)
		tasksRoot.GET("/get/:id", h.GetTask)
		tasksRoot.GET("/get", h.GetAllTasks)
		tasksRoot.POST("/update", h.UpdateTask)
		tasksRoot.GET("/delete/:id", h.DeleteTask)
	}
}

func (h *handler) CreateTask(context *gin.Context) {
	var tasks []task.Task

	jsonData, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		returnError(context, err)
		return
	}

	err = json.Unmarshal(jsonData, &tasks)
	if err != nil {
		returnError(context, err)
		return
	}

	tasksOut, err := h.service.Create(&tasks)
	if err != nil {
		returnError(context, err)
		return
	}

	context.IndentedJSON(http.StatusOK, tasksOut)
}

func (h *handler) GetTask(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		returnError(context, err)
		return
	}

	task, err := h.service.GetById(id)
	if err != nil {
		returnError(context, err)
		return
	}

	context.IndentedJSON(http.StatusOK, task)
}

func (h *handler) GetAllTasks(context *gin.Context) {
	tasks, err := h.service.GetAll()
	if err != nil {
		returnError(context, err)
		return
	}

	context.IndentedJSON(http.StatusOK, tasks)
}

func (h *handler) UpdateTask(context *gin.Context) {
	var tasks []task.Task

	jsonData, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		returnError(context, err)
		return
	}

	err = json.Unmarshal(jsonData, &tasks)
	if err != nil {
		returnError(context, err)
		return
	}

	arrangedTasks, err := h.service.Update(&tasks[0])
	if err != nil {
		returnError(context, err)
		return
	}

	context.IndentedJSON(http.StatusOK, arrangedTasks)
}

func (h *handler) DeleteTask(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		returnError(context, err)
		return
	}

	arrangedTasks, err := h.service.Delete(id)
	if err != nil {
		returnError(context, err)
		return
	}

	context.IndentedJSON(http.StatusOK, arrangedTasks)
}

func returnError(context *gin.Context, err error) {
	context.JSON(http.StatusBadRequest, Error{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	})
}
