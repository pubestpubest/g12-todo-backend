package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/g12-todo-backend/database"
	"github.com/pubestpubest/g12-todo-backend/feature/task/delivery"
	"github.com/pubestpubest/g12-todo-backend/feature/task/repository"
	"github.com/pubestpubest/g12-todo-backend/feature/task/usecase"
)

func TaskRoutes(router *gin.RouterGroup) {
	// NewTaskRepository := repository.NewTaskRepository(database.DB)
	// newTaskUsecase := usecase.NewTaskUsecase(NewTaskRepository)
	taskHandler := delivery.NewTaskHandler(
		usecase.NewTaskUsecase(
			repository.NewTaskRepository(database.DB)))

	taskRoutes := router.Group("/v1/tasks")
	{
		taskRoutes.GET("", taskHandler.GetTaskList)
		taskRoutes.GET("/:id", taskHandler.GetTaskByID)
		taskRoutes.POST("", taskHandler.CreateTask)
		taskRoutes.PUT("/:id", taskHandler.UpdateTask)
		taskRoutes.DELETE("/:id", taskHandler.DeleteTask)
	}
}
