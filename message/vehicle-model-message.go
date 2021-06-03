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
	Name             string  `json:"name,string"`
	Brand            string  `json:"brand,string"`
	Standard         bool    `json:"standard,string"`
	StandardModelID  string  `json:"standardModelID,string"`
	SeatingCapacity  uint32  `json:"seatingCapacity,string"`
	FuelType         string  `json:"fuelType",string`
	FuelEfficiency   float32 `json:"fuelEfficiency,string"`
	FuelTankCapacity uint32  `json:"fuelTankCapacity,string"`
	Displacement     uint32  `json:"displacement,string"`
	Grade            string  `json:"grade,string"`
	WarmUpTime       uint32  `json:"warmUpTime,string"`
}
type VehicleModelRequestWrapper struct {
	Name             string  `json:"name"`
	Brand            string  `json:"brand"`
	Standard         int     `json:"standard"`
	StandardModelID  string  `json:"standardModelID"`
	SeatingCapacity  uint32  `json:"seatingCapacity"`
	FuelType         string  `json:"fuelType",string`
	FuelEfficiency   float32 `json:"fuelEfficiency"`
	FuelTankCapacity uint32  `json:"fuelTankCapacity"`
	Displacement     uint32  `json:"displacement"`
	Grade            string  `json:"grade"`
	WarmUpTime       uint32  `json:"warmUpTime"`
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
	fmt.Println("validation process")
	wrapper := &VehicleModelRequestWrapper{}
	err := json.NewDecoder(r.Body).Decode(wrapper)
	if err != nil {
		fmt.Println("decode error",err)
		return err
	}
	data.Name = wrapper.Name
	data.Brand = wrapper.Brand
	if wrapper.Standard == 0{
		data.Standard = false
	} else {
		data.Standard = true
	}
	data.StandardModelID = wrapper.StandardModelID
	data.SeatingCapacity = wrapper.SeatingCapacity
	data.FuelType = wrapper.FuelType
	data.FuelEfficiency = wrapper.FuelEfficiency
	data.FuelTankCapacity = wrapper.FuelTankCapacity
	data.Displacement = wrapper.Displacement
	data.Grade = wrapper.Grade
	data.WarmUpTime = wrapper.WarmUpTime
	if data.Name == "" {
		fmt.Println("no name")
		return errors.New("EMPTY_NAME")
	}
	if data.Brand == "" {
		fmt.Println("no brand")
		return errors.New("EMPTY_BRAND")
	}
	if data.Grade == "" {
		fmt.Println("no grade")
		return errors.New("EMPTY_GRADE")
	} else if value, ok := pb.VehicleModelGrade_value[data.Grade]; !ok || value <= 0 {
		fmt.Println("no grade")
		return errors.New("INVALID_GRADE")
	}
	if data.WarmUpTime <= 0 {
		fmt.Println("no warm")
		return errors.New("INVALID_WARM_UP_TIME")
	}

	if !data.Standard {
		if data.StandardModelID == "" {
			fmt.Println("no EMPTY_STANDARD_MODEL_ID")
			return errors.New("EMPTY_STANDARD_MODEL_ID")
		}
		if data.FuelType == "" {
			fmt.Println("EMPTY_FUEL_TYPE")
			return errors.New("EMPTY_FUEL_TYPE")
		}
	}

	return nil
}
