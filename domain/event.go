package domain

import (
	"github.com/pubestpubest/g12-todo-backend/models"
	"github.com/pubestpubest/g12-todo-backend/request"
	"github.com/pubestpubest/g12-todo-backend/response"
)

type EventUsecase interface {
	GetEventList(page, limit int) (*response.PaginatedResponse[*response.EventResponse], error)
	GetEventByID(id uint64) (*response.EventResponse, error)
	CreateEvent(req *request.EventRequest) (*response.EventResponse, error)
	UpdateEvent(id uint64, req *request.EventRequest) (*response.EventResponse, error)
	DeleteEvent(id uint64) error
}

type EventRepository interface {
	GetEventList(page, limit int) ([]*models.Events, int64, error)
	GetEventByID(id uint64) (*models.Events, error)
	CreateEvent(event *models.Events) error
	UpdateEvent(event *models.Events) error
	DeleteEvent(id uint64) error
}
