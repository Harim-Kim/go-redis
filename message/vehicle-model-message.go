package message

import (
	pb "OJT/core"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"errors"
)
type VehicleModelListReq struct {
	//ID         []string	`"json:"ID`
	IDs        []string   `json:"IDs,omitempty"`        // 모델 ID[]
	Brand      []string   `json:"brand"`                // 상표
	Brands     []string   `json:"brands,omitempty"`     // 상표[]
	Name       []string   `json:"name,omitempty"`       // 모델명
	//Names 	   []string `"json:"names,omitempty"`
	FuelType   []string   `json:"fuelType,omitempty"`   // 연료 유형
	FuelTypes  []string   `json:"fuelTypes,omitempty"`  // 연료 유형[]
	Grade      []string   `json:"grade,omitempty"`      // 차량 등급
	Grades     []string   `json:"grades,omitempty"`     // 차량 등급[]

}
type VehicleModelDetail struct {
	ID               string  `json:"ID"`
	Name             string  `json:"name"`
	Brand            string  `json:"brand"`
	Standard         bool    `json:"standard,bool"`
	StandardModelID  string  `json:"standardModelID"`
	SeatingCapacity  uint32  `json:"seatingCapacity,int"`
	FuelType         string  `json:"fuelType"`
	FuelEfficiency   float32 `json:"fuelEfficiency,float32"`
	FuelTankCapacity uint32  `json:"fuelTankCapacity,int"`
	Displacement     uint32  `json:"displacement,int"`
	Grade            string  `json:"grade,string"`
	WarmUpTime       uint32  `json:"warmUpTime,int"`
	ImageURL         string  `json:"imageURL"`
}
type VehicleModelRequest struct {
	Name             string  `json:"name"`
	Brand            string  `json:"brand"`
	Standard         bool    `json:"standard,string"`
	StandardModelID  string  `json:"standardModelID"`
	SeatingCapacity  uint32  `json:"seatingCapacity,string"`
	FuelType         string  `json:"fuelType"`
	FuelEfficiency   float32 `json:"fuelEfficiency,string"`
	FuelTankCapacity uint32  `json:"fuelTankCapacity,string"`
	Displacement     uint32  `json:"displacement,string"`
	Grade            string  `json:"grade"`
	WarmUpTime       uint32  `json:"warmUpTime,string"`
}

func  Decoder(r *http.Request) pb.VehicleModelFilter{
	query := r.URL.Query()
	jsonbody, err := json.Marshal(query)
	if err != nil{
		fmt.Println("marshal error:",err)
	}
	middle := VehicleModelListReq{}
	if err := json.Unmarshal(jsonbody, &middle); err != nil {
		fmt.Println("unmarshal middle error:",err)
	}
	filter := pb.VehicleModelFilter{
		Brand: strings.Join(middle.Brand,""),
		Brands: middle.Brands,
		Name: strings.Join(middle.Name,""),
		FuelType: strings.Join(middle.FuelType,""),
		FuelTypes: middle.FuelTypes,
		Grade: strings.Join(middle.Grade,""),
		Grades: middle.Grades,
	}
	return filter
}
func (data *VehicleModelRequest) Validate(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return err
	}

	if data.Name == "" {
		return errors.New("EMPTY_NAME")
	}
	if data.Brand == "" {
		return errors.New("EMPTY_BRAND")
	}
	if data.Grade == "" {
		return errors.New("EMPTY_GRADE")
	} else if value, ok := pb.VehicleModelGrade_value[data.Grade]; !ok || value <= 0 {
		return errors.New("INVALID_GRADE")
	}
	if data.WarmUpTime <= 0 {
		return errors.New("INVALID_WARM_UP_TIME")
	}

	if !data.Standard {
		if data.StandardModelID == "" {
			return errors.New("EMPTY_STANDARD_MODEL_ID")
		}
		if data.FuelType == "" {
			return errors.New("EMPTY_FUEL_TYPE")
		}
	}

	return nil
}
//func validateRequest(r *http.Request) error{
//	err := json.NewDecoder(r.Body).Decode(data)
//	if err != nil {
//		return err
//	}
//
//	if data.Name == "" {
//		return errors.New("EMPTY_NAME")
//	}
//	if data.Brand == "" {
//		return errors.New("EMPTY_BRAND")
//	}
//	if data.Grade == "" {
//		return errors.New("EMPTY_GRADE")
//	} else if value, ok := core.VehicleModelGrade_value[data.Grade]; !ok || value <= 0 {
//		return errors.New("INVALID_GRADE")
//	}
//	if data.WarmUpTime <= 0 {
//		return errors.New("INVALID_WARM_UP_TIME")
//	}
//
//	if !data.Standard {
//		if data.StandardModelID == "" {
//			return errors.New("EMPTY_STANDARD_MODEL_ID")
//		}
//		if data.FuelType == "" {
//			return errors.New("EMPTY_FUEL_TYPE")
//		}
//	}
//
//	return nil
//}