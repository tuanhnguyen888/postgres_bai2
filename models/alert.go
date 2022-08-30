package models

import (
	"time"

	"gorm.io/gorm"
)

type Alert struct {
	Id          uint      `gorm:"primary key;autoIncreament" json:"id"`
	AlertId     *string   `json:"alertID"`
	Category    *string   `json:"category"`
	CloseBy     *string   `json:"closeBy"`
	Create      time.Time `json:"create"`
	Datetime    time.Time `json:"dateTime"`
	Description *string   `json:"description"`
	LastUpdate  time.Time `json:"lastUpdate"`
	Message     *string   `json:"message"`
	Object      *string   `json:"object"`
	ObjectType  *string   `json:"objectType"`
	Owner       *string   `json:"owner"`
	Status      *string   `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Type        *string   `json:"type"`
}

func MigrateAlert(db *gorm.DB) error {
	err := db.AutoMigrate(&Alert{})
	return err
}
