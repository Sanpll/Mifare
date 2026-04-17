package main

import (
	server "mifare/internal/app/apiserver"
	"mifare/internal/handler"
	"mifare/internal/repository"
	"mifare/internal/service"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Mifare REST API
// @version 1.0
// @description API server for Mifare processing

// @host localhost:8888
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description "Bearer {token}"

func main() {
	// считывание конфигов
	if err := initConfig(); err != nil {
		logrus.Fatal(err)
	}

	// подключение SQLite
	db, err := repository.NewSQLiteDB(repository.Config{
		DBPath: viper.GetString("db.DBPath"),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close() // выполнение этого метода откладывается, пока main не завершится

	// подключение миграций
	if err := goose.SetDialect("sqlite3"); err != nil {
		logrus.Fatal(err)
	}
	if err := goose.Up(db.DB, "./migrations"); err != nil {
		logrus.Fatal(err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(viper.GetString("app.api_version"), services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("app.port"), handlers.InitRoutes()); err != nil {
		logrus.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
