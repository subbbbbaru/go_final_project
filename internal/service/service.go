package service

import (
	"github.com/subbbbbaru/go_final_project/internal/models"
	"github.com/subbbbbaru/go_final_project/internal/repository"
)

type Auth interface {
	GenerateToken(password string) (string, error)
	ValideToken(token string) (bool, error)
}
type TodoTask interface {
	Create(task models.Task) (int64, error)
	GetTasks(search string) ([]models.Task, error)
	GetTaskById(taskId int) (models.Task, error)
	Update(task models.Task) (models.Task, error)
	Delete(taskId int) error
	Done(taskId int) error
}

type Service struct {
	Auth
	TodoTask
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth:     NewAuthService(repos.Auth),
		TodoTask: NewTodoTaskService(repos.TodoTask),
	}
}
