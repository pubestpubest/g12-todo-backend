package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/g12-todo-backend/database"
	"github.com/pubestpubest/g12-todo-backend/feature/event/delivery"
	"github.com/pubestpubest/g12-todo-backend/feature/event/repository"
	"github.com/pubestpubest/g12-todo-backend/feature/event/usecase"
)

func EventRoutes(router *gin.RouterGroup) {
	// NewEventRepository := repository.NewEventRepository(database.DB)
	// newEventUsecase := usecase.NewEventUsecase(NewEventRepository)
	eventHandler := delivery.NewEventHandler(
		usecase.NewEventUsecase(
			repository.NewEventRepository(database.DB)))

	eventRoutes := router.Group("/events")
	{
		eventRoutes.GET("", eventHandler.GetEventList)
		eventRoutes.GET("/:id", eventHandler.GetEventByID)
		eventRoutes.POST("", eventHandler.CreateEvent)
		eventRoutes.PUT("/:id", eventHandler.UpdateEvent)
		eventRoutes.DELETE("/:id", eventHandler.DeleteEvent)
	}
}
