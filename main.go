package main

import (
	"fmt"
	"go-final-project/auth"
	"go-final-project/handler"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", handler.HandlerNextDate)
	r.Post("/api/task", auth.Authorization(handler.HandlerPostGetPutTask))
	r.Get("/api/tasks", auth.Authorization(handler.HandlerGetTasks))
	r.Get("/api/task", handler.HandlerPostGetPutTask)
	r.Put("/api/task", handler.HandlerPostGetPutTask)
	r.Post("/api/task/done", auth.Authorization(handler.HandlerDone))
	r.Delete("/api/task", handler.HandlerPostGetPutTask)
	r.Post("/api/signin", auth.SigninHandler)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Println("ошибка запуска сервера:", err)
	}
}
