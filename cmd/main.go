package main

//(!) WaitGroup, File Writing, Change database filling

import (
	"log"
	"os"

	goroutinetask "github.com/2hard4me"
	"github.com/2hard4me/pkg/handlers"
	"github.com/2hard4me/pkg/logging"
	"github.com/2hard4me/pkg/models"
	"github.com/2hard4me/pkg/repository"
	"github.com/2hard4me/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)


func main() {
	logger := logging.GetLogger()

	if err := godotenv.Load(".env"); err != nil {
		logger.Fatal(err)
	}

	if err := initConfig(); err != nil {
		logger.Fatal(err)
	}

	db, err := repository.NewConnection(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASS"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatal("could not load the database")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		logger.Fatal("could not migrate db")
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handlers.NewHandler(services)

	srv := new(goroutinetask.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logger.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
