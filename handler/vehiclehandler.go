package handler

import (
	pb "OJT/core"
	"OJT/message"
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type VehicleModelHandler BaseHandler

// VehicleModelList godoc
// @Summary Show a VehicleModelList
// @Description VehicleModelList를 조회한다
// @Accept  html
// @Produce  json
// @Param name query string false "모델명"
// @Param brand query string false "제조사"
// @Param brands query string false "제조사[]"
// @Param fuelType query string false "연료 유형"
// @Param fuelTypes query string false "연료 유형[]"
// @Param grade query string false "차량 등급"
// @Param grades query string false "차량 등급[]"
// @Success 200 {object} model.Vehiclemodel
// @Router /vehiclemodel [get]
func (h *VehicleModelHandler) VehicleModelList(w http.ResponseWriter, r *http.Request ) {
	fmt.Println("목록 조회")
	
	conn, err := grpc.Dial("localhost:12005", grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Println("did not connect core: %v",err)
	}
	defer conn.Close()

	c := pb.NewVehiclemodelServiceClient(conn)

	//전체
	if len(r.URL.Query()) == 0{
		response, err := c.ListVehicleModel(r.Context(), &pb.ListVehicleModelRequest{Filter: &pb.VehicleModelFilter{}})
		// log 저장
		if err != nil {
			//이렇게 db 정보 넘겨서 create하는 경우 created at, updated at 이 입력 안 됨.
			WriteLog("FAIL", "ListVehicleModel", response.String(),h.CoreDB)
	 		log.Println("could not request: %v", err)
		}else{
			WriteLog("SUCCESS", "ListVehicleModel", response.String(), h.CoreDB)
		}
		w.Write([]byte(response.String()))

		//redis update
		fmt.Println(response.VehicleModels[0])
		for _, vehicle := range response.VehicleModels{
			_, err := h.RedisClient.Get(vehicle.ID).Result()
			if err == redis.Nil {
				h.VehicleModelRedisSet(vehicle.ID,vehicle,time.Hour*1)
			} else if err != nil {
				log.Println("Redis db error",err.Error())
			} else {
				continue
			}
		}

	}else{

		filter2 := message.Decoder(r)

		//reqData2Core := &pb.ListVehicleModelRequest{Filter: &filter2}
		//이렇게 connection 정보 넘기는 경우 error 발생.
		//result, err := (*h.VehicleModelClient).ListVehicleModel(ctx, reqData2Core)
		response2, err := c.ListVehicleModel(r.Context(), &pb.ListVehicleModelRequest{Filter: &filter2})
		//fmt.Println(result.String())
		// log 저장
		if err != nil {
			WriteLog("FAIL", "ListVehicleModel", response2.String(),h.CoreDB)
			log.Println("could not request: %v", err)
		}else{
			WriteLog("SUCCESS", "ListVehicleModel", response2.String(),h.CoreDB)
		}
		w.Write([]byte(response2.String()))
		conn.Close()
		fmt.Println(response2.VehicleModels[0])
		for _, vehicle := range response2.VehicleModels{
			_, err := h.RedisClient.Get(vehicle.ID).Result()
			if err == redis.Nil {
				h.VehicleModelRedisSet(vehicle.ID,vehicle,time.Hour*1)
			} else if err != nil {
				log.Println("Redis db error",err.Error())
			} else {
				continue
			}
		}
	}

	return
}
//  VehicleModelGet godoc
// @Summary Show a  VehicleModelGet
// @Description VehicleModel를 조회한다
// @Accept  html
// @Produce  json
// @Param vehicleID path string true "모델 ID"
// @Success 200 {object} model.Vehiclemodel
// @Router /vehiclemodel/{vehiclemodelID} [get]
func (h *VehicleModelHandler) VehicleModelGet(w http.ResponseWriter, r *http.Request ){
	fmt.Println("모델 조회")
	var vehicle_id uuid.UUID
	queryParams := mux.Vars(r)
	if param, ok := queryParams["vehiclemodelID"]; !ok {
		WriteLog("FAIL", "VehicleModelGet", "cannot find vehiclemodelID values",h.CoreDB)
		return
	} else if uuidValue, err := uuid.Parse(param); err != nil {
		WriteLog("FAIL", "VehicleModelGet", "vehiclemodelID must have uuid format" ,h.CoreDB)
		return
	} else {
		vehicle_id = uuidValue
	}

	conn, err := grpc.Dial("localhost:12005", grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Println("did not connect core: %v",err)
		return
	}
	defer conn.Close()

	c := pb.NewVehiclemodelServiceClient(conn)
	ctx := r.Context()
	result, err := c.GetVehicleModel(ctx, &pb.VehicleModelID{
		ID: vehicle_id.String(),
	})
	if err != nil {
		WriteLog("FAIL", "VehicleModelGet", "vehiclemodelID does not exist!" ,h.CoreDB)
		return
	}
	w.Write([]byte(result.String()))
	//out := &message.VehicleModelDetail{}
	//out.Convert(result)
	//
	w.WriteHeader(http.StatusOK)
	//_ = json.NewEncoder(w).Encode(out)



	//redis test
	keyRedis := "vehicleRedisTest"
	err = h.VehicleModelRedisSet(keyRedis,result,time.Hour*1)
	if err != nil{
		log.Println("Redis set error : %v", err.Error())
	}
	vehicle_detail := &message.VehicleModelDetail{}
	err = h.VehicleModelRedisGet( keyRedis, vehicle_detail)
	if err != nil{
		log.Println("Redis get error : %v", err.Error())
	}
	w.Write([]byte("\n\n---------------------FROM REDIS--------------------------\n\n"))
	// 요부분은 나중에
	vehicle_detail.Grade = string(pb.VehicleModelGrade_value[vehicle_detail.Grade])
	vehicle_from_redis, err := json.Marshal(vehicle_detail)
	if err != nil {
		panic (err)
	}
	w.Write([]byte(string(vehicle_from_redis)))
	return
}
// VehicleModelInsert godoc
// @Summary Show a VehicleModelInsert
// @Description VehicleModel을 생성한다.
// @Accept  html
// @Produce  json
// @Param name body string true "모델명"
// @Param brand body string true "제조사"
// @Param standard body bool false "대표 모델 여부"
// @Param standardModelID body string false "대표 모델 ID"
// @Param seatingCapacity body uint32 false "좌석수"
// @Param fuelType body string false "연료 유형"
// @Param fuelEfficiency body float32 false "연비"
// @Param displacement body uint32 false "주행 거리"
// @Param grade body string true "차량 등급"
// @Param warmUpTime body uint32 false "웜업 시간"
// @Success 200 {object} model.Vehiclemodel
// @Router /vehiclemodel [post]
func (h *VehicleModelHandler) VehicleModelInsert(w http.ResponseWriter, r *http.Request ) {
	fmt.Println("등록 ")
	reqData := &message.VehicleModelRequest{}
	reqData.Brand = "Kia"
	err := reqData.Validate(r)
	if err != nil {
		WriteLog("FAIL", "VehicleModelInsert", "invalid inputs",h.CoreDB)
		return
	}
	conn, err := grpc.Dial("localhost:12005", grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Println("did not connect core: %v",err)
		return
	}
	defer conn.Close()
	c := pb.NewVehiclemodelServiceClient(conn)
	{

		resData, err := c.RegisterVehicleModel(r.Context(), &pb.RegisterVehicleModelRequest{
			Standard:         reqData.Standard,
			StandardModelID:  reqData.StandardModelID,
			Name:             reqData.Name,
			Brand:            reqData.Brand,
			FuelType:         reqData.FuelType,
			SeatingCapacity:  reqData.SeatingCapacity,
			Displacement:     reqData.Displacement,
			FuelTankCapacity: reqData.FuelTankCapacity,
			FuelEfficiency:   reqData.FuelEfficiency,
			Grade:            pb.VehicleModelGrade(pb.VehicleModelGrade_value[reqData.Grade]),
			WarnUpTime:       reqData.WarmUpTime,
		})
		if err != nil {
			WriteLog("FAIL", "VehicleModelInsert", "Core DB Error",h.CoreDB)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(resData.String()))
		WriteLog("SUCCESS", "VehicleModelInsert", resData.String(),h.CoreDB)

		//redis 반영
		key_redis_vehicle := resData.ID
		err = h.VehicleModelRedisSet(key_redis_vehicle,resData,time.Hour*1)
		if err != nil{
			log.Println("Redis set error : %v", err.Error())
		}
	}

}
// VehicleModelUpdate godoc
// @Summary Show a VehicleModelUpdate
// @Description VehicleModel을 수정한다.
// @Accept  html
// @Produce  json
// @Param name body string true "모델명"
// @Param brand body string true "제조사"
// @Param standard body bool false "대표 모델 여부"
// @Param standardModelID body string false "대표 모델 ID"
// @Param seatingCapacity body uint32 false "좌석수"
// @Param fuelType body string false "연료 유형"
// @Param fuelEfficiency body float32 false "연비"
// @Param displacement body uint32 false "주행 거리"
// @Param grade body string true "차량 등급"
// @Param warmUpTime body uint32 false "웜업 시간"
// @Param vehicleID path string true "모델 ID"
// @Success 200 {object} model.Vehiclemodel
// @Router /vehiclemodel/{vehiclemodelID} [patch]
func (h *VehicleModelHandler) VehicleModelUpdate(w http.ResponseWriter, r *http.Request ) {
	fmt.Println("수정")
	conn, err := grpc.Dial("localhost:12005", grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Println("did not connect core: %v",err)
		return
	}
	defer conn.Close()

	c := pb.NewVehiclemodelServiceClient(conn)
	var vehiclemodelID uuid.UUID
	reqData := &message.VehicleModelRequest{}
	{
		queryParams := mux.Vars(r)
		if param, ok := queryParams["vehiclemodelID"]; !ok {
			WriteLog("FAIL", "VehicleModelUpdate", "cannot find vehiclemodelID values",h.CoreDB)
			return
		} else if uuidValue, err := uuid.Parse(param); err != nil {
			WriteLog("FAIL", "VehicleModelUpdate", "vehiclemodelID must have uuid format" ,h.CoreDB)
			return
		} else {
			vehiclemodelID = uuidValue
		}
		reqData.Brand = "Kia"

		err := reqData.Validate(r)
		if err != nil {
			WriteLog("FAIL", "VehicleModelUpdate", "http.StatusBadRequest" ,h.CoreDB)

			return
		}
	}


	{
		ctx := r.Context()
		result, err := c.UpdateVehicleModel(ctx, &pb.UpdateVehicleModelRequest{
			UpdateMask: &field_mask.FieldMask{
				Paths: []string{
					"standard", "standardmodelid", "name", "brand", "fueltype", "fuelefficiency", "fueltankcapacity",
					"seatingcapacity", "displacement", "grade", "warmuptime",
				},
			},
			ID:               vehiclemodelID.String(),
			Standard:         reqData.Standard,
			StandardModelID:  reqData.StandardModelID,
			Name:             reqData.Name,
			Brand:            reqData.Brand,
			FuelType:         reqData.FuelType,
			FuelEfficiency:   reqData.FuelEfficiency,
			FuelTankCapacity: reqData.FuelTankCapacity,
			SeatingCapacity:  reqData.SeatingCapacity,
			Displacement:     reqData.Displacement,
			Grade:            pb.VehicleModelGrade(pb.VehicleModelGrade_value[reqData.Grade]),
			WarmUpTime:       reqData.WarmUpTime,
		})
		if err != nil {
			WriteLog("FAIL", "VehicleModelUpdate", "DB Error" ,h.CoreDB)
			return
		}

		result, err = (*h.VehicleModelClient).GetVehicleModel(ctx, &pb.VehicleModelID{
			ID: result.ID,
		})
		if err != nil {
			WriteLog("FAIL", "VehicleModelUpdate", "DB Error" ,h.CoreDB)
			return
		}
		w.Write([]byte(result.String()))
		WriteLog("SUCCESS", "VehicleModelUpdate", result.String() ,h.CoreDB)
		//out := &message.VehicleModelDetail{}
		//out.Convert(result)

		w.WriteHeader(http.StatusOK)
		//_ = json.NewEncoder(w).Encode(out)

		//redis 반영
		key_redis_vehicle := result.ID
		err = h.VehicleModelRedisSet(key_redis_vehicle,result,time.Hour*1)
		if err != nil{
			log.Println("Redis set error : %v", err.Error())
		}
	}
	return
}
// VehicleModelDelete godoc
// @Summary Show a  VehicleModelDelete
// @Description VehicleModel을 삭제한다.
// @Accept  html
// @Produce  json
// @Param vehicleID path string true "모델 ID"
// @Success 200 {object} model.Vehiclemodel
// @Router /vehiclemodel/{vehiclemodelID} [delete]
func (h *VehicleModelHandler) VehicleModelDelete(w http.ResponseWriter, r *http.Request ) {
	fmt.Println("삭제")
	conn, err := grpc.Dial("localhost:12005", grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Println("did not connect core: %v",err)
		return
	}
	defer conn.Close()

	c := pb.NewVehiclemodelServiceClient(conn)

	var vehiclemodelID uuid.UUID
	queryParams := mux.Vars(r)
	if param, ok := queryParams["vehiclemodelID"]; !ok {
		WriteLog("FAIL", "VehicleModelDelete", "cannot find vehiclemodelID values" ,h.CoreDB)
		return
	} else if uuidValue, err := uuid.Parse(param); err != nil {
		WriteLog("FAIL", "VehicleModelDelete", "vehiclemodelID must have uuid format" ,h.CoreDB)
		return
	} else {
		vehiclemodelID = uuidValue
	}

	ctx := r.Context()
	//count, err := (*h.VehicleClient).CountVehicle(ctx, &core.CountVehicleRequest{
	//	Filter: &pb.VehicleFilter{
	//		ModelID: vehiclemodelID.String(),
	//	},
	//})
	//if err != nil {
	//	replyError(w, h.Logger, utils.ConvertGRPCError(err), err)
	//	return
	//}
	//if count.Count > 0 {
	//	replyError(w, h.Logger, http.StatusConflict, errors.New("CANNOT_DELETE_VEHICLEMODEL_DUE_TO_ASSOCIATED_VEHICLES"))
	//	return
	//}

	result, err := c.DeleteVehicleModel(ctx, &pb.VehicleModelID{
		ID: vehiclemodelID.String(),
	})
	if err != nil {
		WriteLog("FAIL", "VehicleModelDelete", "Core DB Error" ,h.CoreDB)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteLog("SUCCESS", "VehicleModelDelete", result.String(),h.CoreDB)
	w.Write([]byte(result.String()))

	//redis delete
	h.VehicleModelRedisDelete(vehiclemodelID.String())
}
func (h *VehicleModelHandler) VehicleModelCSV (w http.ResponseWriter, r *http.Request ) {
	fmt.Println("csv..")
	conn, err := grpc.Dial("localhost:12005", grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Println("did not connect core: %v",err)
	}
	defer conn.Close()

	c := pb.NewVehiclemodelServiceClient(conn)

	//전체

	response, err := c.ListVehicleModel(r.Context(), &pb.ListVehicleModelRequest{Filter: &pb.VehicleModelFilter{}})
	// log 저장
	if err != nil {
		//이렇게 db 정보 넘겨서 create하는 경우 created at, updated at 이 입력 안 됨.
		WriteLog("FAIL", "ListVehicleModel", response.String(),h.CoreDB)
		log.Println("could not request: %v", err)
	}else{
		WriteLog("SUCCESS", "ListVehicleModel", response.String(), h.CoreDB)
	}
	w.Write([]byte(response.String()))

	//redis update
	file, err := os.Create("D:\\maas\\go\\src\\OJT\\csv_vehicle.csv")
	if err != nil{
		log.Println("CSV File create error")
	}
	wr := csv.NewWriter(bufio.NewWriter(file))
	wr.Write([]string{"ID","Name","Brand","Standard","StandardModelID","SeatingCapacity","FuelType","FuelEfficiency","FuelTankCapacity","Displacement","Grade","WarmUpTime","ImageURL","CreatedAT","UpdatedAT"})


	for _, v := range response.VehicleModels{
		wr.Write([]string{v.GetID(),v.GetName(),v.GetBrand(),strconv.FormatBool(v.GetStandard()),v.GetStandardModelID(),strconv.FormatUint(uint64(v.GetSeatingCapacity()), 10),v.GetFuelType(),strconv.FormatFloat(float64(v.GetFuelEfficiency()),'f',-1,32),strconv.FormatUint(uint64(v.GetFuelTankCapacity()),10),strconv.FormatUint(uint64(v.GetDisplacement()),10),v.GetGrade().String(),strconv.FormatUint(uint64(v.GetWarmUpTime()),10),v.GetImageURL(),strconv.FormatInt(v.GetCreatedAt(),10),strconv.FormatInt(v.GetUpdatedAt(),10)})

	}
	wr.Flush()

	return
}
func (r *VehicleModelHandler) VehicleModelRedisGet( key string, src interface{}) error {

	val, err := r.RedisClient.Get(key).Result()
	if err == redis.Nil || err != nil{
		return err
	}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil{
		return err
	}
	return nil
}
func (r *VehicleModelHandler) VehicleModelRedisSet( key string, value interface{}, expiration time.Duration) error {

	cacheEntry, err := json.Marshal(value)
	if err != nil{
		return err
	}
	err = r.RedisClient.Set( key,  cacheEntry, expiration).Err()
	if err != nil{
		return err
	}
	return nil
}
func (r *VehicleModelHandler) VehicleModelRedisDelete(key string, )  {

	r.RedisClient.Del(key)

	return
}
func (h *VehicleModelHandler)RedisInit() (*redis.Client)  {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil{
		log.Println("redis client error")
		return redisClient
	}
	h.RedisClient = redisClient
	return redisClient

}