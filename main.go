package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/muhammadnurbasari/onesmile-test-service-transaction/dialService"
	_ "github.com/muhammadnurbasari/onesmile-test-service-transaction/docs"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/muhammadnurbasari/onesmile-test-service-transaction/transaction"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// !Important : Comments below are formatted as it is to be read by Swagger tools (Swaggo)
// @title SERVICE TRANSACTION
// @version 1.0.0
// @description API DOCUMENTATION SERVICE TRANSACTION
// @termsOfService
// @contact.name API Support
// @contact.name ABBAS
// @contact.email m.nurbasari@gmail.com
// @BasePath
// @query.collection.format multi
// @securityDefinitions.apikey JWTToken
// @in header
// @name Authorization
func main() {
	var httpAddr = flag.String("http", ":8080", "http listen address")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	err := godotenv.Load("config/.env")

	if err != nil {
		level.Error(logger).Log("err : ", err.Error())
		os.Exit(1)
	}

	DIAL_CREDITCARD_SERVICE := os.Getenv("DIAL_CREDITCARD_SERVICE")
	grpcClientConnCC, err := grpc.Dial(DIAL_CREDITCARD_SERVICE, grpc.WithInsecure())
	if err != nil {
		level.Error(logger).Log("err : ", "could not connect to"+DIAL_CREDITCARD_SERVICE+" error: "+err.Error())
		os.Exit(1)
	}

	DIAL_HISTORY_SERVICE := os.Getenv("DIAL_HISTORY_SERVICE")
	grpcClientConnTrx, err := grpc.Dial(DIAL_HISTORY_SERVICE, grpc.WithInsecure())
	if err != nil {
		level.Error(logger).Log("err : ", "could not connect to"+DIAL_HISTORY_SERVICE+" error: "+err.Error())
		os.Exit(1)
	}

	flag.Parse()
	validateClient := dialService.ServiceCreditCard(grpcClientConnCC)
	transactionClient := dialService.ServiceHistory(grpcClientConnTrx)
	ctx := context.Background()
	srv := transaction.NewService(logger, validateClient, transactionClient)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	})

	endpoints := transaction.MakeEndpoints(srv)

	errs := make(chan error)
	go func() {
		fmt.Println("listening on port", *httpAddr)
		// swagger docs
		router.GET("/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		transaction.NewHttpServerTransaction(ctx, endpoints, router)
		errs <- router.Run(*httpAddr)
	}()

	level.Error(logger).Log("exit", <-errs)

}
