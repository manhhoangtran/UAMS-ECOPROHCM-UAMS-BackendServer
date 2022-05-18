package models

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Area{},
		&Gateway{},
		&Doorlock{},
		&GatewayLog{},
		&Employee{},
		&Student{},
		&Customer{},
		&Scheduler{},
		&SecretKey{},
		&GwNetwork{},
		&DoorlockStatusLog{},
	)
	if err != nil {
		panic(err)
	}
}
