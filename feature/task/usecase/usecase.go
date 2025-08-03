package usecase

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/g12-todo-backend/domain"
	"github.com/pubestpubest/g12-todo-backend/models"
	"github.com/pubestpubest/g12-todo-backend/request"
	"github.com/pubestpubest/g12-todo-backend/response"
)

type taskUsecase struct {
	taskRepository domain.TaskRepository
}

func NewTaskUsecase(taskRepository domain.TaskRepository) domain.TaskUsecase {
	return &taskUsecase{taskRepository: taskRepository}
}

func (u *taskUsecase) GetTaskList(page, limit int) (*response.PaginatedResponse[*response.TaskResponse], error) {
	tasks, total, err := u.taskRepository.GetTaskList(page, limit)
	if err != nil {
		return nil, errors.Wrap(err, "[TaskUsecase.GetTaskList]: Error getting task list")
	}

	var taskResponses []*response.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, &response.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &response.PaginatedResponse[*response.TaskResponse]{
		Data: taskResponses,
		Pagination: response.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (u *taskUsecase) GetTaskByID(id uint64) (*response.TaskResponse, error) {
	task, err := u.taskRepository.GetTaskByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[TaskUsecase.GetTaskByID]: Error getting task")
	}

	if task == nil {
		return nil, errors.New("[TaskUsecase.GetTaskByID]: Task not found")
	}

	return &response.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (u *taskUsecase) CreateTask(req *request.TaskRequest) (*response.TaskResponse, error) {
	task := &models.Task{
		Title:       req.Title,
		Description: &req.Description,
		Status:      req.Status,
	}

	if err := u.taskRepository.CreateTask(task); err != nil {
		return nil, errors.Wrap(err, "[TaskUsecase.CreateTask]: Error creating task")
	}

	return &response.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (u *taskUsecase) UpdateTask(id uint64, req *request.TaskRequest) (*response.TaskResponse, error) {
	task, err := u.taskRepository.GetTaskByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[TaskUsecase.UpdateTask]: Error getting task")
	}

	if task == nil {
		return nil, errors.New("[TaskUsecase.UpdateTask]: Task not found")
	}

	task.Title = req.Title
	task.Description = &req.Description
	task.Status = req.Status

	if err := u.taskRepository.UpdateTask(task); err != nil {
		return nil, errors.Wrap(err, "[TaskUsecase.UpdateTask]: Error updating task")
	}

	return &response.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (u *taskUsecase) DeleteTask(id uint64) error {
	if err := u.taskRepository.DeleteTask(id); err != nil {
		return errors.Wrap(err, "[TaskUsecase.DeleteTask]: Error deleting task")
	}
	return nil
}
