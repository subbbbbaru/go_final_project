package utils

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const taskDateLayout = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	taskdate, err := time.Parse(taskDateLayout, date)
	if err != nil {
		return "", err
	}
	return nextDate(now, taskdate, repeat)
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
					return nextDate.Format(taskDateLayout), nil
				}
			}
		}
		date = time.Date(date.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
	}
}

// func getMonth1(now time.Time, taskDate time.Time, daysOnWeek ...string) (string, error) {

// 	days := []int{}
// 	for _, day := range strings.Split(daysOnWeek[1], ",") {
// 		if day == "-1" {
// 			days = append(days, -1)
// 		} else if day == "-2" {
// 			days = append(days, -2)
// 		} else {
// 			dayInt, err := strconv.Atoi(day)
// 			if err != nil || dayInt < 1 || dayInt > 31 {
// 				return "", errors.New("неверный формат дня в правиле повторения: " + day)
// 			}
// 			days = append(days, dayInt)
// 		}
// 	}

// 	months := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
// 	if len(daysOnWeek) > 2 {
// 		months = []int{}
// 		for _, month := range strings.Split(daysOnWeek[2], ",") {
// 			monthInt, err := strconv.Atoi(month)
// 			if err != nil || monthInt < 1 || monthInt > 12 {
// 				return "", errors.New("неверный формат месяца в правиле повторения: " + month)
// 			}
// 			months = append(months, monthInt)
// 		}
// 	}

// 	// Поиск следующей даты
// 	nextDate := taskDate
// 	for nextDate.Before(now) {
// 		nextDate = nextDate.AddDate(0, 0, 1)
// 	}

// 	for _, day := range days {
// 		if day == -1 {
// 			if nextDate.Day() == daysInMonth(nextDate.Year(), nextDate.Month()) {
// 				nextDate = nextDate.AddDate(0, 0, 1)
// 			}
// 		} else if day == -2 {
// 			if nextDate.Day() == daysInMonth(nextDate.Year(), nextDate.Month())-1 {
// 				nextDate = nextDate.AddDate(0, 0, 1)
// 			}
// 		} else if nextDate.Day() == day {
// 			for _, month := range months {
// 				if nextDate.Month() == time.Month(month) {

// 					return nextDate.Format("20060102"), nil
// 				}
// 			}
// 		}
// 	}

// 	return "", errors.New("не найдена следующая дата")
// }

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func getNextYear(now time.Time, taskDate time.Time) string {
	taskDate = taskDate.AddDate(1, 0, 0)
	for taskDate.Before(now) || taskDate.Equal(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	}
	return taskDate.Format(taskDateLayout)
}
