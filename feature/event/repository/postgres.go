package repository

import (
	"time"

	"github.com/pkg/errors"
	"github.com/pubestpubest/g12-todo-backend/domain"
	"github.com/pubestpubest/g12-todo-backend/models"
	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) domain.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) GetEventList(page, limit int) ([]*models.Events, int64, error) {
	var events []*models.Events
	var total int64

	// Get total count
	if err := r.db.Model(&models.Events{}).Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, "[EventRepository.GetEventList]: Error counting events")
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Find(&events).Error; err != nil {
		return nil, 0, errors.Wrap(err, "[EventRepository.GetEventList]: Error getting event list")
	}

	return events, total, nil
}

func (r *eventRepository) GetEventByID(id uint64) (*models.Events, error) {
	var event models.Events
	if err := r.db.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, errors.Wrap(err, "[EventRepository.GetEventByID]: Error getting event")
	}
	return &event, nil
}

func (r *eventRepository) CreateEvent(event *models.Events) error {
	now := time.Now()
	event.CreatedAt = &now
	event.UpdatedAt = &now

	if err := r.db.Create(event).Error; err != nil {
		return errors.Wrap(err, "[EventRepository.CreateEvent]: Error creating event")
	}
	return nil
}

func (r *eventRepository) UpdateEvent(event *models.Events) error {
	now := time.Now()
	event.UpdatedAt = &now

	if err := r.db.Save(event).Error; err != nil {
		return errors.Wrap(err, "[EventRepository.UpdateEvent]: Error updating event")
	}
	return nil
}

func (r *eventRepository) DeleteEvent(id uint64) error {
	if err := r.db.Delete(&models.Events{}, id).Error; err != nil {
		return errors.Wrap(err, "[EventRepository.DeleteEvent]: Error soft deleting event")
	}
	return nil
}
