package message

type VehicleModelListReq struct {
	IDs        string   `"json:"IDs,omitempty"`        // 모델 ID[]
	Brand      string   `"json:"brand"`                // 상표
	Brands     string   `"json:"brands,omitempty"`     // 상표[]
	Name       string   `"json:"name,omitempty"`       // 모델명
	FuelType   string   `"json:"fuelType,omitempty"`   // 연료 유형
	FuelTypes  string   `"json:"fuelTypes,omitempty"`  // 연료 유형[]
	Grade      string   `"json:"grade,omitempty"`      // 차량 등급
	Grades     string   `"json:"grades,omitempty"`     // 차량 등급[]
	Page       uint64   `"json:"page,omitempty"`       // 페이지 번호
	RowPerPage uint64   `"json:"rowPerPage,omitempty"` // 페이지 당 행수
	SortField  uint64   `"json:"sortField,omitempty"`  // 단일 정렬 필드, 1.name, 2.brand, 3.standard, 4.fuel_type, 5.fuel_efficiency, 6.fuel_tank_capacity, 7.displacement, 8.grade, 9.warm_up_time
	SortOrder  uint64   `"json:"sortOrder,omitempty"`  // 단일 정렬 방향, 1.asc, 2.desc
	SortFields []uint64 `"json:"sortFields,omitempty"` // 다중 정렬 필드 배열, [1,2,3]
	SortOrders []uint64 `"json:"sortOrders,omitempty"` // 다중 정렬 방향 배열, [1,2,3]
}