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

type eventHandler struct {
	eventUsecase domain.EventUsecase
}

func NewEventHandler(eventUsecase domain.EventUsecase) *eventHandler {
	return &eventHandler{eventUsecase: eventUsecase}
}

func (h *eventHandler) GetEventList(c *gin.Context) {
	var paginationReq request.PaginationRequest

	// Set default values
	paginationReq.Page = 1
	paginationReq.Limit = 10

	// Bind query parameters
	if err := c.ShouldBindQuery(&paginationReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	events, err := h.eventUsecase.GetEventList(paginationReq.Page, paginationReq.Limit)
	if err != nil {
		err = errors.Wrap(err, "[EventHandler.GetEventList]: Error getting event list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}

	resp := response.PaginatedResponse[*response.EventResponse]{
		Status:     constant.Success,
		Message:    "List events successfully",
		Data:       events.Data,
		Pagination: events.Pagination,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *eventHandler) GetEventByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := h.eventUsecase.GetEventByID(id)
	if err != nil {
		err = errors.Wrap(err, "[EventHandler.GetEventByID]: Error getting event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}

	resp := response.Response[*response.EventResponse]{
		Status:  constant.Success,
		Message: "Event retrieved successfully",
		Data:    event,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *eventHandler) CreateEvent(c *gin.Context) {
	var req request.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	event, err := h.eventUsecase.CreateEvent(&req)
	if err != nil {
		err = errors.Wrap(err, "[EventHandler.CreateEvent]: Error creating event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}

	resp := response.Response[*response.EventResponse]{
		Status:  constant.Success,
		Message: "Event created successfully",
		Data:    event,
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *eventHandler) UpdateEvent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var req request.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	event, err := h.eventUsecase.UpdateEvent(id, &req)
	if err != nil {
		err = errors.Wrap(err, "[EventHandler.UpdateEvent]: Error updating event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}

	resp := response.Response[*response.EventResponse]{
		Status:  constant.Success,
		Message: "Event updated successfully",
		Data:    event,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *eventHandler) DeleteEvent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := h.eventUsecase.DeleteEvent(id); err != nil {
		err = errors.Wrap(err, "[EventHandler.DeleteEvent]: Error deleting event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}

	resp := response.Response[interface{}]{
		Status:  constant.Success,
		Message: "Event deleted successfully",
		Data:    nil,
	}
	c.JSON(http.StatusOK, resp)
}
