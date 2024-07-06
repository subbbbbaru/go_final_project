package repository

import (
	"database/sql"
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/subbbbbaru/go_final_project/internal/models"
	"github.com/subbbbbaru/go_final_project/utils"
)

const (
	taskDateLayout   = "20060102"
	searchDateLayout = "02.01.2006"
	limitTasks       = 45

	NonSearch   = 0 // нет значения для поиска
	SearchTitle = 1 // поиск по дате
	SearchDate  = 2 // поиск по названию
)

type TodoTaskSQLite struct {
	db *sql.DB
}

func NewTodoTaskSQLite(db *sql.DB) *TodoTaskSQLite {
	return &TodoTaskSQLite{db: db}
}

func (todo *TodoTaskSQLite) Create(task models.Task) (int64, error) {

	query := fmt.Sprintf("INSERT INTO %s (title, date, comment, repeat) VALUES (?, ?, ?, ?) RETURNING id", todoTaskTable)
	row := todo.db.QueryRow(query, task.Title, task.Date, task.Comment, task.Repeat)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (todo *TodoTaskSQLite) GetTasks(search string) ([]models.Task, error) {

	switch typeSearch(search) {
	case NonSearch:
		return todo.nonSearch()
	case SearchDate:
		return todo.dateSearch(search)
	case SearchTitle:
		return todo.titleSearch(search)
	}

	return []models.Task{}, nil
}

func (todo *TodoTaskSQLite) GetTaskById(taskId int) (models.Task, error) {
	var task models.Task

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", todoTaskTable)
	row := todo.db.QueryRow(query, taskId, limitTasks)

	if err := row.Scan(&task.ID, &task.Title, &task.Date, &task.Comment, &task.Repeat); err != nil {
		return task, err
	}

	return task, nil
}

func (todo *TodoTaskSQLite) Update(task models.Task) (models.Task, error) {

	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return models.Task{}, err
	}

	_, err = todo.GetTaskById(id)
	if err != nil {
		return models.Task{}, err
	}

	query := fmt.Sprintf("UPDATE %s SET title = ?, date = ?, Comment = ?, repeat = ? WHERE id = ?", todoTaskTable)

	_, err = todo.db.Exec(query, task.Title, task.Date, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return models.Task{}, err
	}

	return models.Task{}, nil
}

func (todo *TodoTaskSQLite) Delete(taskId int) (models.Task, error) {

	_, err := todo.GetTaskById(taskId)
	if err != nil {
		return models.Task{}, err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", todoTaskTable)

	_, err = todo.db.Exec(query, taskId)
	if err != nil {
		return models.Task{}, err
	}

	return models.Task{}, nil
}

func (todo *TodoTaskSQLite) Done(taskId int) (models.Task, error) {
	task, err := todo.GetTaskById(taskId)
	if err != nil {
		return models.Task{}, err
	}
	if strings.TrimSpace(task.Repeat) == "" {
		return todo.Delete(taskId)
	}
	now := time.Now() //.UTC().Truncate(24 * time.Hour)

	nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return models.Task{}, err
	}

	task.Date = nextDate

	return todo.Update(task)
}

func (todo *TodoTaskSQLite) nonSearch() ([]models.Task, error) {
	tasks := make([]models.Task, 0)
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY date LIMIT ?", todoTaskTable)
	rows, err := todo.db.Query(query, limitTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Date, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func (todo *TodoTaskSQLite) dateSearch(search string) ([]models.Task, error) {

	date, err := time.Parse(searchDateLayout, search)
	if err != nil {
		return nil, err
	}
	dateStr := date.Format(taskDateLayout)

	tasks := make([]models.Task, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE date = ? LIMIT ?", todoTaskTable)
	rows, err := todo.db.Query(query, dateStr, limitTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Date, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func (todo *TodoTaskSQLite) titleSearch(search string) ([]models.Task, error) {
	search = fmt.Sprintf(`%%%s%%`, search)
	tasks := make([]models.Task, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?", todoTaskTable)
	rows, err := todo.db.Query(query, search, search, limitTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Date, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func typeSearch(search string) int {
	if len(search) == 0 {
		return NonSearch
	}

	if _, err := time.Parse(searchDateLayout, search); err == nil {
		return SearchDate
	}
	return SearchTitle
}
