package handler

import (
	pb "OJT/core"
	"OJT/message"
	"OJT/model"
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"time"
)

type VehicleModelHandler BaseHandler

func stringInputConvertToUnintArray(s []string, N int )  []uint64{
	ret := make([]uint64,N)
	for i, v := range s{
		ret[i], _ = strconv.ParseUint(v,10,64)
	}
	return ret
}
func (h *VehicleModelHandler) VehiclesHandler(w http.ResponseWriter, r *http.Request ) {
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


	query := r.URL.Query()

	//전체
	if len(query) == 0{
		response, err := c.ListVehicleModel(ctx, &pb.ListVehicleModelRequest{Filter: &pb.VehicleModelFilter{}})
		log.Printf("result : %v",response)
		// log 저장
		if err != nil {
			h.CoreDB.Create(&model.Log{APIName:"ListVehicleModel",APICallTime:time.Now(),APISuccess:"FAIL", APIResponseName:response.String()})
			log.Fatalf("could not request: %v", err)
		}else{
			h.CoreDB.Create(&model.Log{APIName:"ListVehicleModel",APICallTime:time.Now(),APISuccess:"SUCCESS", APIResponseName:response.String()})
		}
		w.Write([]byte(response.String()))
	}else{
		//filter := pb.VehicleModelFilter{
		//	IDs: strings.Split(query.Get("IDs"),","),
		//	Brand: query.Get("Brand"),
		//	Brands: strings.Split(query.Get("Brands"), ","),
		//	Name: query.Get("Name"),
		//	FuelType: query.Get("FuelType"),
		//	FuelTypes: strings.Split(query.Get("FuelTypes"),","),
		//	Grade: query.Get("Grade"),
		//	Grades: strings.Split(query.Get("Grades"),","),
		//}
		//
		//response, err := c.ListVehicleModel(ctx, &pb.ListVehicleModelRequest{Filter: &filter})
		filter2 := message.Decoder(r)
		//fmt.Println("filter1:",filter)
		//fmt.Println("filter2:",filter2)
		response2, err := c.ListVehicleModel(ctx, &pb.ListVehicleModelRequest{Filter: &filter2})
		//fmt.Println("message using response",response2.String())
		// log 저장
		if err != nil {
			h.CoreDB.Create(&model.Log{APIName:"ListVehicleModel",APICallTime:time.Now(),APISuccess:"FAIL", APIResponseName:response2.String()})
			log.Fatalf("could not request: %v", err)
		}else{
			fmt.Println("creating log",response2.String())
			h.CoreDB.Create(&model.Log{APIName:"ListVehicleModel",APICallTime:time.Now(),APISuccess:"SUCCESS", APIResponseName:response2.String()})
		}
		w.Write([]byte(response2.String()))
	}

	//h.CoreDB.Commit()


	return
}
