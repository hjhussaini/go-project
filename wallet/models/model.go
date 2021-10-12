package models

import (
	"github.com/jinzhu/gorm"
)

type Wallet struct {
	PhoneNo	string	`gorm:"column:phone_no;primaryKey;unique" json:"phone_no"`
	Credit	int	`gorm:"column:credit" json:"credit"`
}

type Model struct {
	database	*gorm.DB
}

func (model *Model) InsertWallet(wallet Wallet) error {
	return model.database.Create(wallet).Error
}

func (model *Model) GetWallet(id string) *Wallet {
	var wallet Wallet

	model.database.Model(&Wallet{}).First(&wallet, &Wallet{PhoneNo: id})

	return &wallet
}

func New(database *gorm.DB) *Model {
	database.SingularTable(true)
	database.AutoMigrate(Wallet{})

	return &Model{database: database}
}
