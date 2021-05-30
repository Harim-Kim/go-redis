package message

import (
	pb "OJT/core"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)
type VehicleModelListReq struct {
	//ID         []string	`"json:"ID`
	IDs        []string   `"json:"IDs,omitempty"`        // 모델 ID[]
	Brand      []string   `"json:"brand"`                // 상표
	Brands     []string   `"json:"brands,omitempty"`     // 상표[]
	Name       []string   `"json:"name,omitempty"`       // 모델명
	//Names 	   []string `"json:"names,omitempty"`
	FuelType   []string   `"json:"fuelType,omitempty"`   // 연료 유형
	FuelTypes  []string   `"json:"fuelTypes,omitempty"`  // 연료 유형[]
	Grade      []string   `"json:"grade,omitempty"`      // 차량 등급
	Grades     []string   `"json:"grades,omitempty"`     // 차량 등급[]

}

func  Decoder(r *http.Request ) pb.VehicleModelFilter{
	query := r.URL.Query()
	//fmt.Println("query:",query)
	jsonbody, err := json.Marshal(query)
	if err != nil{
		fmt.Println("marshal error:",err)
	}
	middle := VehicleModelListReq{}
	if err := json.Unmarshal(jsonbody, &middle); err != nil {
		fmt.Println("unmarshal middle error:",err)
	}
	//fmt.Println(middle)
	filter := pb.VehicleModelFilter{
		IDs: middle.IDs,
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

func validateRequest(r *http.Request){

}