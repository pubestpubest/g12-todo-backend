package repository

import (
	"time"

	"github.com/pkg/errors"
	"github.com/pubestpubest/g12-todo-backend/domain"
	"github.com/pubestpubest/g12-todo-backend/models"
	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) GetTaskList(page, limit int) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var total int64

	// Get total count
	if err := r.db.Model(&models.Task{}).Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "[TaskRepository.GetTaskList]: Error counting tasks")
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, errors.Wrap(err, "[TaskRepository.GetTaskList]: Error getting task list")
	}

	return tasks, total, nil
}

func (r *taskRepository) GetTaskByID(id uint64) (*models.Task, error) {
	var task models.Task
	if err := r.db.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, errors.Wrap(err, "[TaskRepository.GetTaskByID]: Error getting task")
	}
	return &task, nil
}

func (r *taskRepository) CreateTask(task *models.Task) error {
	now := time.Now()
	task.CreatedAt = &now
	task.UpdatedAt = &now

	if err := r.db.Create(task).Error; err != nil {
		return errors.Wrap(err, "[TaskRepository.CreateTask]: Error creating task")
	}
	return nil
}

func (r *taskRepository) UpdateTask(task *models.Task) error {
	now := time.Now()
	task.UpdatedAt = &now

	if err := r.db.Save(task).Error; err != nil {
		return errors.Wrap(err, "[TaskRepository.UpdateTask]: Error updating task")
	}
	return nil
}

func (r *taskRepository) DeleteTask(id uint64) error {
	if err := r.db.Delete(&models.Task{}, id).Error; err != nil {
		return errors.Wrap(err, "[TaskRepository.DeleteTask]: Error soft deleting task")
	}
	return nil
}
