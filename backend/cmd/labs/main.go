package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// import (
// 	"backend/pkg/handler"
// 	"backend/pkg/model"
// 	"backend/pkg/repository"
// 	"backend/pkg/serverity"
// 	"backend/pkg/service"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/joho/godotenv"
// 	_ "github.com/lib/pq"
// 	"github.com/sirupsen/logrus"
// )

func main() {

	// if err := godotenv.Load(); err != nil {
	// 	logrus.Fatalf("error loading env varibles: %s", err.Error())
	// 	return
	// }
	// db, err := repository.NewPostgresDB(model.Config{
	// 	Host:     os.Getenv("DB_HOST"),
	// 	Port:     os.Getenv("DB_PORT"),
	// 	Username: os.Getenv("DB_USERNAME"),
	// 	Password: os.Getenv("DB_PASSWORD"),
	// 	DBName:   os.Getenv("DB_DBNAME"),
	// 	SSLMode:  os.Getenv("DB_SSL_MODE"),
	// })
	// if err != nil {
	// 	logrus.Fatalf("Fatal to connect to DB, because: %s", err.Error())
	// 	return
	// }

	// repo := repository.NewRepo(db)
	// services := service.NewService(repo)
	// handlers := handler.NewHandler(services)

	// srv := new(serverity.Server)

	// go func() {
	// 	if err := srv.Run(handlers.InitRoutes()); err != nil {
	// 		logrus.Fatalf("Problem with start server, because %s", err.Error())
	// 		return
	// 	}
	// }()
	// logrus.Println("backend started")
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	// <-quit
	// logrus.Println("backend shutting down")

	// if err := srv.ShutDown(context.Background()); err != nil {
	// 	logrus.Errorf("error occured on server shutting down: %s", err.Error())
	// 	return
	// }
		// Чтение данных из файла
	data, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	// Создание структуры для хранения данных
	var result map[string]interface{}

	// Преобразование JSON-данных в объект
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("Ошибка при преобразовании JSON:", err)
		return
	}
	fmt.Println(result["1"])
	// Доступ к объекту
	// matrix1 := result["1"].([]interface{})
	// for _, row := range matrix1 {
	// 	fmt.Println(row)
	// }
}
