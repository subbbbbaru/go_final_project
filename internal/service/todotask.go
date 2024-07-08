package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/subbbbbaru/first-sample/pkg/log"

	"github.com/subbbbbaru/go_final_project/internal/models"
	"github.com/subbbbbaru/go_final_project/internal/nextdate"
	"github.com/subbbbbaru/go_final_project/internal/repository"
)

const taskDateLayout = "20060102"

type TodoItemService struct {
	repo repository.TodoTask
}

func NewTodoTaskService(repo repository.TodoTask) *TodoItemService {
	return &TodoItemService{repo: repo}
}

func (todoService *TodoItemService) Create(task models.Task) (int64, error) {
	if task.Title == "" {
		log.Error().Println("task title not found")
		return 0, errors.New("task title not found")
	}

	timeNow := time.Now().Truncate(24 * time.Hour).UTC()

	if task.Date == "" {
		task.Date = timeNow.Format(taskDateLayout)
	}

	date, err := time.Parse(taskDateLayout, task.Date)
	if err != nil {
		log.Error().Println(err)
		return 0, err
	}

	if date.Before(timeNow) {
		if task.Repeat == "" {
			task.Date = timeNow.Format(taskDateLayout)
		} else {
			nextDate, err := nextdate.NextDate(timeNow, task.Date, task.Repeat)
			if err != nil {
				log.Error().Println(err)
				return 0, err
			}
			task.Date = nextDate
		}
	}

	return todoService.repo.Create(task)
}

func (todoService *TodoItemService) GetTasks(search string) ([]models.Task, error) {
	return todoService.repo.GetTasks(search)
}

func (todoService *TodoItemService) GetTaskById(taskId int) (models.Task, error) {
	return todoService.repo.GetTaskById(taskId)
}

func (todoService *TodoItemService) Update(task models.Task) (models.Task, error) {
	if strings.TrimSpace(task.ID) == "" {
		log.Error().Println("task ID not found")
		return models.Task{}, errors.New("task title not found")
	}

	if strings.TrimSpace(task.Title) == "" {
		log.Error().Println("task title not found")
		return models.Task{}, errors.New("task title not found")
	}

	timeNow := time.Now().Truncate(24 * time.Hour).UTC()

	if strings.TrimSpace(task.Date) == "" {
		task.Date = timeNow.Format(taskDateLayout)
	}

	date, err := time.Parse(taskDateLayout, task.Date)
	if err != nil {
		log.Error().Println(err)
		return models.Task{}, err
	}

	if date.Before(timeNow) {
		if task.Repeat == "" {
			task.Date = timeNow.Format(taskDateLayout)
		} else {
			task.Date, err = nextdate.NextDate(timeNow, task.Date, task.Repeat)
			if err != nil {
				log.Error().Println(err)
				return models.Task{}, err
			}
		}
	}
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return models.Task{}, err
	}

	if _, errId := todoService.GetTaskById(id); errId != nil {
		return models.Task{}, errId
	}

	return todoService.repo.Update(task)
}

func (todoService *TodoItemService) Delete(taskId int) error {
	if _, errId := todoService.GetTaskById(taskId); errId != nil {
		return errId
	}
	_, err := todoService.repo.Delete(taskId)
	return err
}

func (todoService *TodoItemService) Done(taskId int) error {
	task, err := todoService.GetTaskById(taskId)
	if err != nil {
		return err
	}
	if strings.TrimSpace(task.Repeat) == "" {
		return todoService.Delete(taskId)
	}
	now := time.Now().UTC().Truncate(24 * time.Hour)
	nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}
	task.Date = nextDate
	_, err = todoService.Update(task)
	return err
}
