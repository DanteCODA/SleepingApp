package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/config"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/consts"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/infrastructure/repositories/repos"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/infrastructure/scraper"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/usecase/assets"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/usecase/checkpoint"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/usecase/price"
)

func main() {
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context) {
	log.Println("lambda handler is called")

	appConf := config.AppConf

	// create new logger
	zap, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal("create app logger failed")
	}
	defer zap.Close()

	// crea