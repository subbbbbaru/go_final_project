package service

import (
	"github.com/subbbbbaru/go_final_project/internal/models"
	"github.com/subbbbbaru/go_final_project/internal/repository"
)

type TodoItemService struct {
	repo repository.TodoTask
}

func NewTodoTaskService(repo repository.TodoTask) *TodoItemService {
	return &TodoItemService{repo: repo}
}

func (todoService *TodoItemService) Create(task models.Task) (int64, error) {
	return todoService.repo.Create(task)
}

func (todoService *TodoItemService) GetTasks(search string) ([]models.Task, error) {
	return todoService.repo.GetTasks(search)
}

func (todoService *TodoItemService) GetTaskById(taskId int) (models.Task, error) {
	return todoService.repo.GetTaskById(taskId)
}

func (todoService *TodoItemService) Update(task models.Task) (models.Task, error) {
	return todoService.repo.Update(task)
}

func (todoService *TodoItemService) Delete(taskId int) (models.Task, error) {
	return todoService.repo.Delete(taskId)
}

func (todoService *TodoItemService) Done(taskId int) (models.Task, error) {
	return todoService.repo.Done(taskId)
}
