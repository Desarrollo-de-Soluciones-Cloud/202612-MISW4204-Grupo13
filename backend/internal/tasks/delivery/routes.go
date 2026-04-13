package delivery

import (
	assignmentsInfrastructure "backend/internal/assignments/infrastructure"
	"backend/internal/tasks/application"
	"backend/internal/tasks/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter) {
	repo := infrastructure.NewTaskRepository()
	if err := repo.AutoMigrate(); err != nil {
		panic(err)
	}
	if err := repo.NormalizeLegacyStatuses(); err != nil {
		panic(err)
	}
	assignmentRepo := assignmentsInfrastructure.NewAssignmentRepository()

	createTask := application.NewCreateTask(repo, assignmentRepo, nil)
	listTasks := application.NewListTasks(repo)
	getTaskByID := application.NewGetTaskByID(repo)
	updateTask := application.NewUpdateTask(repo, assignmentRepo, nil)
	deleteTask := application.NewDeleteTask(repo, nil)
	handler := NewTaskHandler(createTask, listTasks, getTaskByID, updateTask, deleteTask)

	tasks := r.Group("/tasks")
	{
		tasks.POST("", handler.CreateTask)
		tasks.GET("", handler.ListTasks)
		tasks.GET("/:id", handler.GetTaskByID)
		tasks.PUT("/:id", handler.UpdateTask)
		tasks.DELETE("/:id", handler.DeleteTask)
	}
}
