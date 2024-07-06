package repository

import (
	"database/sql"

	"github.com/subbbbbaru/go_final_project/internal/models"
)

type Auth interface {
	GetPassword() (string, error)
}

type TodoTask interface {
	Create(task models.Task) (int64, error)
	GetTasks(search string) ([]models.Task, error)
	GetTaskById(taskId int) (models.Task, error)
	Update(task models.Task) (models.Task, error)
	Delete(taskId int) (models.Task, error)
	Done(taskId int) (models.Task, error)
}

type Repository struct {
	Auth
	TodoTask
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Auth:     NewAuthFromEnv(),
		TodoTask: NewTodoTaskSQLite(db),
	}
}
