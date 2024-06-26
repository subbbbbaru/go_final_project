package service

import (
	"time"

	"github.com/subbbbbaru/go_final_project/internal/models"
	"github.com/subbbbbaru/go_final_project/internal/repository"
)

//	type Auth interface {
//		CreateUser(user todo.User) (int, error)
//		GenerateToken(username, password string) (string, error)
//		ParseToken(token string) (int, error)
//	}
type TodoTask interface {
	NextDate(now time.Time, date string, repeat string) (string, error)
	Create( /*userId int64,*/ task models.Task) (int64, error)
	GetTasks( /*userId int64,*/ search string) ([]models.Task, error)
	GetTaskById( /*userId int64,*/ taskId int) (models.Task, error)
	Update( /*userId int64,*/ task models.Task) (models.Task, error)
	Delete( /*userId int64,*/ taskId int) (models.Task, error)
	Done( /*userId int64,*/ taskId int) (models.Task, error)
}

type Service struct {
	// Auth
	TodoTask
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		// Auth:     NewAuthService(repos.Auth),
		TodoTask: NewTodoTaskService(repos.TodoTask),
	}
}
