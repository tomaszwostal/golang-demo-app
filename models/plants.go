package models

import (
	"gorm.io/gorm"
)

type Plants struct {
	ID      uint    `gorm:"primary_key, auto_increment" json:"id"`
	Name    *string `json:"name"`
	Species *string `json:"species"`
	Plan    *string `json:"plan"`
}

func MigratePlants(db *gorm.DB) error {
	err := db.AutoMigrate(&Plants{})
	return err
}
