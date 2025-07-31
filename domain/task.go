package domain

import (
	"github.com/pubestpubest/g12-todo-backend/models"
	"github.com/pubestpubest/g12-todo-backend/request"
	"github.com/pubestpubest/g12-todo-backend/response"
)

type TaskUsecase interface {
	GetTaskList(page, limit int) (*response.PaginatedResponse[*response.TaskResponse], error)
	GetTaskByID(id uint64) (*response.TaskResponse, error)
	CreateTask(req *request.TaskRequest) (*response.TaskResponse, error)
	UpdateTask(id uint64, req *request.TaskRequest) (*response.TaskResponse, error)
	DeleteTask(id uint64) error
}

type TaskRepository interface {
	GetTaskList(page, limit int) ([]*models.Task, int64, error)
	GetTaskByID(id uint64) (*models.Task, error)
	CreateTask(task *models.Task) error
	UpdateTask(task *models.Task) error
	DeleteTask(id uint64) error
}
