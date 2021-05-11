package model

import (
	"time"
)

type Log struct{
	APIName string `gorm:"column:api_name;varchar(30);index"`
	APICallTime time.Time
	APISuccess string `gorm:"column:api_success;varchar(10);index"`//SUCESS, FAIL
	APIResponseName string `gorm:"column:api_response_name;varchar(10);index"`
}
func(Log) TableName() string{
	return "log"
}