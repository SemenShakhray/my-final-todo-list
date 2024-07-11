package main

import (
	"fmt"
	"go-final-project/auth"
	"go-final-project/config"
	"go-final-project/db"
	"go-final-project/interal/handler"
	"go-final-project/interal/storage"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

var Password string

func main() {
	env := config.GetEnv()
	Password = os.Getenv("TODO_PASSWORD")
	fmt.Println("Приложение запущено на порту", env.Port)

	dataBase := db.CreateDataBase()
	defer dataBase.Close()
	store := storage.NewStore(dataBase)

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", handler.HandlerNextDate)
	r.Post("/api/task", auth.Authorization(handler.HandlerPostGetPutTask(store)))
	r.Get("/api/tasks", auth.Authorization(handler.HandlerGetTasks(store)))
	r.Get("/api/task", handler.HandlerPostGetPutTask(store))
	r.Put("/api/task", handler.HandlerPostGetPutTask(store))
	r.Post("/api/task/done", auth.Authorization(handler.HandlerDone(store)))
	r.Delete("/api/task", handler.HandlerPostGetPutTask(store))
	r.Post("/api/signin", auth.SigninHandler)

	err := http.ListenAndServe(":"+env.Port, r)
	if err != nil {
		fmt.Println("ошибка запуска сервера:", err)
	}
}
