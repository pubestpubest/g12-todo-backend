package usecase

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/g12-todo-backend/domain"
	"github.com/pubestpubest/g12-todo-backend/models"
	"github.com/pubestpubest/g12-todo-backend/request"
	"github.com/pubestpubest/g12-todo-backend/response"
)

type eventUsecase struct {
	eventRepository domain.EventRepository
}

func NewEventUsecase(eventRepository domain.EventRepository) domain.EventUsecase {
	return &eventUsecase{eventRepository: eventRepository}
}

func (u *eventUsecase) GetEventList(page, limit int) (*response.PaginatedResponse[*response.EventResponse], error) {
	events, total, err := u.eventRepository.GetEventList(page, limit)
	if err != nil {
		return nil, errors.Wrap(err, "[EventUsecase.GetEventList]: Error getting event list")
	}

	var eventResponses []*response.EventResponse
	for _, event := range events {
		eventResponses = append(eventResponses, &response.EventResponse{
			ID:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			Complete:    &event.Complete,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
			Location:    event.Location,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
		})
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &response.PaginatedResponse[*response.EventResponse]{
		Data: eventResponses,
		Pagination: response.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (u *eventUsecase) GetEventByID(id uint64) (*response.EventResponse, error) {
	event, err := u.eventRepository.GetEventByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[EventUsecase.GetEventByID]: Error getting event")
	}

	if event == nil {
		return nil, errors.New("event not found")
	}

	return &response.EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Complete:    &event.Complete,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
	}, nil
}

func (u *eventUsecase) CreateEvent(req *request.EventRequest) (*response.EventResponse, error) {

	if req.StartTime.After(req.EndTime) || req.StartTime.Equal(req.EndTime) {
		return nil, errors.New("[EventUsecase.CreateEvent]: startTime must be before endTime")
	}

	event := &models.Events{
		Title:       req.Title,
		Description: &req.Description,
		Complete:    *req.Complete,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	if err := u.eventRepository.CreateEvent(event); err != nil {
		return nil, errors.Wrap(err, "[EventUsecase.CreateEvent]: Error creating event")
	}

	return &response.EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Complete:    &event.Complete,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
	}, nil
}

func (u *eventUsecase) UpdateEvent(id uint64, req *request.EventRequest) (*response.EventResponse, error) {
	if req.StartTime.After(req.EndTime) || req.StartTime.Equal(req.EndTime) {
		return nil, errors.New("[EventUsecase.UpdateEvent]: startTime must be before endTime")
	}

	event, err := u.eventRepository.GetEventByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[EventUsecase.UpdateEvent]: Error getting event")
	}

	if event == nil {
		return nil, errors.New("event not found")
	}

	event.Title = req.Title
	event.Description = &req.Description
	event.Complete = *req.Complete
	event.Location = req.Location
	event.StartTime = req.StartTime
	event.EndTime = req.EndTime

	if err := u.eventRepository.UpdateEvent(event); err != nil {
		return nil, errors.Wrap(err, "[EventUsecase.UpdateEvent]: Error updating event")
	}

	return &response.EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Complete:    &event.Complete,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
	}, nil
}

func (u *eventUsecase) DeleteEvent(id uint64) error {
	if err := u.eventRepository.DeleteEvent(id); err != nil {
		return errors.Wrap(err, "[EventUsecase.DeleteEvent]: Error deleting event")
	}
	return nil
}
