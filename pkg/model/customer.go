package model

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name string
}
