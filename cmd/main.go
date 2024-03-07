package main

import (
	"context"
	"fmt"
	todo "ruby"
	"ruby/pkg/client/mongodb"
	"ruby/pkg/handler"
	"ruby/pkg/service"
	kb "ruby/pkg/user/db"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err!= nil{
		logrus.Fatalf("ошибка инициализации конфига , error initializing: %s", err.Error())
	}


 
  mongoDBclient, err := mongodb.NewClient(context.Background(), viper.GetString("mongodb.host"), viper.GetString("mongodb.port"),"", "", viper.GetString("mongodb.database"), "")
  if err != nil{
  	panic(err)
  }
 storage := kb.NewStorage(mongoDBclient, viper.GetString("mongodb.collection"))

	fmt.Println(storage)
	services:= service.NewService()
	handlers:= handler.NewHandler(services, storage)
	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err!=nil{
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}