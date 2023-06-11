package model

import (
	"gorm.io/gorm"
)

type CustomerStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *CustomerStore {
	return &CustomerStore{
		db: db,
	}
}

func (s *CustomerStore) CreateCustomer(name string) (*Customer, error) {
	c := &Customer{
		Name: name,
	}
	s.db.Create(c)
	return c, nil
}
