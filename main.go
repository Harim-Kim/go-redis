package main

import (
	pb "OJT/core"
	"OJT/handler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"log"
	"net/http"
)
var (
	CoreDB 				  *gorm.DB
	VehicleModelClient    *pb.VehiclemodelServiceClient
	//client 			= 	  &redisClient{}
)

func init(){
	//db 연결을 위한 connection
	db, err := gorm.Open("mysql","root:root@tcp(127.0.0.1:3306)/vehicle_model?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		panic("failed to open database")
	}
	CoreDB = db

}

// @title Hellow
// @version 1.1.1
// @description This is vehicle Api Document Server(Swagger)
// @contact.email harimkim@hyundai-autoever.com
// @host localhost:18080
// @BasePath /
func main(){
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