package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	myLog "github.com/subbbbbaru/first-sample/pkg/log"
	"github.com/subbbbbaru/go_final_project/configs"
	"github.com/subbbbbaru/go_final_project/internal/api"
	"github.com/subbbbbaru/go_final_project/internal/handlers"
	"github.com/subbbbbaru/go_final_project/internal/repository"
	"github.com/subbbbbaru/go_final_project/internal/service"
)

const webDir = "web"

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		myLog.Info().Println("No .env file found")
	}
}

func main() {
	// Инициализация логгеров
	myLog.InitLoggers(os.Stdout, os.Stderr)

	config := configs.New()

	dbsqlite, err := repository.NewSQLite3DB(config.DB.Name)
	if err != nil {
		myLog.Error().Fatalln(err)
	}

	repos := repository.NewRepository(dbsqlite)
	services := service.NewService(repos)
	handlers := handlers.NewHandler(services)

	server := new(api.Server)

	mux := http.NewServeMux()

	pass, err := repos.GetPassword()
	if err != nil {
		myLog.Error().Fatalln(err)
	}

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", handlers.NextDayHandler)
	mux.HandleFunc("/api/task", handlers.UserIdentity(handlers.TaskHandler, pass))
	mux.HandleFunc("/api/tasks", handlers.UserIdentity(handlers.GetTasksHandler, pass))
	mux.HandleFunc("/api/task/done", handlers.UserIdentity(handlers.DoneTaskHandler, pass))
	mux.HandleFunc("/api/signin", handlers.SignIn)

	if err := server.Run(config.Server.Port, mux); err != nil {
		myLog.Error().Fatalf("Error while running http server %s", err.Error())
	}
}
