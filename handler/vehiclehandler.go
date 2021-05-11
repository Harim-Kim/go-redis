package handler

import (
	pb "OJT/core"
	"OJT/model"
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

)

type VehicleModelHandler BaseHandler

func LandingHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello world"))
	return
}
func stringInputConvertToUnintArray(s []string, N int )  []uint64{
	ret := make([]uint64,N)
	for i, v := range s{
		ret[i], _ = strconv.ParseUint(v,10,64)
	}
	return ret
}
func VehiclesHandler(w http.ResponseWriter, r *http.Request ) {
	fmt.Println("Working,,")
	conn, err := grpc.Dial("localhost:12005", grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect core: %v",err)
	}
	defer conn.Close()

	c := pb.NewVehiclemodelServiceClient(conn)
	// grpc gateway와 grpc 서버 이어 주기 위한 context 선언
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()

	//log처리를 위한 DB Connect
	db, err := gorm.Open("mysql","root:root@tcp(127.0.0.1:3306)/vehicle_model?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		panic("failed to open database")
	}

	query := r.URL.Query()
	//전체
	if len(query) == 0{
		response, err := c.ListVehicleModel(ctx, &pb.ListVehicleModelRequest{Filter: &pb.VehicleModelFilter{}})
		log.Printf("result : %v",response)
		// log 저장
		if err != nil {
			db.Create(&model.Log{APIName:"ListVehicleModel",APICallTime:time.Now(),APISuccess:"FAIL", APIResponseName:response.String()})
			log.Fatalf("could not request: %v", err)
		}else{
			db.Create(&model.Log{APIName:"ListVehicleModel",APICallTime:time.Now(),APISuccess:"SUCCESS", APIResponseName:response.String()})
		}
		w.Write([]byte(response.String()))
	}else{
		//데이터 전처리
		var pageReq uint64
		if query.Get("Page") != ""{
			temp, errReq := strconv.ParseUint(query.Get("Page"),10,64)
			if errReq != nil{
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			pageReq = temp
		}
		var rowPerPage uint64
		if query.Get("RowPerPage") != ""{
			temp, errReq := strconv.ParseUint(query.Get("RowPerPage"),10,64)
			if errReq != nil{
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			rowPerPage = temp
		}


		var sortField uint64
		if query.Get("SortFiled") != "" {
			temp, errReq := strconv.ParseUint(query.Get("SortField"),10,64)
			if errReq != nil{
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			sortField = temp
		}

		var sortOrder uint64
		if query.Get("SortOrder") != ""{
			temp, errReq := strconv.ParseUint(query.Get("SortOrder"),10,64)
			if errReq != nil{
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			sortOrder = temp
		}
		var sortOrders []uint64
		var sortFields []uint64
		if query.Get("SortOrders") == ""{
			sortOrders = nil
		}else{
			strArray := strings.Split(query.Get("SortOrders"),",")
			sortOrders = stringInputConvertToUnintArray(strArray, len(strArray))
		}
		if query.Get("SortFields") == ""{
			sortFields = nil
		}else{
			strArray := strings.Split(query.Get("SortFields"),",")
			stringInputConvertToUnintArray(strArray, len(strArray))
		}
		filter := pb.VehicleModelFilter{
			IDs: strings.Split(query.Get("IDs"),","),
			Brand: query.Get("Brand"),
			Brands: strings.Split(query.Get("Brands"), ","),
			Name: query.Get("Name"),
			FuelType: query.Get("FuelType"),
			FuelTypes: strings.Split(query.Get("FuelTypes"),","),
			Grade: query.Get("Grade"),
			Grades: strings.Split(query.Get("Grades"),","),
			Page: pageReq,
			RowPerPage:  rowPerPage,
			SortField: sortField,
			SortFields: sortFields,
			SortOrder: sortOrder,
			SortOrders: sortOrders,
		}

		response, err := c.ListVehicleModel(ctx, &pb.ListVehicleModelRequest{Filter: &filter})

		// log 저장
		if err != nil {
			db.Create(&model.Log{APIName:"ListVehicleModel",APICallTime:time.Now(),APISuccess:"FAIL", APIResponseName:response.String()})
			log.Fatalf("could not request: %v", err)
		}else{
			db.Create(&model.Log{APIName:"ListVehicleModel",APICallTime:time.Now(),APISuccess:"SUCCESS", APIResponseName:response.String()})
		}
		w.Write([]byte(response.String()))
	}

	db.Commit()

	db.Close()
	return
}
