package model

import (
	"time"
	"github.com/google/uuid"

)

type Vehiclemodel struct {
	ID               uuid.UUID `gorm:"column:ID;type:char(36);primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `sql:"index"`
	CreatedBy        uuid.UUID  `gorm:"type:char(36)"`
	CreatorName      string     `gorm:"type:varchar(40)"`

	Name             string     `gorm:"column:name;not null"`      // 모델명
	Brand            string     `gorm:"column:brand"`              // 제조사 이름

	Standard         bool       `gorm:"column:standard"`           // 표준 모델 여부
	StandardModelID  string     `gorm:"column:standard_model_id"`  // 표준 모델ID

	SeatingCapacity  uint32     `gorm:"column:seating_capacity"`   // 탑승인원
	FuelType         string     `gorm:"column:fuel_type;"`         // 연료 유형 : 전기
	FuelEfficiency   float32    `gorm:"column:fuel_efficiency"`    // 연비
	FuelTankCapacity uint32     `gorm:"column:fuel_tank_capacity"` // 연료탱크용량 liter
	Displacement     uint32     `gorm:"column:displacement;"`      // 배기량 cc
	Grade            string     `gorm:"column:grade"`              // 차량 등록 경형/소형/준중형/중형/대형/SUV
	WarmUpTime       uint32     `gorm:"column:warm_up_time"`       // 차량 이용 예비시간
	ImageURL         string     `gorm:"column:image_url"`
	Image            Image      `gorm:"foreignkey:ID;association_autoupdate:false;association_autocreate:false"`
}
func(Vehiclemodel) TableName() string{
	return "vehicle_model_uuid"
}