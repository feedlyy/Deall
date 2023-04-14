package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func init() {
	viper.SetConfigFile(`../config/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// Setup Logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	// Connect to MongoDB
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	mongoUri := fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	timeoutCtx := viper.GetInt(`context.timeout`)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutCtx)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// test ping mongo
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	logrus.Info("Got ping from mongodb")

	serverPort := viper.GetString(`server.address`)
	handler := httprouter.New()

	logrus.Infof("Server run on localhost%v", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, handler))
}
