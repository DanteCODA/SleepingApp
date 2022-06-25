package models

import (
	"context"
	"time"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AssetP