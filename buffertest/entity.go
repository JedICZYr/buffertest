package one

import (
	"time"
)

type DataStruct struct {
	Data            string  	`validate:"required" gorm:"-"`
	ReceiveTime		time.Time	`validate:"required" gorm:"-"`
}