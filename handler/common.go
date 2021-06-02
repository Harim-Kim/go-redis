package handler

import (
	"OJT/model"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"time"
	pb "OJT/core"
)

type BaseHandler struct {
	CoreDB		*gorm.DB
	VehicleModelClient    *pb.VehiclemodelServiceClient
	RedisClient 		  *redis.Client
}

func WriteLog(code string, APIName string, response string, db *gorm.DB) {
	db.Create(&model.Log{APIName:APIName,APICallTime:time.Now(),APISuccess:code, APIResponseName:response})
	return
}

