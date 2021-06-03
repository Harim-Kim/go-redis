package main

import (
	pb "OJT/core"
	"OJT/handler"
	"context"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	"google.golang.org/grpc"
	"log"
	"net/http"
)
var (
	CoreDB 				  *gorm.DB
	VehicleModelClient    *pb.VehiclemodelServiceClient
	client 			= 	  &redisClient{}
)
type redisClient struct{
	c *redis.Client
}
func RedisInit() (*redis.Client)  {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil{
		log.Println("redis client error")
		return redisClient
	}
	client.c = redisClient
	return redisClient

}
func init(){
	//db 연결을 위한 connection
	db, err := gorm.Open("mysql","root:root@tcp(127.0.0.1:3306)/vehicle_model?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		panic("failed to open database")
	}
	CoreDB = db

}

// @title open API(Swagger)
// @version 1.1.1
// @description This is Open Api Document Server(Swagger)
// @contact.email harimkim@hyundai-autoever.com
// @host localhost:18080
// @BasePath /
func main(){
	RedisInit()
	r := mux.NewRouter()
	HandleRoutes(r)
	makeConnection()
	// Swagger
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	http.ListenAndServe(":18080", r)
	CoreDB.Close()
}
func makeConnection(){
	conn, err := grpc.Dial("localhost:12005", grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect core: %v",err)
	}
	defer conn.Close()

	c := pb.NewVehiclemodelServiceClient(conn)
	// grpc gateway와 grpc 서버 이어 주기 위한 context 선언
	//ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	//defer cancel()
	VehicleModelClient = &c
}
func HandleRoutes(router *mux.Router){
	//vehicle
	vh := &handler.VehicleModelHandler{
		CoreDB: CoreDB,
		VehicleModelClient: VehicleModelClient,

	}
	vh.RedisInit()
	router.HandleFunc("/vehiclemodel", vh.VehicleModelList).Methods(http.MethodGet)                                // 차량 모델 목록
	router.HandleFunc("/vehiclemodel/{vehiclemodelID}", vh.VehicleModelGet).Methods(http.MethodGet) 				// 차량 모델 조회
	router.HandleFunc("/vehiclemodel", vh.VehicleModelInsert).Methods(http.MethodPost)                             // 차량 모델 등록
	router.HandleFunc("/vehiclemodel/{vehiclemodelID}", vh.VehicleModelUpdate).Methods(http.MethodPatch)           // 차량 모델 수정
	router.HandleFunc("/vehiclemodel/{vehiclemodelID}", vh.VehicleModelDelete).Methods(http.MethodDelete)          // 차량 모델 삭제
	router.HandleFunc("/vehiclemodel_csv", vh.VehicleModelCSV).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}