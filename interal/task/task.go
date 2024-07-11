package task

import (
	"fmt"
	"go-final-project/repeat"
	"strconv"
	"time"
)

const layout = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (t *Task) CheckID() error {
	if t.ID == "" {
		return fmt.Errorf(`{"error":"Не указан индификатор задачи"}`)
	}
	_, err := strconv.ParseInt(t.ID, 10, 32)
	if err != nil {
		return fmt.Errorf(`{"error":"Указан невозможный индификатор задачи"}`)
	}
	return nil
}

func (t *Task) CheckTitle() error {
	if t.Title == "" {
		return fmt.Errorf(`{"error":"Не указан заголовок задачи"}`)
	}
	return nil
}

func (t *Task) CheckData() (time.Time, error) {
	if t.Date == "" {
		t.Date = time.Now().Format(layout)
	}
	parseDate, err := time.Parse(layout, t.Date)
	if err != nil {
		return time.Time{}, fmt.Errorf(`{"error":"Дата указана в неверном формате"}`)
	}
	return parseDate, nil
}

func (t *Task) CheckRepeat(parseDate time.Time) (string, error) {
	if t.Repeat != "" {
		nextDate, err := repeat.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			return "", fmt.Errorf(`{"error":"Неверное правило повторения"}`)
		}
		if parseDate.Before(time.Now()) && t.Date != time.Now().Format(layout) {
			t.Date = nextDate
		}
	} else if parseDate.Before(time.Now()) {
		t.Date = time.Now().Format(layout)
	}
	return t.Date, nil
}
