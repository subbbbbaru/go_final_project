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

	// http.Handle("/", http.FileServer(http.Dir(webDir)))

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", handlers.NextDayHandler)
	mux.HandleFunc("/api/task", handlers.TaskHandler)
	mux.HandleFunc("/api/tasks", handlers.GetTasksHandler)
	mux.HandleFunc("/api/task/done", handlers.DoneTaskHandler)

	//go func() {
	if err := server.Run(config.Server.Port, mux); err != nil {
		log.Fatalf("Error while running http server %s", err.Error())
	}
	//}()

	// if err := server.Shutdown(context.Background()); err != nil {
	// 	log.Fatalf("error occured on server shutting down: %s", err.Error())
	// }
}
