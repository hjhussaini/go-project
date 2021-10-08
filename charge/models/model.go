package models

import (
	"github.com/jinzhu/gorm"
)

type Charge struct {
	PhoneNo	string	`gorm:"column:phone_no;primaryKey" json:"phone_no"`
	Code	int	`gorm:"column:code;primaryKey;autoIncrement:false" json:"code"`
}

type ChargeList	[]Charge

type Model struct {
	database	*gorm.DB
}

func (model *Model) CountCharge(code int) int {
	var count int

	model.database.Model(&Charge{}).Where("code = ?", code).Count(&count)

	return count
}

func (model *Model) InsertCharge(charge Charge) error {
	return model.database.Create(charge).Error
}

func (model *Model) GetChargeList(code int) ChargeList {
	var list ChargeList

	model.database.Model(&Charge{}).Where("code = ?", code).Find(&list)

	return list
}

func New(database *gorm.DB) *Model {
	database.SingularTable(true)
	database.AutoMigrate(Charge{})

	return &Model{database: database}
}
