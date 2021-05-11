package main

import (
	"OJT/handler"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)
var (
	CoreDB *gorm.DB
	ctx = context.Background()
)
func ExampleClient(){
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	err := rdb.Set(ctx,"key", "value",0).Err()
	if err != nil{
		panic(err)
	}

	val,err := rdb.Get(ctx,"key").Result()
	if err != nil{
		panic(err)
	}
	fmt.Println("key",val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil{
		fmt.Println("Key2 deosn't exist")
	}else if err != nil{
		panic(err)
	}else{
		fmt.Println("key2", val2)
	}
}
func Init(){
	//log처리를 위한 DB Connect
	db, err := gorm.Open("mysql","root:root@tcp(127.0.0.1:3306)/vehicle_model?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		panic("failed to open database")
	}
	CoreDB = db


	CoreDB.Close()
}
func main(){
	//Init()
	ExampleClient()
	r := mux.NewRouter()
	r.HandleFunc("/",handler.LandingHandler).Methods(http.MethodGet) // Test
	r.HandleFunc("/vehicleModel",handler.VehiclesHandler).Methods(http.MethodGet)
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	http.ListenAndServe(":18080", r)
	fmt.Println("Don't Die..")
}

