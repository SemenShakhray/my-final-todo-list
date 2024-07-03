package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-final-project/db"
	repeat "go-final-project/rules-repeat"
	"log"
	"net/http"
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

type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func HandlerPostGetPutTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	db := db.CreateDataBase()
	defer db.Close()

	switch {
	case r.Method == http.MethodPost:
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}
		if task.Title == "" {
			http.Error(w, `{"error":"не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}
		if task.Date == "" {
			task.Date = time.Now().Format(layout)
		}
		parseDate, err := time.Parse(layout, task.Date)
		if err != nil {
			http.Error(w, `{"error":"дата указана в неверном формате"}`, http.StatusBadRequest)
			return
		}
		if task.Repeat != "" {
			nextDate, err := repeat.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"неверное правило повторения"}`, http.StatusBadRequest)
				return
			}
			if parseDate.Before(time.Now()) && task.Date != time.Now().Format(layout) {
				task.Date = nextDate
			}
		} else if parseDate.Before(time.Now()) {
			task.Date = time.Now().Format(layout)
		}
		query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
		res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка добавления задачи в базу данных"}`, http.StatusBadRequest)
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Ошибка получения ID вставленной записи"}`, http.StatusInternalServerError)
			return
		}

		resp := Response{ID: fmt.Sprintf("%d", id)}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}

	case r.Method == http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"нет индификатора задачи"}`, http.StatusBadRequest)
			return
		}
		row := db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id)
		err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(task); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}

	case r.Method == http.MethodPut:
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}
		if task.ID == "" {
			http.Error(w, `{"error":"не указан индификатор задачи"}`, http.StatusBadRequest)
			return
		}
		_, err = strconv.ParseInt(task.ID, 10, 32)
		if err != nil {
			http.Error(w, `{"error":"указан невозможный индификатор задачи"}`, http.StatusBadRequest)
			return
		}
		if task.Title == "" {
			http.Error(w, `{"error":"не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}
		if task.Date == "" {
			task.Date = time.Now().Format(layout)
		}
		parseDate, err := time.Parse(layout, task.Date)
		if err != nil {
			http.Error(w, `{"error":"дата указана в неверном формате"}`, http.StatusBadRequest)
			return
		}
		if task.Repeat != "" {
			nextDate, err := repeat.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"неверное правило повторения"}`, http.StatusBadRequest)
				return
			}
			if parseDate.Before(time.Now()) && task.Date != time.Now().Format(layout) {
				task.Date = nextDate
			}
		} else if parseDate.Before(time.Now()) {
			task.Date = time.Now().Format(layout)
		}
		query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
		_, err = db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			http.Error(w, `{"error":"Ошибка обновления задачи в базе данных"}`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	case r.Method == http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"не указан индификатор задачи"}`, http.StatusBadRequest)
			return
		}
		_, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			http.Error(w, `{"error":"указан невозможный индификатор задачи"}`, http.StatusBadRequest)
			return
		}
		query := "DELETE FROM scheduler WHERE id = ?"
		_, err = db.Exec(query, id)
		if err != nil {
			http.Error(w, `{"error":"не получается удалить задачу"}`, http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}

}

func HandlerGetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	var rows *sql.Rows
	var err error
	db := db.CreateDataBase()
	defer db.Close()
	search := r.URL.Query().Get("search")
	if search == "" {
		rows, err = db.Query("SELECT * FROM scheduler ORDER BY date LIMIT ?", 20)
	} else if date, error := time.Parse("02.01.2006", search); error == nil {
		query := "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = db.Query(query, date.Format(layout), 20)
	} else {
		search = "%" + search + "%"
		query := "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		rows, err = db.Query(query, search, search, 20)
	}
	if err != nil {
		http.Error(w, `{"error":"ошибка запроса"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"ошибка сканирования запроса"}`, http.StatusBadRequest)
			return
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		tasks = []Task{}
	}
	resp := map[string][]Task{
		"tasks": tasks,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
		return
	}
}

func HandlerNextDate(w http.ResponseWriter, r *http.Request) {
	db := db.CreateDataBase()
	defer db.Close()
	strnow := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	strRepeat := r.URL.Query().Get("repeat")

	now, err := time.Parse(layout, strnow)
	if err != nil {
		log.Fatal(err)
	}
	nextdate, err := repeat.NextDate(now, date, strRepeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write([]byte(nextdate))
}

func HandlerDone(w http.ResponseWriter, r *http.Request) {
	var task Task
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"не указан индификатор задачи"}`, http.StatusBadRequest)
		return
	}
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"указан невозможный индификатор задачи"}`, http.StatusBadRequest)
		return
	}
	db := db.CreateDataBase()
	defer db.Close()
	row := db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id)
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusBadRequest)
		return
	}
	if task.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id=?", task.ID)
		if err != nil {
			http.Error(w, `{"error":"не получается удалить задачу"}`, http.StatusBadRequest)
		}
	} else {
		next, err := repeat.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"неверное правило повторения"}`, http.StatusBadRequest)
			return
		}
		task.Date = next
	}
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = db.Exec(query, task.Date, task.ID)
	if err != nil {
		http.Error(w, `{"error":"Ошибка добавления задачи в базу данных"}`, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
		return
	}
}
