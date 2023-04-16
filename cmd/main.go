package main

import (
	usrHandler "Deall/handler"
	"Deall/helpers"
	md "Deall/middleware"
	"Deall/repository"
	"Deall/service"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
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
	dbName := viper.GetString(`database.name`)
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

	// Get the student collection
	collection := client.Database(dbName).Collection(helpers.UsersCollection)

	// Create a unique index on the name field
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	indexView := collection.Indexes()
	_, err = indexView.CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal("Failed to create index:", err)
	}

	userRepo := repository.NewUserRepository(client.Database(dbName))
	tokenRepo := repository.NewTokenRepository(client.Database(dbName))
	userService := service.NewUserService(userRepo, tokenRepo)
	userHandler := usrHandler.NewUserHandler(userService, time.Duration(timeoutCtx)*time.Second)
	middleware := md.NewMiddleware(userRepo, tokenRepo)

	serverPort := viper.GetString(`server.address`)
	handler := httprouter.New()

	handler.POST("/api/v1/login", userHandler.Login)
	handler.POST("/api/v1/user", middleware.AuthMiddleware(userHandler.RegistUser, []string{"Admin"}))
	handler.DELETE("/api/v1/user/:id", middleware.AuthMiddleware(userHandler.DeleteUser, []string{"Admin"}))
	handler.PATCH("/api/v1/user/:id", middleware.AuthMiddleware(userHandler.UpdateUser, []string{"Admin"}))
	handler.GET("/api/v1/users", middleware.AuthMiddleware(userHandler.GetAllUser, []string{"User", "Admin"}))
	handler.GET("/api/v1/user", middleware.AuthMiddleware(userHandler.GetUser, []string{"User", "Admin"}))

	logrus.Infof("Server run on localhost%v", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, handler))
}
