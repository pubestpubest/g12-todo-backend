package usecase

import (
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/pubestpubest/g12-todo-backend/models"
	"github.com/pubestpubest/g12-todo-backend/request"
)

// mockEventRepository implements domain.EventRepository for testing
type mockEventRepository struct {
	events        []*models.Events
	shouldError   bool
	errorMessage  string
	getByIDResult *models.Events
	getByIDError  error
}

func newMockEventRepository() *mockEventRepository {
	return &mockEventRepository{
		events:      make([]*models.Events, 0),
		shouldError: false,
	}
}

func (m *mockEventRepository) GetEventList(page, limit int) ([]*models.Events, int64, error) {
	if m.shouldError {
		return nil, 0, errors.New(m.errorMessage)
	}

	total := int64(len(m.events))

	// Calculate pagination
	start := (page - 1) * limit
	end := start + limit

	if start >= len(m.events) {
		return []*models.Events{}, total, nil
	}

	if end > len(m.events) {
		end = len(m.events)
	}

	return m.events[start:end], total, nil
}

func (m *mockEventRepository) GetEventByID(id uint64) (*models.Events, error) {
	if m.getByIDError != nil {
		return nil, m.getByIDError
	}

	if m.getByIDResult != nil {
		return m.getByIDResult, nil
	}

	for _, event := range m.events {
		if event.ID == id {
			return event, nil
		}
	}
	return nil, nil
}

func (m *mockEventRepository) CreateEvent(event *models.Events) error {
	if m.shouldError {
		return errors.New(m.errorMessage)
	}

	// Simulate auto-increment ID
	event.ID = uint64(len(m.events) + 1)
	now := time.Now()
	event.CreatedAt = &now
	event.UpdatedAt = &now

	m.events = append(m.events, event)
	return nil
}

func (m *mockEventRepository) UpdateEvent(event *models.Events) error {
	if m.shouldError {
		return errors.New(m.errorMessage)
	}

	for i, e := range m.events {
		if e.ID == event.ID {
			now := time.Now()
			event.UpdatedAt = &now
			m.events[i] = event
			return nil
		}
	}
	return errors.New("event not found")
}

func (m *mockEventRepository) DeleteEvent(id uint64) error {
	if m.shouldError {
		return errors.New(m.errorMessage)
	}

	for i, event := range m.events {
		if event.ID == id {
			m.events = append(m.events[:i], m.events[i+1:]...)
			return nil
		}
	}
	return errors.New("event not found")
}

// Helper function to create test events
func createTestEvent(id uint64, title, description, location string, complete bool, startTime, endTime time.Time) *models.Events {
	now := time.Now()
	return &models.Events{
		ID:          id,
		Title:       title,
		Description: &description,
		Complete:    complete,
		CreatedAt:   &now,
		UpdatedAt:   &now,
		Location:    location,
		StartTime:   startTime,
		EndTime:     endTime,
	}
}

// Helper function to create test request
func createTestEventRequest(title, description, location string, complete bool, startTime, endTime time.Time) *request.EventRequest {
	return &request.EventRequest{
		Title:       title,
		Description: description,
		Location:    location,
		StartTime:   startTime,
		EndTime:     endTime,
		Complete:    &complete,
	}
}

// Helper function to get deterministic test times
func getTestTimes() (time.Time, time.Time) {
	baseTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	return baseTime, baseTime.Add(time.Hour)
}

func TestNewEventUsecase(t *testing.T) {
	mockRepo := newMockEventRepository()
	usecase := NewEventUsecase(mockRepo)

	if usecase == nil {
		t.Error("Expected usecase to be created, got nil")
	}
}

func TestEventUsecase_GetEventList(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		mockEvents    []*models.Events
		shouldError   bool
		errorMessage  string
		expectedError bool
		expectedTotal int
		expectedPages int
		expectedData  int
	}{
		{
			name:  "successful get event list",
			page:  1,
			limit: 2,
			mockEvents: []*models.Events{
				createTestEvent(1, "Event 1", "Description 1", "Location 1", false, time.Now(), time.Now().Add(time.Hour)),
				createTestEvent(2, "Event 2", "Description 2", "Location 2", true, time.Now(), time.Now().Add(time.Hour)),
				createTestEvent(3, "Event 3", "Description 3", "Location 3", false, time.Now(), time.Now().Add(time.Hour)),
			},
			expectedError: false,
			expectedTotal: 3,
			expectedPages: 2,
			expectedData:  2,
		},
		{
			name:          "repository error",
			page:          1,
			limit:         10,
			shouldError:   true,
			errorMessage:  "database error",
			expectedError: true,
		},
		{
			name:          "empty result",
			page:          1,
			limit:         10,
			mockEvents:    []*models.Events{},
			expectedError: false,
			expectedTotal: 0,
			expectedPages: 0,
			expectedData:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockEventRepository()
			mockRepo.events = tt.mockEvents
			mockRepo.shouldError = tt.shouldError
			mockRepo.errorMessage = tt.errorMessage

			usecase := NewEventUsecase(mockRepo)
			result, err := usecase.GetEventList(tt.page, tt.limit)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result, got nil")
				return
			}

			if result.Pagination.Total != tt.expectedTotal {
				t.Errorf("Expected total %d, got %d", tt.expectedTotal, result.Pagination.Total)
			}

			if result.Pagination.TotalPages != tt.expectedPages {
				t.Errorf("Expected total pages %d, got %d", tt.expectedPages, result.Pagination.TotalPages)
			}

			if len(result.Data) != tt.expectedData {
				t.Errorf("Expected data length %d, got %d", tt.expectedData, len(result.Data))
			}
		})
	}
}

func TestEventUsecase_GetEventByID(t *testing.T) {
	tests := []struct {
		name          string
		eventID       uint64
		mockEvent     *models.Events
		mockError     error
		expectedError bool
		expectedNil   bool
	}{
		{
			name:    "successful get event by id",
			eventID: 1,
			mockEvent: createTestEvent(1, "Test Event", "Test Description", "Test Location", false,
				time.Now(), time.Now().Add(time.Hour)),
			expectedError: false,
			expectedNil:   false,
		},
		{
			name:          "repository error",
			eventID:       1,
			mockError:     errors.New("database error"),
			expectedError: true,
		},
		{
			name:          "event not found",
			eventID:       999,
			mockEvent:     nil,
			expectedError: true,
			expectedNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockEventRepository()
			mockRepo.getByIDResult = tt.mockEvent
			mockRepo.getByIDError = tt.mockError

			usecase := NewEventUsecase(mockRepo)
			result, err := usecase.GetEventByID(tt.eventID)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result, got nil")
				return
			}

			if result.ID != tt.mockEvent.ID {
				t.Errorf("Expected ID %d, got %d", tt.mockEvent.ID, result.ID)
			}

			if result.Title != tt.mockEvent.Title {
				t.Errorf("Expected title %s, got %s", tt.mockEvent.Title, result.Title)
			}
		})
	}
}

func TestEventUsecase_CreateEvent(t *testing.T) {
	tests := []struct {
		name           string
		request        *request.EventRequest
		shouldError    bool
		errorMessage   string
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "successful create event",
			request: createTestEventRequest("Test Event", "Test Description", "Test Location", false,
				time.Now(), time.Now().Add(time.Hour)),
			expectedError: false,
		},
		{
			name: "start time after end time",
			request: func() *request.EventRequest {
				startTime, endTime := getTestTimes()
				return createTestEventRequest("Test Event", "Test Description", "Test Location", false,
					endTime, startTime) // Swap them to make start after end
			}(),
			expectedError:  true,
			expectedErrMsg: "startTime must be before endTime",
		},
		{
			name: "start time equal to end time",
			request: func() *request.EventRequest {
				baseTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
				return createTestEventRequest("Test Event", "Test Description", "Test Location", false,
					baseTime, baseTime) // Same exact time
			}(),
			expectedError:  true,
			expectedErrMsg: "startTime must be before endTime",
		},
		{
			name: "repository error",
			request: createTestEventRequest("Test Event", "Test Description", "Test Location", false,
				time.Now(), time.Now().Add(time.Hour)),
			shouldError:   true,
			errorMessage:  "database error",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockEventRepository()
			mockRepo.shouldError = tt.shouldError
			mockRepo.errorMessage = tt.errorMessage

			usecase := NewEventUsecase(mockRepo)
			result, err := usecase.CreateEvent(tt.request)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				if tt.expectedErrMsg != "" && !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("Expected error message to contain '%s', got: %v", tt.expectedErrMsg, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result, got nil")
				return
			}

			if result.Title != tt.request.Title {
				t.Errorf("Expected title %s, got %s", tt.request.Title, result.Title)
			}

			if result.Location != tt.request.Location {
				t.Errorf("Expected location %s, got %s", tt.request.Location, result.Location)
			}
		})
	}
}

func TestEventUsecase_UpdateEvent(t *testing.T) {
	existingEvent := createTestEvent(1, "Original Event", "Original Description", "Original Location", false,
		time.Now(), time.Now().Add(time.Hour))

	tests := []struct {
		name           string
		eventID        uint64
		request        *request.EventRequest
		mockEvent      *models.Events
		mockGetError   error
		shouldError    bool
		errorMessage   string
		expectedError  bool
		expectedErrMsg string
		setupEvents    []*models.Events
	}{
		{
			name:    "successful update event",
			eventID: 1,
			request: createTestEventRequest("Updated Event", "Updated Description", "Updated Location", true,
				time.Now(), time.Now().Add(time.Hour)),
			mockEvent:     existingEvent,
			expectedError: false,
			setupEvents:   []*models.Events{existingEvent},
		},
		{
			name:    "start time after end time",
			eventID: 1,
			request: func() *request.EventRequest {
				startTime, endTime := getTestTimes()
				return createTestEventRequest("Updated Event", "Updated Description", "Updated Location", true,
					endTime, startTime) // Swap them to make start after end
			}(),
			expectedError:  true,
			expectedErrMsg: "startTime must be before endTime",
		},
		{
			name:    "event not found",
			eventID: 999,
			request: createTestEventRequest("Updated Event", "Updated Description", "Updated Location", true,
				time.Now(), time.Now().Add(time.Hour)),
			mockEvent:      nil,
			expectedError:  true,
			expectedErrMsg: "event not found",
		},
		{
			name:    "repository get error",
			eventID: 1,
			request: createTestEventRequest("Updated Event", "Updated Description", "Updated Location", true,
				time.Now(), time.Now().Add(time.Hour)),
			mockGetError:  errors.New("database error"),
			expectedError: true,
		},
		{
			name:    "repository update error",
			eventID: 1,
			request: createTestEventRequest("Updated Event", "Updated Description", "Updated Location", true,
				time.Now(), time.Now().Add(time.Hour)),
			mockEvent:     existingEvent,
			shouldError:   true,
			errorMessage:  "update failed",
			expectedError: true,
			setupEvents:   []*models.Events{existingEvent},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockEventRepository()
			mockRepo.getByIDResult = tt.mockEvent
			mockRepo.getByIDError = tt.mockGetError
			mockRepo.shouldError = tt.shouldError
			mockRepo.errorMessage = tt.errorMessage
			if tt.setupEvents != nil {
				mockRepo.events = tt.setupEvents
			}

			usecase := NewEventUsecase(mockRepo)
			result, err := usecase.UpdateEvent(tt.eventID, tt.request)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				if tt.expectedErrMsg != "" && !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("Expected error message to contain '%s', got: %v", tt.expectedErrMsg, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result, got nil")
				return
			}

			if result.Title != tt.request.Title {
				t.Errorf("Expected title %s, got %s", tt.request.Title, result.Title)
			}
		})
	}
}

func TestEventUsecase_DeleteEvent(t *testing.T) {
	existingEvent := createTestEvent(1, "Test Event", "Test Description", "Test Location", false,
		time.Now(), time.Now().Add(time.Hour))

	tests := []struct {
		name          string
		eventID       uint64
		shouldError   bool
		errorMessage  string
		expectedError bool
		setupEvents   []*models.Events
	}{
		{
			name:          "successful delete event",
			eventID:       1,
			expectedError: false,
			setupEvents:   []*models.Events{existingEvent},
		},
		{
			name:          "repository error",
			eventID:       1,
			shouldError:   true,
			errorMessage:  "delete failed",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockEventRepository()
			mockRepo.shouldError = tt.shouldError
			mockRepo.errorMessage = tt.errorMessage
			if tt.setupEvents != nil {
				mockRepo.events = tt.setupEvents
			}

			usecase := NewEventUsecase(mockRepo)
			err := usecase.DeleteEvent(tt.eventID)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
