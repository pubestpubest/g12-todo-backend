package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/pubestpubest/g12-todo-backend/constant"
	"github.com/pubestpubest/g12-todo-backend/domain"
	"github.com/pubestpubest/g12-todo-backend/request"
	"github.com/pubestpubest/g12-todo-backend/response"
	"github.com/pubestpubest/g12-todo-backend/utils"
	log "github.com/sirupsen/logrus"
)

type taskHandler struct {
	taskUsecase domain.TaskUsecase
}

func NewTaskHandler(taskUsecase domain.TaskUsecase) *taskHandler {
	return &taskHandler{taskUsecase: taskUsecase}
}

func (h *taskHandler) GetTaskList(c *gin.Context) {
	var paginationReq request.PaginationRequest

	// Set default values
	paginationReq.Page = 1
	paginationReq.Limit = 10

	// Bind query parameters
	if err := c.ShouldBindQuery(&paginationReq); err != nil {
		err = errors.Wrap(err, "[TaskHandler.GetTaskList]: Error binding query parameters")
		log.Warn(err)
		resp := response.PaginatedResponse[interface{}]{
			Status:  constant.Failed,
			Message: utils.StandardError(err),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	tasks, err := h.taskUsecase.GetTaskList(paginationReq.Page, paginationReq.Limit)
	if err != nil {
		err = errors.Wrap(err, "[TaskHandler.GetTaskList]: Error getting task list")
		log.Error(err)
		resp := response.PaginatedResponse[interface{}]{
			Status:  constant.Failed,
			Message: utils.StandardError(err),
			Data:    nil,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := response.PaginatedResponse[*response.TaskResponse]{
		Status:     constant.Success,
		Message:    "List tasks successfully",
		Data:       tasks.Data,
		Pagination: tasks.Pagination,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *taskHandler) GetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.taskUsecase.GetTaskByID(id)
	if err != nil {
		err = errors.Wrap(err, "[TaskHandler.GetTaskByID]: Error getting task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}

	resp := response.Response[*response.TaskResponse]{
		Status:  constant.Success,
		Message: "Task retrieved successfully",
		Data:    task,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *taskHandler) CreateTask(c *gin.Context) {
	var req request.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	task, err := h.taskUsecase.CreateTask(&req)
	if err != nil {
		err = errors.Wrap(err, "[TaskHandler.CreateTask]: Error creating task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}

	resp := response.Response[*response.TaskResponse]{
		Status:  constant.Success,
		Message: "Task created successfully",
		Data:    task,
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *taskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req request.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	task, err := h.taskUsecase.UpdateTask(id, &req)
	if err != nil {
		err = errors.Wrap(err, "[TaskHandler.UpdateTask]: Error updating task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}

	resp := response.Response[*response.TaskResponse]{
		Status:  constant.Success,
		Message: "Task updated successfully",
		Data:    task,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *taskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := h.taskUsecase.DeleteTask(id); err != nil {
		err = errors.Wrap(err, "[TaskHandler.DeleteTask]: Error deleting task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}

	resp := response.Response[interface{}]{
		Status:  constant.Success,
		Message: "Task deleted successfully",
		Data:    nil,
	}
	c.JSON(http.StatusOK, resp)
}
