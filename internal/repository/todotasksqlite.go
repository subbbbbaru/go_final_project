package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/subbbbbaru/go_final_project/internal/models"
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

func (todo *TodoTaskSQLite) NextDate(now time.Time, date string, repeat string) (string, error) {
	taskdate, err := time.Parse(taskDateLayout, date)
	if err != nil {
		return "", err
	}
	return nextDate(now, taskdate, repeat)
}

func (todo *TodoTaskSQLite) Create( /*userId int64,*/ task models.Task) (int64, error) {

	query := fmt.Sprintf("INSERT INTO %s (title, date, comment, repeat) VALUES (?, ?, ?, ?) RETURNING id", todoTaskTable)
	row := todo.db.QueryRow(query, task.Title, task.Date, task.Comment, task.Repeat)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (todo *TodoTaskSQLite) GetTasks( /*userId int64,*/ search string) ([]models.Task, error) {

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

func (todo *TodoTaskSQLite) GetTaskById( /*userId int64,*/ taskId int) (models.Task, error) {
	var task models.Task

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", todoTaskTable)
	row := todo.db.QueryRow(query, taskId, limitTasks)

	if err := row.Scan(&task.ID, &task.Title, &task.Date, &task.Comment, &task.Repeat); err != nil {
		return task, err
	}

	return task, nil
}

func (todo *TodoTaskSQLite) Update( /*userId int64,*/ task models.Task) (models.Task, error) {

	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return models.Task{}, err
	}
	log.Println("Update this ID = ", id)

	_, err = todo.GetTaskById(id)
	if err != nil {
		return models.Task{}, err
	}

	query := fmt.Sprintf("UPDATE %s SET title = ?, date = ?, Comment = ?, repeat = ? WHERE id = ?", todoTaskTable)

	_, err = todo.db.Exec(query, task.Title, task.Date, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return models.Task{}, err
	}

	log.Println("Updated task = ", task)

	return models.Task{}, nil
}

func (todo *TodoTaskSQLite) Delete( /*userId int64,*/ taskId int) (models.Task, error) {

	_, err := todo.GetTaskById(taskId)
	if err != nil {
		return models.Task{}, err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", todoTaskTable)

	_, err = todo.db.Exec(query, taskId)
	if err != nil {
		return models.Task{}, err
	}

	log.Println("Delete task with ID = ", taskId)
	return models.Task{}, nil
}

func (todo *TodoTaskSQLite) Done( /*userId int64,*/ taskId int) (models.Task, error) {
	task, err := todo.GetTaskById(taskId)
	if err != nil {
		return models.Task{}, err
	}
	if strings.TrimSpace(task.Repeat) == "" {
		return todo.Delete(taskId)
	}
	now := time.Now().UTC().Truncate(24 * time.Hour)

	nextDate, err := todo.NextDate(now, task.Date, task.Repeat)
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

func nextDate(now time.Time, taskDate time.Time, repeat string) (string, error) {
	var result string
	if len(repeat) == 0 {
		return "", errors.New("invalid repeat")
	}

	nextDayRules := strings.Fields(repeat)

	switch nextDayRules[0] {
	case "d":
		return getNextDay(now, taskDate, nextDayRules...)
	case "w":
		return getNextWeek(now, nextDayRules...)
	case "m":
		return getMonth(now, taskDate, nextDayRules...)
		// return "", errors.New("month not realized") // getNextMonth(now, taskDate, nextDayRules...)
	case "y":
		result = getNextYear(now, taskDate)
		return result, nil
	}

	return "", errors.New("invalid repeat")
}

func getNextDay(now time.Time, taskDate time.Time, day ...string) (string, error) {
	if len(day) != 2 {
		return "", errors.New("invalid day")
	}
	var result string
	days, err := strconv.Atoi(day[1])

	if err != nil {
		return result, err
	}
	if days > 0 && days < 401 {
		taskDate = taskDate.AddDate(0, 0, days)
		for taskDate.Before(now) || taskDate.Equal(now) {
			taskDate = taskDate.AddDate(0, 0, days)
		}
		result = taskDate.Format(taskDateLayout)
		return result, nil
	}
	return result, errors.New("wrong next day")
}

func getNextWeek(now time.Time, daysOnWeek ...string) (string, error) {
	if len(daysOnWeek) != 2 {
		return "", errors.New("wrong next week")
	}

	days := strings.Split(daysOnWeek[1], ",")
	if len(days) > 8 {
		return "", errors.New("wrong parse days on week")
	}

	weekDays := make([]int, len(days))
	for idx, day := range days {
		dayInt, err := strconv.Atoi(day)
		if err != nil {
			return "", err
		}
		if dayInt < 0 || dayInt > 7 {
			return "", fmt.Errorf("day out of week day: day is %d", dayInt)
		}
		weekDays[idx] = dayInt
	}
	sort.Ints(weekDays)

	currDay := int(now.Weekday())
	if currDay == 0 {
		currDay = 7
	}

	for _, day := range weekDays {
		if day > currDay {
			nextDay := now.AddDate(0, 0, day-currDay)
			return nextDay.Format(taskDateLayout), nil
		}
	}

	nextDay := now.AddDate(0, 0, 7+weekDays[0]-currDay)

	return nextDay.Format(taskDateLayout), nil
}

func getMonth(now time.Time, date time.Time, parts ...string) (string, error) {

	// Парсинг дней
	daysStr := parts[1]
	dayStrs := strings.Split(daysStr, ",")
	days := []int{}
	for _, ds := range dayStrs {
		day, err := strconv.Atoi(ds)
		if err != nil {
			return "", fmt.Errorf("invalid day in repeat rule: %v", err)
		}
		if day == -1 || day == -2 || (day >= 1 && day <= 31) {
			days = append(days, day)
		} else {
			return "", fmt.Errorf("day out of month: day is %d", day)
		}
	}
	sort.Ints(days)

	// Парсинг месяцев, если указаны
	months := map[int]bool{}
	if len(parts) > 2 {
		monthsStr := parts[2]
		monthStrs := strings.Split(monthsStr, ",")
		for _, ms := range monthStrs {
			month, err := strconv.Atoi(ms)
			if err != nil {
				return "", fmt.Errorf("invalid month in repeat rule: %v", err)
			}
			months[month] = true
		}
	}

	// Функция для получения следующей даты
	getNextDate := func(year int, month time.Month) (time.Time, error) {
		// for _, day := range days {
		// 	var targetDate time.Time
		// 	if day > 0 {
		// 		targetDate = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		// 	} else if day == -1 {
		// 		targetDate = time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
		// 	} else if day == -2 {
		// 		targetDate = time.Date(year, month+1, -1, 0, 0, 0, 0, time.UTC)
		// 	} else {
		// 		return time.Time{}, fmt.Errorf("invalid day: %d", day)
		// 	}

		// 	if targetDate.After(now) {
		// 		return targetDate, nil
		// 	}
		// }
		daysInCurrentMonth := daysInMonth(year, month)
		var possibleDates []time.Time
		for _, day := range days {
			var targetDate time.Time
			if day > 0 {
				if day > daysInCurrentMonth {
					continue
				}
				targetDate = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
			} else if day == -1 {
				targetDate = time.Date(year, month, daysInCurrentMonth, 0, 0, 0, 0, time.UTC)
			} else if day == -2 {
				targetDate = time.Date(year, month, daysInCurrentMonth-1, 0, 0, 0, 0, time.UTC)
			} else {
				return time.Time{}, fmt.Errorf("invalid day: %d", day)
			}
			possibleDates = append(possibleDates, targetDate)
		}
		sort.Slice(possibleDates, func(i, j int) bool {
			return possibleDates[i].Before(possibleDates[j])
		})

		for _, d := range possibleDates {
			if d.After(now) {
				return d, nil
			}
		}

		return time.Time{}, fmt.Errorf("no valid next date found")
	}

	// Поиск следующей даты
	for {
		for m := date.Month(); m <= 12; m++ {
			if len(months) == 0 || months[int(m)] {
				nextDate, err := getNextDate(date.Year(), m)
				if err == nil {
					return nextDate.Format("20060102"), nil
				}
			}
		}
		date = time.Date(date.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	}
}

func getMonth1(now time.Time, taskDate time.Time, daysOnWeek ...string) (string, error) {

	days := []int{}
	for _, day := range strings.Split(daysOnWeek[1], ",") {
		if day == "-1" {
			days = append(days, -1)
		} else if day == "-2" {
			days = append(days, -2)
		} else {
			dayInt, err := strconv.Atoi(day)
			if err != nil || dayInt < 1 || dayInt > 31 {
				return "", errors.New("неверный формат дня в правиле повторения: " + day)
			}
			days = append(days, dayInt)
		}
	}

	months := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	if len(daysOnWeek) > 2 {
		months = []int{}
		for _, month := range strings.Split(daysOnWeek[2], ",") {
			monthInt, err := strconv.Atoi(month)
			if err != nil || monthInt < 1 || monthInt > 12 {
				return "", errors.New("неверный формат месяца в правиле повторения: " + month)
			}
			months = append(months, monthInt)
		}
	}

	// Поиск следующей даты
	nextDate := taskDate
	for nextDate.Before(now) {
		nextDate = nextDate.AddDate(0, 0, 1)
	}

	for _, day := range days {
		if day == -1 {
			if nextDate.Day() == daysInMonth(nextDate.Year(), nextDate.Month()) {
				nextDate = nextDate.AddDate(0, 0, 1)
			}
		} else if day == -2 {
			if nextDate.Day() == daysInMonth(nextDate.Year(), nextDate.Month())-1 {
				nextDate = nextDate.AddDate(0, 0, 1)
			}
		} else if nextDate.Day() == day {
			for _, month := range months {
				if nextDate.Month() == time.Month(month) {

					return nextDate.Format("20060102"), nil
				}
			}
		}
	}

	return "", errors.New("не найдена следующая дата")
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// func getNextMonth(now time.Time, taskDate time.Time, daysOnMonth ...string) (string, error) {
// 	if len(daysOnMonth) != 2 || len(daysOnMonth) != 3 {
// 		return "", errors.New("wrong next week")
// 	}
// 	if len(daysOnMonth) == 3 {
// 		days := strings.Split(daysOnMonth[1], ",")
// 		if len(days) > 8 {
// 			return "", errors.New("wrong parse days on week")
// 		}
// 		return "", nil
// 	}

// 	if len(daysOnMonth) == 2 {

// 	}
// 	return 0, errors.New("wrong next day")
// }

// func calcMonthWithDay(now time.Time, taskDate time.Time, days string) (string, error) {
// 	daysList := strings.Split(days, ",")

// 	monthDays := make([]int, len(daysList))
// 	for idx, day := range days {
// 		dayInt, err := strconv.Atoi(daysList)
// 		if err != nil {
// 			return "", err
// 		}
// 		monthDays[idx] = dayInt
// 	}
// 	sort.Ints(monthDays)

// 	if number > 0 && number < 32 {
// 		return number, nil
// 	}

// 	if number == -1 || number == -2 {

// 	}

// }

// func calcMonthWithDayAndMonth(now time.Time, taskDate time.Time, days, months string) (string, error) {

// }

func getNextYear(now time.Time, taskDate time.Time) string {
	taskDate = taskDate.AddDate(1, 0, 0)
	for taskDate.Before(now) || taskDate.Equal(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	}
	return taskDate.Format(taskDateLayout)
}
