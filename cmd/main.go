package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/subbbbbaru/go_final_project/configs"
	"github.com/subbbbbaru/go_final_project/internal/api"
	"github.com/subbbbbaru/go_final_project/internal/handlers"
	"github.com/subbbbbaru/go_final_project/internal/repository"
	"github.com/subbbbbaru/go_final_project/internal/service"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	config := configs.New()

	dbsqlite, err := repository.NewSQLite3DB(config.DB.Name)
	if err != nil {
		log.Fatalf("[Error] %v", err)
	}

	repos := repository.NewRepository(dbsqlite)
	services := service.NewService(repos)
	handlers := handlers.NewHandler(services)

	server := new(api.Server)

	mux := http.NewServeMux()

	webDir := "web"

	pass, err := repos.GetPassword()
	if err != nil {
		log.Fatalf("[Error] %v", err)
	}

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", handlers.UserIdentity(handlers.NextDayHandler, pass))
	mux.HandleFunc("/api/task", handlers.UserIdentity(handlers.TaskHandler, pass))
	mux.HandleFunc("/api/tasks", handlers.UserIdentity(handlers.GetTasksHandler, pass))
	mux.HandleFunc("/api/task/done", handlers.UserIdentity(handlers.DoneTaskHandler, pass))
	mux.HandleFunc("/api/signin", handlers.SignIn)

	if err := server.Run(config.Server.Port, mux); err != nil {
		log.Fatalf("Error while running http server %s", err.Error())
	}
}
