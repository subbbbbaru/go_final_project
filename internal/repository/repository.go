package repository

import (
	"database/sql"
	"time"

	"github.com/subbbbbaru/go_final_project/internal/models"
)

type Auth interface {
	//	CreateUser(user string) (int, error)
	GetPassword() (string, error)
}

type TodoTask interface {
	NextDate(now time.Time, date string, repeat string) (string, error)
	Create( /*userId int64,*/ task models.Task) (int64, error)
	GetTasks( /*userId int64,*/ search string) ([]models.Task, error)
	GetTaskById( /*userId int64,*/ taskId int) (models.Task, error)
	Update( /*userId,*/ task models.Task) (models.Task, error)
	Delete( /*userId int64,*/ taskId int) (models.Task, error)
	Done( /*userId int64,*/ taskId int) (models.Task, error)
}

type Repository struct {
	Auth
	TodoTask
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Auth:     NewAuthFromEnv(),
		TodoTask: NewTodoTaskSQLite(db), // NewTodoTaskSQLite(db),
	}
}
